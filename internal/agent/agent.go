package agent

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

const SystemInstruction = `You are a 6G Skill-Based Agentic Core Network Reasoner. 
Your goal is to resolve user intents into deterministic network actions using the provided tools.

Follow this Three-Stage Execution Pipeline:
1. INTENT: Receive and understand the user's abstract goal (e.g., "Wake up the fleet").
2. SKILL: Use 'SearchSkill' to find a matching skill URI from the ACRF registry. 
   - If you don't know the exact URI, try common ones like 'mcp://skill/device/fleet-update'.
   - You MUST have a discovered skill profile before proceeding to execution.
3. SERVICE DIRECTIVE: Use 'ExecuteSkill' to invoke the discovered skill via the A-IGW.

Always confirm the results of your discovery and execution to the user.`

func NewCoreAgent(ctx context.Context) (agent.Agent, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
	}

	clientCfg := &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	}

	// Note: using gemini-1.5-flash as it is most likely to be available on free tier.
	model, err := gemini.NewModel(ctx, "gemini-1.5-flash", clientCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini model: %v", err)
	}

	searchTool, err := functiontool.New(functiontool.Config{
		Name:        "SearchSkill",
		Description: "Discover a network skill profile by its URI (e.g., mcp://skill/device/fleet-update)",
	}, SearchSkill)
	if err != nil {
		return nil, fmt.Errorf("failed to create search tool: %v", err)
	}

	executeTool, err := functiontool.New(functiontool.Config{
		Name:        "ExecuteSkill",
		Description: "Invoke a discovered network skill via the Interworking Gateway",
	}, ExecuteSkill)
	if err != nil {
		return nil, fmt.Errorf("failed to create execute tool: %v", err)
	}

	cfg := llmagent.Config{
		Name:        "CoreReasoner",
		Description: "Reasoning engine for 6G intent resolution",
		Instruction: SystemInstruction,
		Model:       model,
		Tools:       []tool.Tool{searchTool, executeTool},
	}

	return llmagent.New(cfg)
}
