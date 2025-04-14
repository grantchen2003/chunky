package byterange

import (
	"fmt"
	"testing"
)

func TestByteRange_NewByteRange(t *testing.T) {
	tests := []struct {
		name              string
		startByte         int
		endByte           int
		expectedByteRange ByteRange
		expectedError     error
	}{
		{
			name:              "invalid input",
			startByte:         15,
			endByte:           5,
			expectedByteRange: ByteRange{},
			expectedError:     fmt.Errorf("endByte: 5 is less than startByte: 15"),
		},
		{
			name:              "valid input",
			startByte:         0,
			endByte:           10,
			expectedByteRange: ByteRange{StartByte: 0, EndByte: 10},
			expectedError:     nil,
		},
		{
			name:              "equal startByte and endByte",
			startByte:         5,
			endByte:           5,
			expectedByteRange: ByteRange{StartByte: 5, EndByte: 5},
			expectedError:     nil,
		},
	}

	for _, tt := range tests {
		test := func(t *testing.T) {
			byteRange, err := NewByteRange(tt.startByte, tt.endByte)

			if (err == nil) != (tt.expectedError == nil) || (err != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("NewByteRange(%d, %d) error = '%v', expected error = '%v'",
					tt.startByte, tt.endByte, err, tt.expectedError)
			}

			if byteRange != tt.expectedByteRange {
				t.Errorf("NewByteRange(%d, %d) = %v, want %v",
					tt.startByte, tt.endByte, byteRange, tt.expectedByteRange)
			}
		}

		t.Run(tt.name, test)
	}
}

func TestByteRange_intersects(t *testing.T) {
	tests := []struct {
		name     string
		br1      ByteRange
		br2      ByteRange
		expected bool
	}{
		{
			name:     "Partial intersect: br2 starts inside br1",
			br1:      ByteRange{0, 10},
			br2:      ByteRange{5, 15},
			expected: true,
		},
		{
			name:     "Partial intersect: br1 starts inside br2",
			br1:      ByteRange{5, 15},
			br2:      ByteRange{0, 10},
			expected: true,
		},
		{
			name:     "Intersecting at boundary: br1 ends where br2 starts",
			br1:      ByteRange{0, 5},
			br2:      ByteRange{5, 10},
			expected: true,
		},
		{
			name:     "Intersecting at boundary: br2 ends where br1 starts",
			br1:      ByteRange{5, 10},
			br2:      ByteRange{0, 5},
			expected: true,
		},
		{
			name:     "br2 is fully contained within br1",
			br1:      ByteRange{0, 10},
			br2:      ByteRange{2, 8},
			expected: true,
		},
		{
			name:     "br1 is fully contained within br2",
			br1:      ByteRange{2, 8},
			br2:      ByteRange{0, 10},
			expected: true,
		},
		{
			name:     "Disjoint: br1 ends before br2 starts",
			br1:      ByteRange{0, 4},
			br2:      ByteRange{5, 10},
			expected: false,
		},
		{
			name:     "Disjoint: br2 ends before br1 starts",
			br1:      ByteRange{5, 10},
			br2:      ByteRange{0, 4},
			expected: false,
		},
		{
			name:     "Identical ranges",
			br1:      ByteRange{3, 7},
			br2:      ByteRange{3, 7},
			expected: true,
		},
	}

	for _, tt := range tests {
		test := func(t *testing.T) {
			result := intersects(tt.br1, tt.br2)
			if result != tt.expected {
				t.Errorf("For %s, expected %v, but got %v", tt.name, tt.expected, result)
			}
		}
		t.Run(tt.name, test)
	}
}

