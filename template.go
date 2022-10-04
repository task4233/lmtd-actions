package lmtd

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed templates/genreInfo.md.tpl
var genreInfoTemplate string

func (l LMTd) generateMarkdown(genreInfo GenreInfo) (string, error) {
	// templateは引数で渡せるようにした方が綺麗なので、複数ファイルのembedをするようになったらそうする
	tpl, err := template.New("genreInfo").Parse(genreInfoTemplate)
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
