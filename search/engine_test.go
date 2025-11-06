package search

import (
	"math"
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []float64
		expected float64
		delta    float64
	}{
		{
			name:     "identical vectors",
			a:        []float64{1.0, 2.0, 3.0},
			b:        []float64{1.0, 2.0, 3.0},
			expected: 1.0,
			delta:    0.0001,
		},
		{
			name:     "orthogonal vectors",
			a:        []float64{1.0, 0.0},
			b:        []float64{0.0, 1.0},
			expected: 0.0,
			delta:    0.0001,
		},
		{
			name:     "opposite vectors",
			a:        []float64{1.0, 0.0},
			b:        []float64{-1.0, 0.0},
			expected: -1.0,
			delta:    0.0001,
		},
		{
			name:     "different lengths",
			a:        []float64{1.0, 2.0},
			b:        []float64{1.0, 2.0, 3.0},
			expected: 0.0,
			delta:    0.0001,
		},
		{
			name:     "zero vector",
			a:        []float64{0.0, 0.0},
			b:        []float64{1.0, 2.0},
			expected: 0.0,
			delta:    0.0001,
		},
		{
			name:     "similar vectors",
			a:        []float64{1.0, 2.0, 3.0},
			b:        []float64{1.1, 2.1, 2.9},
			expected: 0.999,
			delta:    0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cosineSimilarity(tt.a, tt.b)
			if math.Abs(result-tt.expected) > tt.delta {
				t.Errorf("cosineSimilarity(%v, %v) = %f, want %f (±%f)",
					tt.a, tt.b, result, tt.expected, tt.delta)
			}
		})
	}
}

func TestCosineSimilarityEdgeCases(t *testing.T) {
	// Empty vectors
	result := cosineSimilarity([]float64{}, []float64{})
	if result != 0.0 {
		t.Errorf("Empty vectors should return 0.0, got %f", result)
	}

	// Nil vectors
	result = cosineSimilarity(nil, nil)
	if result != 0.0 {
		t.Errorf("Nil vectors should return 0.0, got %f", result)
	}

	// Single element vectors
	result = cosineSimilarity([]float64{5.0}, []float64{3.0})
	if result != 1.0 {
		t.Errorf("Single positive element vectors should return 1.0, got %f", result)
	}
}

func BenchmarkCosineSimilarity(b *testing.B) {
	// Test with typical embedding dimensions (384 for all-minilm)
	vec1 := make([]float64, 384)
	vec2 := make([]float64, 384)

	for i := range vec1 {
		vec1[i] = float64(i) * 0.1
		vec2[i] = float64(i) * 0.11 // Slightly different
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cosineSimilarity(vec1, vec2)
	}
}

// Test helper to create mock embeddings
func createMockEmbedding(pattern string) []float64 {
	// Create deterministic embeddings based on pattern
	// This simulates different types of content
	embedding := make([]float64, 10) // Smaller for testing

	switch pattern {
	case "meeting":
		embedding = []float64{0.1, 0.9, 0.2, 0.8, 0.1, 0.7, 0.3, 0.6, 0.2, 0.5}
	case "python":
		embedding = []float64{0.8, 0.1, 0.9, 0.2, 0.7, 0.3, 0.6, 0.4, 0.8, 0.1}
	case "funny":
		embedding = []float64{0.3, 0.7, 0.1, 0.6, 0.4, 0.8, 0.2, 0.9, 0.1, 0.7}
	case "deadline":
		embedding = []float64{0.2, 0.8, 0.3, 0.7, 0.1, 0.9, 0.4, 0.6, 0.3, 0.8}
	default:
		// Random-ish pattern
		for i := range embedding {
			embedding[i] = float64(i) * 0.1
		}
	}

	return embedding
}

func TestSearchRelevance(t *testing.T) {
	// Test that similar content types have higher similarity
	meetingEmb1 := createMockEmbedding("meeting")
	meetingEmb2 := createMockEmbedding("meeting")
	pythonEmb := createMockEmbedding("python")

	// Same type should have high similarity (use delta for floating point comparison)
	sameSimilarity := cosineSimilarity(meetingEmb1, meetingEmb2)
	if math.Abs(sameSimilarity-1.0) > 0.0001 {
		t.Errorf("Identical embeddings should have similarity ~1.0, got %f", sameSimilarity)
	}

	// Different types should have lower similarity
	diffSimilarity := cosineSimilarity(meetingEmb1, pythonEmb)
	if diffSimilarity >= sameSimilarity {
		t.Errorf("Different content should have lower similarity than same content")
	}
}
