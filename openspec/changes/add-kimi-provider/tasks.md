## 1. OpenAI Provider Implementation

- [x] 1.1 Create `internal/agent/openai` package
- [x] 1.2 Implement Moonshot/OpenAI request/response types
- [x] 1.3 Implement the `GenerateContent` method for the `model.LLM` interface
- [x] 1.4 Add tool-calling mapping logic (adk-go tools -> OpenAI tools)

## 2. Configuration and Environment

- [x] 2.1 Update `internal/config/config.go` to include Kimi/Moonshot variables
- [x] 2.2 Update `.env.example` with placeholders for Kimi configuration
- [x] 2.3 Ensure `AGENTIC_LLM_PROVIDER` can be "gemini" or "kimi"

## 3. Agent Integration

- [x] 3.1 Refactor `internal/agent/agent.go` to support factory-based provider selection
- [x] 3.2 Implement Kimi initialization in `NewCoreAgent`
- [x] 3.3 Verify that the Mock Agent still works correctly

## 4. Verification

- [x] 4.1 Run system integration tests using Kimi (if key provided)
- [x] 4.2 Verify successful end-to-end "Fleet Wake-Up" flow with Kimi
