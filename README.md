# [i] [f]or[g]ot cli

a cli tool to help you remember what is the command to do what you are trying to do

deadsimple and dumb config as well

would be great to have a quick interactive cli rather than remember dozens of aliases

bringing the ux of fzf's ctrl+r and nvim telescope to the shell

## usage

the ux should be very simple to fzf's ctrl+r behavior
- type `ifg` and enter
- enters an interactive session where users can type keywords
- if the user doesnt type any keyword, the user can up/down arrow or "j" "k" to select up and down from a list of all commands
- if the user types keywords, the cli matches the command / command title / command description and gives a list where user can also select up/down j/k
- when user presses enter, the command will appear, just like ctrl+r in fzf, for the user to modify it a bit, before pressing enter to execute

## tech stack

go

## installation

```bash

# prerequisites - have go installed
go install github.com/gjtiquia/ifg
```

## config

configuration lives in ~/.ifg/config.sh

the `.sh` extension is simply for convenience of IDE highlighting

an example configuration

```bash
# copy to clipboard (MacOS)
# this command copies to clipboard
# $ echo "hi" | pbcopy
pbcopy

# paste from clipboard (MacOS)
# $ pbpaste >> file.txt
pbpaste

# copy to clipboard (Linux)
# $ echo "hi" | xclip sel -clip
xclip sel -clip
```

super straightforward config
- separated by empty newline
- first comment line (if any) is title
- second comment line and subsequent lines are descriptions (if any)
- last line is the command
- the keyword search should be simper simple too, searches thru each "entry block" for matches
- keyword search should be "order-agnostic", eg. typing "macos copy" should match "copy to clipboard (MacOS)"


