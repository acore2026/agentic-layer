## ADDED Requirements

### Requirement: Vector Generation
The system SHALL be able to convert text strings into numerical vector embeddings using an external embedding provider (Gemini).

#### Scenario: Successful Embedding
- **WHEN** a text string is provided to the embedding client
- **THEN** it SHALL return a slice of float32 values representing the semantic vector.

### Requirement: Similarity Calculation
The system SHALL calculate the Cosine Similarity between two embedding vectors to determine their semantic closeness.

#### Scenario: High Similarity
- **WHEN** two vectors representing similar concepts are compared
- **THEN** the resulting score SHALL be close to 1.0.

### Requirement: Search Thresholding
The semantic discovery process SHALL only return results that meet or exceed a configurable similarity threshold (default 0.75).

#### Scenario: No Match Found
- **WHEN** no registered skill has a similarity score greater than 0.75 for a given query
- **THEN** the discovery process SHALL return a "not found" result.
