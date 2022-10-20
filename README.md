# piaotong

[![GoDoc](https://pkg.go.dev/badge/github.com/sdcxtech/piaotong)](https://pkg.go.dev/github.com/sdcxtech/piaotong)

Golang SDK for 票通

## Usage
Open blue invoice

```go
c := piaotong.New(piaotong.Config{})
resp, _ := c.OpenBlueInvoice(ctx, &piaotong.OpenBlueInvoiceRequest{})
fmt.Println(resp)
```


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