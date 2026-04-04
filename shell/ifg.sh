# ifg - interactive command finder
# Add to ~/.bashrc OR ~/.zshrc:
#   source "$(ifg --sh)"

ifg() {
	local cmd=$(command ifg)
	if [[ -n "$cmd" ]]; then
		history -s "$cmd"
		echo "Command: $cmd"
		echo "Press UP to access from history"
	fi
}
