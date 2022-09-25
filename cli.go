package lmtd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type CLI struct {
	Stdout io.Writer
	Stderr io.Writer
}

func (c CLI) Run(ctx context.Context, version string, args []string) error {
	targetDir := "/Users/deadbeef/work/beginners-lecture-ctf"
	genres, err := os.ReadDir(targetDir)
	if err != nil {
		log.Fatalln(err)
	}

	for _, genre := range genres {
		if !genre.IsDir() {
			continue
		}
		// .から始まるディレクトリは関係ないのでスルーする
		if strings.HasPrefix(genre.Name(), ".") {
			continue
		}

		genreDir := filepath.Join(targetDir, genre.Name())
		problems, err := os.ReadDir(genreDir)
		if err != nil {
			log.Fatalln(err)
		}

		infos := []ProblemInfo{}
		for _, problem := range problems {
			// genre用directory直下にある.lmtdはcategory.yamlしか含まれないので無視する
			if problem.Name() == ".lmtd" {
				continue
			}

			info, err := extractInfo(filepath.Join(genreDir, problem.Name()))
			if err != nil {
				log.Fatalln(err)
			}

			infos = append(infos, *info)
		}

		// 問題が1問もないと判断してスルーする
		if len(infos) == 0 {
			continue
		}

		fmt.Println(genre.Name(), "\tinfos:\t", infos)

		// TODO: infosからREADME.mdを生成する
	}

	return nil
}

type Difficulty string

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
	ProblemName string
	Difficulty  Difficulty
	Order       int
	Point       int
}

func extractInfo(target string) (*ProblemInfo, error) {
	lmtdDir := filepath.Join(target, ".lmtd")

	// challenge.yaml -> order, tag(difficulty)
	buf, err := os.ReadFile(filepath.Join(lmtdDir, "challenge.yaml"))
	if err != nil {
		return nil, err
	}
	chalData := Challenge{}
	err = yaml.Unmarshal(buf, &chalData)
	if err != nil {
		return nil, err
	}
	// TODO: spec.Tagsからdifficultyのみを抜き出す
	var difficulty Difficulty = ""

	// flag.yaml -> point
	buf, err = os.ReadFile(filepath.Join(lmtdDir, "flag.yaml"))
	if err != nil {
		return nil, err
	}
	flagData := Flag{}
	err = yaml.Unmarshal(buf, &flagData)
	if err != nil {
		return nil, err
	}

	return &ProblemInfo{
		ProblemName: chalData.Spec.Name,
		Difficulty:  difficulty,
		Order:       chalData.Spec.Order,
		Point:       flagData.Spec.Point,
	}, nil
}
