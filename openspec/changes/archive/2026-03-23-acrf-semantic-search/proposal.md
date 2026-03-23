## Why

Current skill discovery in ACRF relies on exact-string matching or simple keyword search, which is insufficient for 6G intent-driven networking. Natural language intents (e.g., "ensure device is reachable") may not explicitly contain the skill ID or tags. Implementing **Semantic Search** using vector embeddings and Cosine Similarity enables the ACRF to match intents to skills based on meaning, fulfilling the "Semantic Discovery" requirements of the 3GPP 6G Agentic Core proposal.

## What Changes

- **Update `SkillProfile` Model**: Add a `Description` field to provide semantic context for each skill.
- **Implement Gemini Embedding Client**: Add a lightweight HTTP client to fetch 768-dimensional vectors from the Gemini Embedding API (`text-embedding-004`).
- **Refactor `InMemoryRegistry`**: 
    - Change storage to a slice of structs containing both the `SkillProfile` and its embedding vector.
    - Fetch embeddings automatically during `Register`.
    - Perform semantic similarity search during `Discover` using Cosine Similarity.
- **Implement Math Helper**: Add a pure Go `cosineSimilarity` function.
- **Add Search Threshold**: Implement a similarity threshold (default `0.75`) to prevent low-confidence matches.

## Capabilities

### New Capabilities
- `semantic-matching-engine`: The core capability for generating embeddings and performing similarity-based discovery.

### Modified Capabilities
- `skill-registry`: Upgrade requirements to support semantic discovery instead of just identity-based lookup.
- `agentic-models`: Update `SkillProfile` to include the mandatory `Description` field.

## Impact

- **Performance**: Discovery now requires a network call to the Gemini API (unless cached) and a linear scan of vectors (acceptable for current MVP scale).
- **Configuration**: Requires `AGENTIC_GEMINI_API_KEY` to be set.
- **Breaking Change**: `SkillProfile` now requires a `Description` field for meaningful search.
