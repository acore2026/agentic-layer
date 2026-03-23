package registry

import (
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float32
	}{
		{
			name:     "Identical vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{1, 0, 0},
			expected: 1.0,
		},
		{
			name:     "Orthogonal vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{0, 1, 0},
			expected: 0.0,
		},
		{
			name:     "Opposite vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{-1, 0, 0},
			expected: -1.0,
		},
		{
			name:     "Simple similar vectors",
			a:        []float32{1, 1},
			b:        []float32{1, 0},
			expected: 0.70710677, // 1/sqrt(2)
		},
		{
			name:     "Empty vectors",
			a:        []float32{},
			b:        []float32{},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cosineSimilarity(tt.a, tt.b)
			if mathAbs(got-tt.expected) > 1e-6 {
				t.Errorf("cosineSimilarity() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func mathAbs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
