# lmtd-actions
lmtd-actionsはlmtd用repositoryにおける各ジャンルの情報をまとめてmarkdownに書き出すツールです。

## Installation
```bash
go install github.com/task4233/lmtd-actions/cmd/lmtd-actions@latest
```

## Usage
```bash
lmtd-actions <target dir>
```

## Example
2022_beginners_ctf

```bash
$ lmtd-actions 2022_beginnersctf_ctf
$ cat web/README.md
## web

|問題名|難易度|order|points|
|:-:|:-:|:-:|:-:|
|gallery|easy|2|500|
|Ironhand|medium|3|500|
|serial|medium|3|500|
|textex|easy|2|500|
|Util|beginner|1|500|
```