func TestByteRange_mergePair(t *testing.T) {
	tests := []struct {
		name          string
		br1           ByteRange
		br2           ByteRange
		expected      ByteRange
		expectedError error
	}{
		{
			name:          "Partial intersect: br2 starts inside br1",
			br1:           ByteRange{0, 10},
			br2:           ByteRange{5, 15},
			expected:      ByteRange{0, 15},
			expectedError: nil,
		},
		{
			name:          "Partial intersect: br1 starts inside br2",
			br1:           ByteRange{5, 15},
			br2:           ByteRange{0, 10},
			expected:      ByteRange{0, 15},
			expectedError: nil,
		},
		{
			name:          "Intersecting at boundary: br1 ends where br2 starts",
			br1:           ByteRange{0, 5},
			br2:           ByteRange{5, 10},
			expected:      ByteRange{0, 10},
			expectedError: nil,
		},
		{
			name:          "Intersecting at boundary: br2 ends where br1 starts",
			br1:           ByteRange{5, 10},
			br2:           ByteRange{0, 5},
			expected:      ByteRange{0, 10},
			expectedError: nil,
		},
		{
			name:          "br2 is fully contained within br1",
			br1:           ByteRange{0, 10},
			br2:           ByteRange{2, 8},
			expected:      ByteRange{0, 10},
			expectedError: nil,
		},
		{
			name:          "br1 is fully contained within br2",
			br1:           ByteRange{2, 8},
			br2:           ByteRange{0, 10},
			expected:      ByteRange{0, 10},
			expectedError: nil,
		},
		{
			name:          "Disjoint: br1 ends before br2 starts",
			br1:           ByteRange{0, 4},
			br2:           ByteRange{5, 10},
			expected:      ByteRange{},
			expectedError: fmt.Errorf("byte ranges {0 4} and {5 10} don't intersect"),
		},
		{
			name:          "Disjoint: br2 ends before br1 starts",
			br1:           ByteRange{5, 10},
			br2:           ByteRange{0, 4},
			expected:      ByteRange{},
			expectedError: fmt.Errorf("byte ranges {5 10} and {0 4} don't intersect"),
		},
		{
			name:          "Identical ranges",
			br1:           ByteRange{3, 7},
			br2:           ByteRange{3, 7},
			expected:      ByteRange{3, 7},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		test := func(t *testing.T) {
			result, err := mergePair(tt.br1, tt.br2)

			if result != tt.expected {
				t.Errorf("For %s, expected '%v', but got '%v'", tt.name, tt.expected, result)
			}

			if (err == nil) != (tt.expectedError == nil) || (err != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("For %s, expected '%v', but got '%v'", tt.name, tt.expectedError, err)
			}
		}

		t.Run(tt.name, test)
	}
}

func TestByterange_sortByteRanges(t *testing.T) {
	tests := []struct {
		name     string
		input    []ByteRange
		expected []ByteRange
	}{
		{
			name: "Already sorted",
			input: []ByteRange{
				{StartByte: 0, EndByte: 5},
				{StartByte: 5, EndByte: 10},
				{StartByte: 10, EndByte: 15},
			},
			expected: []ByteRange{
				{StartByte: 0, EndByte: 5},
				{StartByte: 5, EndByte: 10},
				{StartByte: 10, EndByte: 15},
			},
		},
		{
			name: "Unsorted by start byte",
			input: []ByteRange{
				{StartByte: 10, EndByte: 15},
				{StartByte: 0, EndByte: 5},
				{StartByte: 5, EndByte: 10},
			},
			expected: []ByteRange{
				{StartByte: 0, EndByte: 5},
				{StartByte: 5, EndByte: 10},
				{StartByte: 10, EndByte: 15},
			},
		},
		{
			name: "Same start byte, different end byte",
			input: []ByteRange{
				{StartByte: 5, EndByte: 10},
				{StartByte: 5, EndByte: 7},
				{StartByte: 5, EndByte: 12},
				{StartByte: 5, EndByte: 7},
			},
			expected: []ByteRange{
				{StartByte: 5, EndByte: 7},
				{StartByte: 5, EndByte: 7},
				{StartByte: 5, EndByte: 10},
				{StartByte: 5, EndByte: 12},
			},
		},
		{
			name:     "Empty input",
			input:    []ByteRange{},
			expected: []ByteRange{},
		},
	}

	for _, tt := range tests {
		test := func(t *testing.T) {
			sortByteRanges(tt.input)
			for i := range tt.input {
				if tt.input[i] != tt.expected[i] {
					t.Errorf("Mismatch at index %d: got %v, want %v", i, tt.input[i], tt.expected[i])
				}
			}
		}

		t.Run(tt.name, test)
	}
}

func TestByteRange_merge(t *testing.T) {
	tests := []struct {
		name     string
		input    []ByteRange
		expected []ByteRange
	}{
		{
			name:     "Empty input",
			input:    []ByteRange{},
			expected: nil,
		},
		{
			name: "No ranges intersect",
			input: []ByteRange{
				{StartByte: 0, EndByte: 2},
				{StartByte: 3, EndByte: 5},
				{StartByte: 6, EndByte: 8},
			},
			expected: []ByteRange{
				{StartByte: 0, EndByte: 2},
				{StartByte: 3, EndByte: 5},
				{StartByte: 6, EndByte: 8},
			},
		},
		{
			name: "All ranges intersect",
			input: []ByteRange{
				{StartByte: 1, EndByte: 5},
				{StartByte: 3, EndByte: 7},
				{StartByte: 6, EndByte: 10},
			},
			expected: []ByteRange{
				{StartByte: 1, EndByte: 10},
			},
		},
		{
			name: "Mix of intersecting and disjoint ranges",
			input: []ByteRange{
				{StartByte: 0, EndByte: 3},
				{StartByte: 2, EndByte: 5},
				{StartByte: 10, EndByte: 15},
				{StartByte: 14, EndByte: 20},
			},
			expected: []ByteRange{
				{StartByte: 0, EndByte: 5},
				{StartByte: 10, EndByte: 20},
			},
		},
		{
			name: "Intersecting touching ranges",
			input: []ByteRange{
				{StartByte: 0, EndByte: 5},
				{StartByte: 5, EndByte: 10},
			},
			expected: []ByteRange{
				{StartByte: 0, EndByte: 10},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Merge(tt.input)
			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d merged ranges, got %d", len(tt.expected), len(result))
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("at index %d: got %v, expected %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}
