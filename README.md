# [i] [f]or[g]ot

for when you are trying to rmb that command

an interactive cli with shell integration support

dead-simple config format with fuzzy search

## installation

```bash
go install github.com/gjtiquia/ifg@latest
```

<details>
    <summary>Command 'ifg' not found</summary>

```bash
# make sure this is in your .bashrc / .zshrc
export PATH=$PATH:$HOME/go/bin
```
</details>

other useful commands:

```bash
# checks what is the latest available version on go proxy cache
go list -m github.com/gjtiquia/ifg@latest

# checks what is the latest version directly from GitHub
GOPROXY=direct go list -m github.com/gjtiquia/ifg@latest

# installs latest version directly from GitHub
GOPROXY=direct go install github.com/gjtiquia/ifg@latest

# installs binary at current directory instead of a global install
GOBIN=$(pwd) go install github.com/gjtiquia/ifg@latest
```

### shell integration

```bash
# add to `~/.bashrc` or `~/.zshrc`
eval "$(ifg --sh)"
```

this adds the command to the history instead of just printing it out,
which is useful as you can access the command by simply pressing UP

## usage

### interactive cli

```bash
ifg

# then type to fuzzy search
# select the command with arrrow keys and press enter

# vim keys are supported as well
# escape and navigate with j/k. enter back to insert mode with i/I/a/A
```

### web server

you can also serve `ifg` on the web,
accessible via browser or `curl`,
returning html or plain text respectively,
its like having a cli cheatsheet readily available whenever you need

```bash
# serve web server at port 5432
ifg web --port 5432

# list all entries
curl -L your-ifg-domain.com

# fuzzy search entry
curl -L your-ifg-domain.com/your-query
```

## config

### config directory location
- if `XDG_CONFIG_HOME` is set: `$XDG_CONFIG_HOME/ifg/`
- if `XDG_CONFIG_HOME` is not set: `~/.ifg/`

if no config directory exists,
the config directory will be automatically created on first running `ifg`,
with a default config called `config.sh`

### config directory structure

all `*.sh` files in the config directory are read, sorted alphabetically by path:

```
~/.ifg/
├── 01-git.sh
├── 02-docker.sh
├── personal/
│   └── scripts.sh
└── work/
    ├── 01-ssh.sh
    └── 02-deploy.sh
```

number prefixes are optional.
use them for custom ordering.

subdirectories are supported.

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

check out the default `config.sh` for a more concrete example

## development

```bash
# runs tests
go test ./...

# runs main.go
go run .

# builds binary at project root named `ifg`
go build .

# builds binary at ~/go/bin for global usage
go install .
```

### testing shell integration

if you do not want to "polute" the global install, you may test it this way

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
