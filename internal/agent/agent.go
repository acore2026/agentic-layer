package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"log"
	"os"
	"strings"

	"github.com/google/6g-agentic-core/internal/agent/openai"
	"github.com/google/6g-agentic-core/internal/events"
	"github.com/google/6g-agentic-core/pkg/models"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

const SystemInstruction = `You are a 6G Skill-Based Agentic Core Network Reasoner. 
Your goal is to resolve user intents into deterministic network actions using the provided tools.

You MUST follow this Three-Stage Execution Pipeline:
1. INTENT: Receive the user's abstract goal (e.g., "Wake up the fleet").
2. SKILL DISCOVERY: You MUST call the 'SearchSkill' tool to find a matching skill URI from the ACRF registry
   - If the user provides a natural language goal, use key terms from that goal as the 'skill_id' to search.
   - Example: For "Wake up the fleet", search for 'mcp://skill/device/fleet-update' or simply 'fleet-update'.
   - Even if you think you know the URI, you MUST verify it exists via 'SearchSkill'.
   - NEVER ask the user for the URI. That is YOUR job to find.
3. SERVICE DIRECTIVE: Once you have the skill profile from 'SearchSkill', you MUST call 'ExecuteSkill' with that URI to invoke the action via the A-IGW.

DO NOT say you cannot help until you have at least tried to search for a skill.
Always confirm the final results of the execution to the user.`

// emitEvent sends a JSON event to the SSE broker if initialized.
func emitEvent(broker *events.Broker, eventType string, data any) {
	if broker == nil {
		return
	}
	payload := map[string]any{
		"type": eventType,
		"data": data,
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("SSE: Failed to marshal event: %v", err)
		return
	}
	
	// Non-blocking send
	select {
	case broker.Notifier <- jsonBytes:
	default:
		log.Println("SSE: Broker notifier channel full, dropping event")
	}
}

