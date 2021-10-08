# dendrite-lint

Simple linter to check if all defined endpoints are used on the server side. (inspired by [this tutorial](https://disaev.me/p/writing-useful-go-analysis-linter))

Only checks the ```inthttp``` package at the moment.

## Usage

Clone this repository:

```bash
cd dendrite-lint
go install ./cmd/...

cd $yourDendritePath
dendrite-lint ./...
# or
go vet -vettool=$(which dendrite-lint) ./...
```

It should display something along

```bash
userapi/inthttp/client.go:52:2: declared 'QueryKeyBackupPath' endpoint, but not used in internal api server
userapi/inthttp/client.go:40:2: declared 'PerformKeyBackupPath' endpoint, but not used in internal api server
```

if there's an endpoint defined but not used.