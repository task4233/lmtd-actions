# lmtd-checker
lmtd-checkerはlmtd用repositoryにおける各ジャンルの情報をまとめてmarkdownに書き出すツールです。

## Installation
```bash
go install github.com/task4233/lmtd-checker/cmd/lmtd-checker@latest
```

## Usage
```bash
lmtd-checker <target dir>
```

## Example
2022_beginners_ctf

```bash
$ lmtd-checker 2022_beginnersctf_ctf
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
