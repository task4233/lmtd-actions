{{ range $genreInfo := .GenreInfos }}## {{ $genreInfo.Name }}

|問題名|難易度|order|points|
|:-:|:-:|:-:|:-:|{{ range $info := $genreInfo.ProblemInfos }}
|{{ $info.Name }}|{{ $info.Difficulty }}|{{ $info.Order }}|{{ $info.Points }}|{{ end }}
{{ end }}