package lmtd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExtractDifficulty(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		tags           []string
		wantDifficulty Difficulty
	}{
		"got correctly": {
			tags: []string{
				"author:task4233",
				"beginner",
			},
			wantDifficulty: "beginner",
		},
		"got correctly when difficulty is empty": {
			tags: []string{
				"author:task4233",
			},
			wantDifficulty: "",
		},
		"got correctly when difficulty is typoed": {
			tags: []string{
				"author:task4233",
				"beginer", // typoed
			},
			wantDifficulty: "",
		},
	}

	lmtd := LMTd{}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			gotDiff := lmtd.extractDifficulty(tt.tags)
			if gotDiff != tt.wantDifficulty {
				t.Fatalf("failed extractDifficulty: want: %s, got: %s", tt.wantDifficulty, gotDiff)
			}
		})
	}
}

func TestExtractInfo(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		targetPath string
		wantInfo   *ProblemInfo
		wantError  bool
	}{
		"got correctly": {
			targetPath: "./testdata/normal",
			wantInfo: &ProblemInfo{
				Name:       "test",
				Difficulty: "beginner",
				Order:      55,
				Points:     100,
			},
			wantError: false,
		},
	}

	lmtd := LMTd{}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			gotInfo, err := lmtd.extractInfo(tt.targetPath)
			if (err != nil) != tt.wantError {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.wantInfo, gotInfo); diff != "" {
				t.Errorf("extractInfo (-want +got) =\n%s\n", diff)
			}
		})
	}

}

func TestGenrateMarkdown(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		genreInfo    GenreInfo
		wantMarkdown string
		wantError    bool
	}{
		"generated correctly": {
			genreInfo: GenreInfo{
				Name: "web",
				ProblemInfos: []ProblemInfo{
					{
						Name:       "test",
						Difficulty: "beginner",
						Order:      1,
						Points:     100,
					},
					{
						Name:       "test2",
						Difficulty: "easy",
						Order:      2,
						Points:     200,
					},
				},
			},
			wantMarkdown: `## web

|問題名|難易度|order|points|
|:-:|:-:|:-:|:-:|
|test|beginner|1|100|
|test2|easy|2|200|
`,
			wantError: false,
		},
	}

	lmtd := LMTd{}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			gotMarkdown, err := lmtd.generateMarkdown(tt.genreInfo)
			if (err != nil) != tt.wantError {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.wantMarkdown, gotMarkdown); diff != "" {
				t.Errorf("generateMarkdown (-want +got) =\n%s\n", diff)
			}
		})
	}
}
