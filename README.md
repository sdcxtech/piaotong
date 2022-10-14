# piaotong

[![GoDoc](https://godoc.org/github.com/sdcxtech/piaotong?status.svg)](https://godoc.org/github.com/sdcxtech/piaotong)


Golang SDK for 票通


## Development

### Install [pre-commit](https://pre-commit.com/)

```sh
# if you havn't install pre-commit on your mac
brew install pre-commit

# install pre-commit into your git hooks
pre-commit install
```

### Install [commitlint](https://commitlint.js.org/)

```sh
npm install -g commitlint
npm install -g @commitlint/config-conventional
```

### Install [golangci-lint](https://golangci-lint.run/)

```sh
brew install golangci/tap/golangci-lint
brew upgrade golangci/tap/golangci-lint
```

Add below json snappet to your vscode settings:

```json
    "go.lintTool": "golangci-lint",
    "go.lintFlags": [
        "--fast"
      ]
```