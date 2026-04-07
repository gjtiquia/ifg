# [i] [f]or[g]ot

```
ifg - [i] [f]or[g]ot

type to search: magick

---

  # convert from jpg to png
  # supports many other formats too
  magick rose.jpg rose.png

> # basic info
> # like the image size
> magick identify image.png

  # resize
  # can be used for basic image upscaling
  magick input.png -resize 200% output.png
```

## intro

for when you are trying to rmb that command

curate your own personal library of commands

dead-simple config format with fuzzy search

`fzf`, `telescope.nvim`, `⌘K`, `Ctrl+P` inspired UX

## installation

### go install (recommended)

prerequisite: [Go](https://go.dev/doc/install)

```bash
go install github.com/gjtiquia/ifg@latest
```

<details>
    <summary><code>Command 'ifg' not found</code></summary>

```bash
# make sure this is in your .bashrc / .zshrc
export PATH=$PATH:$HOME/go/bin
```
</details>

<details>
    <summary>other useful commands:</summary>

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
</details>

### download binary

skip this if you have downloaded `ifg` via `go install`

```bash
# TODO
```

### shell integration

this adds the command to the history instead of just printing it out,
which is useful as you can access the command by simply pressing UP

```bash
# for Linux and MacOS:
# add to your `~/.bashrc` or `~/.zshrc`
eval "$(ifg --sh)"

# for Windows:
# TODO
```

## usage

### interactive cli

```bash
ifg

# then type to fuzzy search
# select the command with arrrow keys and press enter

# vim keys are supported as well
# escape and navigate with j/k. go back to insert mode with i/I/a/A
```

### web server (experimental)

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

can be a dead simple single file config

```
~/.ifg/
└── config.sh
```

or a collection of files, organized to your own liking.

all `*.sh` files in the config directory are read, sorted alphabetically by path:

```
~/.ifg/
├── git.sh
├── docker.sh
├── personal/
│   └── scripts.sh
└── work/
    ├── 01-ssh.sh
    └── 02-deploy.sh
```

number prefixes are optional.
use them for custom ordering.

subdirectories are supported.

feel free to check out [my personal config](https://github.com/gjtiquia/.ifg) for more examples

### config format

its just bash comments and one-liners. entries are seperated by empty lines

```bash
# an optional title
# an optional description
# as many lines as you want
echo "the command you want to remember"

# another title
echo "another command"

echo "titles and descriptions are overrated"
```

check out the default [`/shell/config.sh`](./shell/config.sh) for a more concrete example

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

if you do not want to "pollute" the global install, you may test it this way

```bash
# builds binary at project root named `ifg`
go build .

# adds project root to PATH temporarily
export PATH="$PWD:$PATH"

# loads wrapper to current shell
eval "$(ifg --sh)"
```

## ai-usage disclosure

ai is heavily used for generating code in this project

this project aims to build something i need in a short amount of time, but also to
- explore a viable dev workflow with ai
- explore [opencode](https://opencode.ai/)
- explore the capabilities of the GLM 5 model via [opencode go](https://opencode.ai/go) 

some thoughts
- its not as "fun" as writing by hand, but i have to admit its much faster
- golang is a really great language to vibecode with
  - idiomatic go is pretty "standard" over the past decade
  - stdlib works great for a lot of things
  - strong typing, compile time errors
  - dead simple unit testing
  - these feedback loops increase the quality of code that ai generates
- cli (or text-based) tools are great for to vibecode for
  - much more predictable UI design, and easier to enforce a certain style
  - very easy to do "integration tests" via unit tests
- while claude and gpt are amazing, GLM 5 is very much good enough for me

## license

MIT
