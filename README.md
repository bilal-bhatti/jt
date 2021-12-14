# jt

[![Go Report Card](https://goreportcard.com/badge/github.com/bilal-bhatti/jt)](https://goreportcard.com/report/github.com/bilal-bhatti/jt)


cli json tools

jt is a cli wrapper around `jq` and `json-path` implementations for transforming json documents.

```
github.com/itchyny/gojq
github.com/ohler55/ojg
```

## install
``` sh
brew tap bilal-bhatti/homebrew-taps
brew install jt
```
or
``` sh
go install github.com/bilal-bhatti/jt/cmd/jt@latest
```

## commands
get started

``` sh
jt commands
jt help template
jt help apply
```

### template 
generate a jq based transformation template

``` sh
echo '{"x":"y"}' | jt
cat examples/i.json | jt
```

### apply
apply template to input json

``` sh
cat examples/i.json | jt apply -t examples/t.json
jt apply -t examples/t.json -i examples/i.json
jt apply -i examples/servers.json -t examples/template.json
```

