package lmtd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type LMTd struct {
	Stdout io.Writer
	Stderr io.Writer
}

func (l LMTd) Run(ctx context.Context, version string, args []string) error {
	if len(args) == 0 {
		return errors.New("lmtd-cheker <target directory path>")
	}
	targetDir := args[0]
	genres, err := os.ReadDir(targetDir)
	if err != nil {
		return err
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
			// ReadDirの実行だけで後続の処理を落としたくないので、
			// stderrでロギングすることでエラーハンドリングしたことにする
			fmt.Fprintf(l.Stderr, "failed os.ReadDir: genre: %s, err: %s\n", genre.Name(), err.Error())
			continue
		}

		problemInfos := make([]ProblemInfo, 0, len(problems))
		for _, problem := range problems {
			// genre用directory直下にある.lmtdはcategory.yamlしか含まれないので無視する
			if problem.Name() == ".lmtd" || problem.Name() == "README.md" {
				continue
			}

			info, err := l.extractInfo(filepath.Join(genreDir, problem.Name()))
			if err != nil {
				// 情報収集しただけで後続の処理を落としたくないので、
				// stderrでロギングすることでエラーハンドリングしたことにする
				fmt.Fprintf(l.Stderr, "failed extractInfo: genre: %s, err: %s\n", genre.Name(), err.Error())
				continue
			}

			problemInfos = append(problemInfos, *info)
		}

		// 問題が1問もないと判断してスルーする
		if len(problemInfos) == 0 {
			continue
		}

		// TODO: infosからREADME.mdを生成する
		genreInfo := GenreInfo{
			Name:         genre.Name(),
			ProblemInfos: problemInfos,
		}
		markdown, err := l.generateMarkdown(genreInfo)
		if err != nil {
			return fmt.Errorf("failed extractInfo: %s", err.Error())
		}

		outputPath := filepath.Join(genreDir, "README.md")
		if err = l.appendInfo(markdown, outputPath); err != nil {
			return fmt.Errorf("failed appendInfo: %s", err.Error())
		}
	}

	return nil
}

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

func (l LMTd) extractInfo(targetPath string) (*ProblemInfo, error) {
	lmtdDir := filepath.Join(targetPath, ".lmtd")

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
	var difficulty Difficulty = l.extractDifficulty(chalData.Spec.Tags)

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

	p := &ProblemInfo{
		Name:       chalData.Spec.Name,
		Difficulty: difficulty,
		Order:      chalData.Spec.Order,
		Points:     flagData.Spec.Point,
	}
	if err := p.Validate(); err != nil {
		// TODO: OrderとPointsの対応関係が取れていない場合は、Pointsを元にOrderを上書きする
		// 上書きするとフォーマットが変わりそうで嫌な気もする
		// ひとまずロギングで対処
		fmt.Fprintf(l.Stderr, "order of %s/.lmtd/challenge.yaml should be %d, but got %d\n", targetPath, p.Points/10, p.Order)
		p.Order = p.Points / 10

		// buf, err := yaml.Marshal(p)
		// if err != nil {
		// 	return p, fmt.Errorf("failed yaml.Marshal: data(%v): %s", *p, err.Error())
		// }
	}

	return p, nil
}

// O(len(tags) * len(difficulties))だが、
// 前者も後者も高々5程度に収まるはずなので実質定数倍
func (l LMTd) extractDifficulty(tags []string) Difficulty {
	for _, tag := range tags {
		// spec.Tagsからdifficultyのみを抜き出す
		for _, diff := range difficulties {
			if tag == string(diff) {
				return diff
			}
		}
	}

	// difficultyが空の可能性があるのでerrorは返さない
	return ""
}

func (l LMTd) generateMarkdown(genreInfo GenreInfo) (string, error) {
	// 追加で表示したい情報が増える可能性がありそうなので、複数ファイル指定可能にしておく
	tpl, err := template.ParseFiles([]string{"templates/genreInfo.md.tpl"}...)
	if err != nil {
		return "", err
	}

	writer := &bytes.Buffer{}
	err = tpl.Execute(writer, struct {
		GenreInfo GenreInfo
	}{
		GenreInfo: genreInfo,
	})
	if err != nil {
		return "", err
	}

	return writer.String(), nil
}

func (l LMTd) appendInfo(markdown string, outputPath string) error {
	// ファイルが存在しない時に新規作成したくない場合は、os.O_WRONLY|os.O_APPENDにする
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprint(file, markdown)
	return err
}
