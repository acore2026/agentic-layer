## 1. Model Updates

- [x] 1.1 Add `Description` field to `SkillProfile` in `pkg/models/models.go`

## 2. Math and Embedding Utilities

- [x] 2.1 Implement `cosineSimilarity` function in `internal/registry/embeddings.go`
- [x] 2.2 Implement Gemini embedding client in `internal/registry/embeddings.go`
- [x] 2.3 Add unit tests for `cosineSimilarity` to ensure mathematical correctness

## 3. Registry Refactoring

- [x] 3.1 Update `InMemoryRegistry` struct to use a slice of entries (Profile + Embedding) in `internal/registry/registry.go`
- [x] 3.2 Update `Register` method to fetch and store embeddings for new skills
- [x] 3.3 Update `Discover` method to implement the semantic search logic with similarity threshold (0.75)

## 4. Verification

- [x] 4.1 Run system integration tests to verify semantic discovery works with natural language intents
- [x] 4.2 Verify that low-confidence queries (below 0.75) correctly return a 404
