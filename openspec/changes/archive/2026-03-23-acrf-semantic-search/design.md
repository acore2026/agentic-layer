## Context

The current ACRF uses an exact-match registry. 3GPP TS 23.501/23.502 equivalents for 6G are evolving towards semantic discovery. This design introduces a vector-based search mechanism within the `InMemoryRegistry` to support natural language intent resolution.

## Goals / Non-Goals

**Goals:**
- Implement Cosine Similarity search for skill discovery.
- Integrate with Gemini Embedding API (`text-embedding-004`).
- Maintain a thread-safe in-memory store for profiles and vectors.
- Ensure backwards compatibility where possible (though `Description` becomes practically required).

**Non-Goals:**
- Persistent vector database (e.g., Pinecone, Milvus) - we will stay in-memory for the MVP.
- Complex NLP preprocessing (stemming, lemmatization) - we rely on the embedding model's quality.

## Decisions

- **Embedding Client**: A lightweight function using `net/http` to call the Gemini API.
    - *Alternative*: Using a full GenAI SDK.
    - *Rationale*: Minimal dependencies are preferred for this microservice.
- **Storage Structure**: 
    ```go
    type skillEntry struct {
        Profile   models.SkillProfile
        Embedding []float32
    }
    ```
    The `InMemoryRegistry` will hold `[]skillEntry` instead of `map[string]SkillProfile`.
- **Search Logic**:
    1. Embed the search query.
    2. Iterate through all `skillEntry` items.
    3. Calculate Cosine Similarity.
    4. Keep track of the entry with the highest score.
    5. Return result if `score > 0.75`.
- **Concatenation Strategy**: For embedding a registered skill, we will use `fmt.Sprintf("%s: %s", profile.SkillID, profile.Description)`. This ensures both the unique ID and the descriptive text contribute to the vector representation.

## Risks / Trade-offs

- **[Risk] API Latency**: Every discovery request now incurs a network hop to the Gemini API.
    - *Mitigation*: In a production environment, we would use local embedding models or aggressive caching. For the MVP, we accept the latency.
- **[Risk] Cost**: Embedding calls cost tokens.
    - *Mitigation*: Skill registration happens once; search frequency should be monitored.
- **[Trade-off] Linear Scan**: $O(N)$ search time.
    - *Rationale*: For < 1000 skills, a linear scan of float32 slices is extremely fast in Go and simpler than implementing an HNSW index.
