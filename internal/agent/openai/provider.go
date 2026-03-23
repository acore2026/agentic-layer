package openai

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"log"
	"net/http"
	"sync"

	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

type Provider struct {
	APIKey  string
	BaseURL string
	Model   string

	// Side-channel to remember reasoning content for Moonshot multi-turn consistency
	mu             sync.Mutex
	reasoningCache map[string]string
}

func (p *Provider) Name() string {
	return "openai-compatible"
}

type openAIRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
	Tools    []openAITool    `json:"tools,omitempty"`
}

type openAIMessage struct {
	Role             string           `json:"role"`
	Content          string           `json:"content,omitempty"`
	ReasoningContent string           `json:"reasoning_content,omitempty"`
	ToolCalls        []openAIToolCall `json:"tool_calls,omitempty"`
	ToolCallID       string           `json:"tool_call_id,omitempty"`
}

func (m *openAIMessage) getFingerprint() string {
	h := sha256.New()
	h.Write([]byte(m.Role))
	h.Write([]byte(m.Content))
	for _, tc := range m.ToolCalls {
		h.Write([]byte(tc.Function.Name))
		var args map[string]any
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err == nil {
			canon, _ := json.Marshal(args)
			h.Write(canon)
		} else {
			h.Write([]byte(tc.Function.Arguments))
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

type openAITool struct {
	Type     string             `json:"type"`
	Function openAIFunctionDesc `json:"function"`
}

type openAIFunctionDesc struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parameters  any    `json:"parameters,omitempty"`
}

type openAIToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function openAIFunctionCall `json:"function"`
}

type openAIFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type openAIResponse struct {
	Choices []struct {
		Message openAIMessage `json:"message"`
	} `json:"choices"`
}

func (p *Provider) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	if p.reasoningCache == nil {
		p.mu.Lock()
		if p.reasoningCache == nil {
			p.reasoningCache = make(map[string]string)
		}
		p.mu.Unlock()
	}

	return func(yield func(*model.LLMResponse, error) bool) {
		if stream {
			yield(nil, fmt.Errorf("streaming not supported yet by custom openai provider"))
			return
		}

		oaReq := openAIRequest{
			Model: p.Model,
		}

		toolCounts := make(map[string]int)

		for i, c := range req.Contents {
			if c.Role == "model" {
				msg := openAIMessage{Role: "assistant"}
				toolCallInThisTurn := 0
				
				for _, part := range c.Parts {
					if part.Text != "" {
						msg.Content = part.Text
					}
					if part.FunctionCall != nil {
						callID := fmt.Sprintf("call_%s_%d", part.FunctionCall.Name, toolCallInThisTurn)
						toolCallInThisTurn++
						
						argsBytes, _ := json.Marshal(part.FunctionCall.Args)
						msg.ToolCalls = append(msg.ToolCalls, openAIToolCall{
							ID:   callID,
							Type: "function",
							Function: openAIFunctionCall{
								Name:      part.FunctionCall.Name,
								Arguments: string(argsBytes), 
							},
						})
					}
				}
				
				fp := msg.getFingerprint()
				p.mu.Lock()
				if rc, ok := p.reasoningCache[fp]; ok {
					msg.ReasoningContent = rc
				}
				p.mu.Unlock()
				
				oaReq.Messages = append(oaReq.Messages, msg)
				toolCounts = make(map[string]int)
				continue
			}

			hasToolResponse := false
			for _, part := range c.Parts {
				if part.FunctionResponse != nil {
					idx := toolCounts[part.FunctionResponse.Name]
					toolCounts[part.FunctionResponse.Name]++
					
					callID := fmt.Sprintf("call_%s_%d", part.FunctionResponse.Name, idx)
					respBytes, _ := json.Marshal(part.FunctionResponse.Response)
					
					oaReq.Messages = append(oaReq.Messages, openAIMessage{
						Role:       "tool",
						ToolCallID: callID,
						Content:    string(respBytes),
					})
					hasToolResponse = true
				}
			}

			if !hasToolResponse {
				content := ""
				if len(c.Parts) > 0 {
					content = c.Parts[0].Text
				}
				if content != "" || i == 0 {
					oaReq.Messages = append(oaReq.Messages, openAIMessage{
						Role:    "user",
						Content: content,
					})
				}
			}
		}

		oaReq.Tools = []openAITool{
			{
				Type: "function",
				Function: openAIFunctionDesc{
					Name:        "SearchSkill",
					Description: "Discover a network skill profile by its URI or keyword",
					Parameters: map[string]any{
						"type": "object",
						"properties": map[string]any{
							"skill_id": map[string]any{
								"type":        "string",
								"description": "The unique URI or keyword of the skill to discover",
							},
						},
						"required": []string{"skill_id"},
					},
				},
			},
			{
				Type: "function",
				Function: openAIFunctionDesc{
					Name:        "ExecuteSkill",
					Description: "Invoke a discovered network skill via the Interworking Gateway",
					Parameters: map[string]any{
						"type": "object",
						"properties": map[string]any{
							"skill_id": map[string]any{
								"type":        "string",
								"description": "The URI of the skill to execute",
							},
						},
						"required": []string{"skill_id"},
					},
				},
			},
		}

		body, _ := json.Marshal(oaReq)
		url := p.BaseURL + "/chat/completions"
		log.Printf("[OpenAI] Request: %s", string(body))
		
		httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
		httpReq.Header.Set("Authorization", "Bearer "+p.APIKey)
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("X-Msh-Next-Thinking", "0")

		resp, err := http.DefaultClient.Do(httpReq)
		if err != nil {
			log.Printf("[OpenAI] Request failed: %v", err)
			yield(nil, err)
			return
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		log.Printf("[OpenAI] Response (%d): %s", resp.StatusCode, string(respBody))

		if resp.StatusCode != http.StatusOK {
			yield(nil, fmt.Errorf("OpenAI API error (%d): %s", resp.StatusCode, string(respBody)))
			return
		}

		var res openAIResponse
		if err := json.Unmarshal(respBody, &res); err != nil {
			yield(nil, err)
			return
		}

		if len(res.Choices) == 0 {
			yield(nil, fmt.Errorf("empty response from OpenAI"))
			return
		}

		choice := res.Choices[0].Message
		role := choice.Role
		if role == "assistant" {
			role = "model"
		}

		if choice.ReasoningContent != "" {
			fp := choice.getFingerprint()
			p.mu.Lock()
			p.reasoningCache[fp] = choice.ReasoningContent
			p.mu.Unlock()
		}

		response := &model.LLMResponse{
			Content: &genai.Content{
				Role: role,
			},
		}

		if choice.Content != "" {
			response.Content.Parts = append(response.Content.Parts, &genai.Part{Text: choice.Content})
		}

		if len(choice.ToolCalls) > 0 {
			for _, tc := range choice.ToolCalls {
				var args map[string]any
				if tc.Function.Arguments != "" {
					json.Unmarshal([]byte(tc.Function.Arguments), &args)
				}
				response.Content.Parts = append(response.Content.Parts, &genai.Part{
					FunctionCall: &genai.FunctionCall{
						Name: tc.Function.Name,
						Args: args,
					},
				})
			}
		}

		yield(response, nil)
	}
}
