# [i] [f]or[g]ot

for when you are trying to rmb that command

an interactive cli with shell integration support

dead-simple config format with fuzzy search

## installation

```bash
go install github.com/gjtiquia/ifg
```

### shell integration

```bash
# add to `~/.bashrc` or `~/.zshrc`
eval "$(ifg --sh)"
```

this adds the command to the history instead of just printing it out,
which is useful as you can access the command by simply pressing UP

## usage

```bash
ifg
```

## config

### config location
- if `XDG_CONFIG_HOME` is set: `$XDG_CONFIG_HOME/ifg/config.sh`
- if `XDG_CONFIG_HOME` is not set: `~/.ifg/config.sh`

if no config exists, a default one will be created on first running `ifg`

### config format

```bash
# an optional title
# an optional description
# as many lines as you want
echo "the command you want to remember"

# another title
echo "another command"

echo "titles and descriptions are overrated"
```

## development

```bash
# runs main.go
go run .

# builds binary at project root named `ifg`
go build .

# runs tests
go test ./...
```

### testing shell integration

```bash
# builds binary at project root named `ifg`
go build .

# adds project root to PATH temporarily
export PATH="$PWD:$PATH"

# loads wrapper to current shell
eval "$(ifg --sh)"
```

## license

MIT

