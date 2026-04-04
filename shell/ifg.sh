# ifg - [i] [f]or[g]ot
# Add to ~/.bashrc OR ~/.zshrc:
#   eval "$(ifg --sh)"

ifg() {
	if [[ $# -gt 0 ]]; then
		command ifg "$@"
		return $?
	fi

	local cmd=$(command ifg)
	if [[ -n "$cmd" ]]; then
		history -s "$cmd"
		echo "# ifg - [i] [f]or[g]ot"
		echo ""
		echo "  $cmd"
		echo ""
		echo "[press UP to access from history]"
		echo ""
	fi
}
