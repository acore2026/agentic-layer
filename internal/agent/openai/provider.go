package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"

	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

type Provider struct {
	APIKey  string
	BaseURL string
	Model   string
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
	Role       string           `json:"role"`
	Content    string           `json:"content,omitempty"`
	ToolCalls  []openAIToolCall `json:"tool_calls,omitempty"`
	ToolCallID string           `json:"tool_call_id,omitempty"`
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
	return func(yield func(*model.LLMResponse, error) bool) {
		if stream {
			yield(nil, fmt.Errorf("streaming not supported yet by custom openai provider"))
			return
		}

		oaReq := openAIRequest{
			Model: p.Model,
		}

		for _, c := range req.Contents {
			msg := openAIMessage{Role: c.Role}
			if c.Role == "model" {
				msg.Role = "assistant"
			}
			
			hasContent := false
			for _, part := range c.Parts {
				if part.Text != "" {
					msg.Content = part.Text
					hasContent = true
				}
				if part.FunctionCall != nil {
					argsBytes, _ := json.Marshal(part.FunctionCall.Args)
					msg.ToolCalls = append(msg.ToolCalls, openAIToolCall{
						ID:   "call_" + part.FunctionCall.Name,
						Type: "function",
						Function: openAIFunctionCall{
							Name:      part.FunctionCall.Name,
							Arguments: string(argsBytes), 
						},
					})
					hasContent = true
				}
				if part.FunctionResponse != nil {
					msg.Role = "tool"
					msg.ToolCallID = "call_" + part.FunctionResponse.Name
					respBytes, _ := json.Marshal(part.FunctionResponse.Response)
					msg.Content = string(respBytes)
					hasContent = true
				}
			}
			
			if hasContent {
				oaReq.Messages = append(oaReq.Messages, msg)
			}
		}

		// Tool Injection
		oaReq.Tools = []openAITool{
			{
				Type: "function",
				Function: openAIFunctionDesc{
					Name:        "SearchSkill",
					Description: "Discover a network skill profile by its URI",
					Parameters: map[string]any{
						"type": "object",
						"properties": map[string]any{
							"skill_id": map[string]any{
								"type":        "string",
								"description": "The unique URI of the skill to discover (e.g., mcp://skill/device/fleet-update)",
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
								"description": "The URI of the skill to execute (must have been discovered first)",
							},
						},
						"required": []string{"skill_id"},
					},
				},
			},
		}

		body, _ := json.Marshal(oaReq)
		url := p.BaseURL + "/chat/completions"
		
		httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
		httpReq.Header.Set("Authorization", "Bearer "+p.APIKey)
		httpReq.Header.Set("Content-Type", "application/json")
		// Disable thinking to avoid 400 errors regarding missing reasoning_content
		httpReq.Header.Set("X-Msh-Next-Thinking", "0")

		resp, err := http.DefaultClient.Do(httpReq)
		if err != nil {
			yield(nil, err)
			return
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
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
