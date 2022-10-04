package lmtd

import "fmt"

type Difficulty string

var difficulties = []Difficulty{
	"beginner",
	"easy",
	"medium",
	"hard",
}

// TODO: ChallengeとFlagのownerChallengeUniqueKeyが一致していることを確認する
type Challenge struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Spec       struct {
		Name  string   `yaml:"name"`
		Order int      `yaml:"order"`
		Tags  []string `yaml:"tags"`
	} `yaml:"spec"`
	// TODO: https://github.com/SECCON/beginners-lecture-ctf/blob/main/misc/addition_master/.lmtd/challenge.yaml
}

type Flag struct {
	ApiVersion string `yaml:"apiVersion"`
	Spec       struct {
		Point int `yaml:"point"`
	} `yaml:"spec"`
	// TOOD: https://github.com/SECCON/beginners-lecture-ctf/blob/main/misc/addition_master/.lmtd/flag.yaml
}

// ProblemInfo は各問題の情報を保持
type ProblemInfo struct {
	Name       string
	Difficulty Difficulty
	Order      int
	Points     int
}

func (p *ProblemInfo) Validate() error {
	if p.Points/10 != p.Order {
		return fmt.Errorf("points(%d) should be one-tenth of order(%d)", p.Points, p.Order)
	}

	return nil
}

// GenreInfo は各ジャンルの情報を保持
type GenreInfo struct {
	Name         string
	ProblemInfos []ProblemInfo
}