func NewCoreAgent(ctx context.Context, broker *events.Broker) (agent.Agent, error) {
	if os.Getenv("AGENTIC_USE_MOCK_AGENT") == "true" {
		log.Println("Using MockCoreAgent as requested by AGENTIC_USE_MOCK_AGENT=true")
		return NewMockCoreAgent(broker)
	}

	provider := os.Getenv("AGENTIC_LLM_PROVIDER")
	if provider == "" {
		provider = "gemini"
	}

	var llm model.LLM
	var err error

	switch provider {
	case "kimi":
		apiKey := os.Getenv("AGENTIC_KIMI_API_KEY")
		baseURL := os.Getenv("AGENTIC_KIMI_BASE_URL")
		modelName := os.Getenv("AGENTIC_KIMI_MODEL")
		if baseURL == "" {
			baseURL = "https://api.moonshot.cn/v1"
		}
		if modelName == "" {
			modelName = "kimi-k2.5"
		}
		llm = &openai.Provider{
			APIKey:  apiKey,
			BaseURL: baseURL,
			Model:   modelName,
		}
		log.Printf("Initialized Kimi provider (Model: %s)", modelName)

	case "gemini":
		apiKey := os.Getenv("AGENTIC_GEMINI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("AGENTIC_GEMINI_API_KEY environment variable is required for Gemini")
		}
		clientCfg := &genai.ClientConfig{
			APIKey:  apiKey,
			Backend: genai.BackendGeminiAPI,
		}
		llm, err = gemini.NewModel(ctx, "gemini-1.5-flash", clientCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create gemini model: %v", err)
		}
		log.Println("Initialized Gemini provider")

	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", provider)
	}

	searchTool, err := functiontool.New(functiontool.Config{
		Name:        "SearchSkill",
		Description: "Discover a network skill profile by its URI or keyword",
	}, func(ctx tool.Context, input SearchSkillInput) (string, error) {
		emitEvent(broker, "tool_call_started", map[string]string{"tool": "SearchSkill", "skill_id": input.SkillID})
		res, err := SearchSkill(ctx, input)
		emitEvent(broker, "tool_call_completed", map[string]any{"tool": "SearchSkill", "result": res, "error": err})
		return res, err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create search tool: %v", err)
	}

	executeTool, err := functiontool.New(functiontool.Config{
		Name:        "ExecuteSkill",
		Description: "Invoke a discovered network skill via the Interworking Gateway",
	}, func(ctx tool.Context, input ExecuteSkillInput) (string, error) {
		emitEvent(broker, "tool_call_started", map[string]string{"tool": "ExecuteSkill", "skill_id": input.SkillID})
		res, err := ExecuteSkill(ctx, input)
		emitEvent(broker, "tool_call_completed", map[string]any{"tool": "ExecuteSkill", "result": res, "error": err})
		return res, err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create execute tool: %v", err)
	}

	cfg := llmagent.Config{
		Name:        "CoreReasoner",
		Description: "Reasoning engine for 6G intent resolution",
		Instruction: SystemInstruction,
		Model:       llm,
		Tools:       []tool.Tool{searchTool, executeTool},
	}

	// Wrap the core agent logic to emit start/end events
	baseAgent, err := llmagent.New(cfg)
	if err != nil {
		return nil, err
	}

	return agent.New(agent.Config{
		Name: "CoreReasonerWithEvents",
		Run: func(ctx agent.InvocationContext) iter.Seq2[*session.Event, error] {
			prompt := ""
			if ctx.UserContent() != nil && len(ctx.UserContent().Parts) > 0 {
				prompt = ctx.UserContent().Parts[0].Text
			}
			emitEvent(broker, "reasoning_started", map[string]string{"prompt": prompt})

			return func(yield func(*session.Event, error) bool) {
				var finalResponse string
				for event, err := range baseAgent.Run(ctx) {
					if err != nil {
						emitEvent(broker, "reasoning_error", map[string]string{"error": err.Error()})
						if !yield(nil, err) { return }
						continue
					}
					if event.Content != nil && len(event.Content.Parts) > 0 {
						finalResponse += event.Content.Parts[0].Text
					}
					if !yield(event, nil) { return }
				}
				emitEvent(broker, "reasoning_completed", map[string]string{"response": finalResponse})
			}
		},
	})
}

// NewMockCoreAgent creates a deterministic mock agent for testing.
func NewMockCoreAgent(broker *events.Broker) (agent.Agent, error) {
	return agent.New(agent.Config{
		Name:        "CoreReasoner",
		Description: "Mock reasoning engine for 6G intent resolution",
		Run: func(ctx agent.InvocationContext) iter.Seq2[*session.Event, error] {
			return func(yield func(*session.Event, error) bool) {
				prompt := ""
				if ctx.UserContent() != nil && len(ctx.UserContent().Parts) > 0 {
					prompt = ctx.UserContent().Parts[0].Text
				}

				emitEvent(broker, "reasoning_started", map[string]string{"prompt": prompt})
				log.Printf("[MockAgent] Processing prompt: %s", prompt)

				// Determine SkillID based on prompt
				skillID := "mcp://skill/device/fleet-update" // Default
				if strings.Contains(strings.ToLower(prompt), "turbo") {
					skillID = "mcp://skill/qos/turbo-mode"
				} else if strings.Contains(strings.ToLower(prompt), "interruption") || strings.Contains(strings.ToLower(prompt), "v2x") {
					skillID = "mcp://skill/reliability/path-diversity"
				} else if strings.Contains(strings.ToLower(prompt), "drone") || strings.Contains(strings.ToLower(prompt), "flight") {
					skillID = "mcp://skill/edge/secure-flight"
				} else if strings.Contains(strings.ToLower(prompt), "pizza") {
					skillID = "mcp://skill/cook/pizza"
				}

				// 1. Discovery Step
				emitEvent(broker, "discovery_started", map[string]string{"query": skillID})
				log.Printf("[MockAgent] Calling SearchSkill for %s", skillID)
				profile, err := SearchSkill(context.Background(), SearchSkillInput{SkillID: skillID})
				if err != nil {
					emitEvent(broker, "discovery_failed", map[string]string{"error": err.Error()})
					yield(nil, fmt.Errorf("discovery failed: %v", err))
					return
				}
				
				if strings.Contains(profile, "not found") {
					emitEvent(broker, "skill_not_found", map[string]string{"skill_id": skillID})
					log.Printf("[MockAgent] Discovery result: %s", profile)
					finalMsg := fmt.Sprintf("Mock result: I couldn't find a skill for %s", skillID)
					event := session.NewEvent(ctx.InvocationID())
					event.Content = &genai.Content{
						Parts: []*genai.Part{{Text: finalMsg}},
						Role:  "model",
					}
					emitEvent(broker, "reasoning_completed", map[string]string{"response": finalMsg})
					yield(event, nil)
					return
				}
				
				log.Printf("[MockAgent] Discovery result: %s", profile)
				var profileObj models.SkillProfile
				json.Unmarshal([]byte(profile), &profileObj)
				emitEvent(broker, "skill_discovered", profileObj)

				// 2. Invocation Step
				emitEvent(broker, "invocation_started", map[string]string{"skill_id": skillID})
				log.Printf("[MockAgent] Calling ExecuteSkill for %s", skillID)
				result, err := ExecuteSkill(context.Background(), ExecuteSkillInput{SkillID: skillID})
				if err != nil {
					emitEvent(broker, "invocation_failed", map[string]string{"error": err.Error()})
					finalMsg := fmt.Sprintf("Mock result: execution failed for %s: %v", skillID, err)
					event := session.NewEvent(ctx.InvocationID())
					event.Content = &genai.Content{
						Parts: []*genai.Part{{Text: finalMsg}},
						Role:  "model",
					}
					emitEvent(broker, "reasoning_completed", map[string]string{"response": finalMsg})
					yield(event, nil)
					return
				}
				log.Printf("[MockAgent] Invocation result: %s", result)
				emitEvent(broker, "invocation_completed", map[string]string{"result": result})

				// 3. Final Response
				finalMsg := fmt.Sprintf("Mock result: successfully triggered %s", skillID)
				event := session.NewEvent(ctx.InvocationID())
				event.Content = &genai.Content{
					Parts: []*genai.Part{{Text: finalMsg}},
					Role:  "model",
				}
				emitEvent(broker, "reasoning_completed", map[string]string{"response": finalMsg})
				yield(event, nil)
			}
		},
	})
}
