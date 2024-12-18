package ascii

import (
	"strings"
	"testing"

	"github.com/github/gh-skyline/types"
)

func TestGenerateASCII(t *testing.T) {
	tests := []struct {
		name          string
		grid          [][]types.ContributionDay
		user          string
		year          int
		includeHeader bool
		wantErr       bool
	}{
		{
			name:          "empty grid",
			grid:          [][]types.ContributionDay{},
			user:          "testuser",
			year:          2023,
			includeHeader: false,
			wantErr:       true,
		},
		{
			name:          "valid grid",
			grid:          makeTestGrid(3, 7),
			user:          "testuser",
			year:          2023,
			includeHeader: false,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateASCII(tt.grid, tt.user, tt.year, tt.includeHeader)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateASCII() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Existing validation code...
				if !strings.Contains(result, "testuser") {
					t.Error("Generated ASCII should contain username")
				}
				if !strings.Contains(result, "2023") {
					t.Error("Generated ASCII should contain year")
				}
				if !strings.Contains(result, string(EmptyBlock)) {
					t.Error("Generated ASCII should contain empty blocks")
				}
			}
		})
	}
}

// Helper function to create test grid
func makeTestGrid(weeks, days int) [][]types.ContributionDay {
	grid := make([][]types.ContributionDay, weeks)
	for i := range grid {
		grid[i] = make([]types.ContributionDay, days)
		for j := range grid[i] {
			grid[i][j] = types.ContributionDay{ContributionCount: i * j}
		}
	}
	return grid
}

func TestGetBlock(t *testing.T) {
	tests := []struct {
		name         string
		normalized   float64
		dayIdx       int
		nonZeroIdx   int
		expectedRune rune
	}{
		{"empty block", 0.0, 0, 1, EmptyBlock},
		{"single low block", 0.2, 0, 1, FoundationLow},
		{"single medium block", 0.5, 0, 1, FoundationMed},
		{"single high block", 0.8, 0, 1, FoundationHigh},
		{"foundation low", 0.2, 0, 2, FoundationLow},
		{"middle high", 0.8, 1, 3, MiddleHigh},
		{"top medium", 0.5, 2, 3, TopMed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getBlock(tt.normalized, tt.dayIdx, tt.nonZeroIdx)
			if result != tt.expectedRune {
				t.Errorf("getBlock(%f, %d, %d) = %c, want %c",
					tt.normalized, tt.dayIdx, tt.nonZeroIdx,
					result, tt.expectedRune)
			}
		})
	}
}
