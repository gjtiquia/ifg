package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gjtiquia/ifg/internal/config"
	"github.com/gjtiquia/ifg/internal/search"
)

type Server struct {
	Port    int
	Entries []config.Entry
}

func NewServer(port int) (*Server, error) {
	configDir := config.GetConfigDir()
	entries, err := config.LoadConfig(configDir)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}
	return &Server{
		Port:    port,
		Entries: entries,
	}, nil
}

func (s *Server) Start() error {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/search", s.handleSearch)
	addr := fmt.Sprintf(":%d", s.Port)
	return http.ListenAndServe(addr, nil)
}

func isCLIRequest(r *http.Request) bool {
	ua := r.Header.Get("User-Agent")
	return strings.Contains(ua, "curl/") ||
		strings.Contains(ua, "Wget/") ||
		strings.Contains(ua, "HTTPie/")
}

func formatEntriesText(entries []config.Entry) string {
	var b strings.Builder
	for _, e := range entries {
		if e.Title != "" && e.Title != e.Command {
			b.WriteString("# ")
			b.WriteString(e.Title)
			b.WriteString("\n")
		}
		for _, desc := range e.Description {
			b.WriteString("# ")
			b.WriteString(desc)
			b.WriteString("\n")
		}
		b.WriteString(e.Command)
		b.WriteString("\n\n")
	}
	return b.String()
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path != "/" {
		path = strings.TrimPrefix(path, "/")
	} else {
		path = ""
	}

	if isCLIRequest(r) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if path == "" {
			w.Write([]byte(formatEntriesText(s.Entries)))
			return
		}
		query := strings.ReplaceAll(path, "-", " ")
		filtered := search.Match(s.Entries, query)
		if len(filtered) == 0 {
			fmt.Fprintf(w, "# No commands found for: %s\n", path)
			return
		}
		w.Write([]byte(formatEntriesText(filtered)))
		return
	}

	if path != "" {
		query := strings.ReplaceAll(path, "-", " ")
		http.Redirect(w, r, "/?q="+query, http.StatusFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}

type stateResponse struct {
	SelectedIdx int         `json:"selectedIdx"`
	Entries     []entryJSON `json:"entries"`
}

type entryJSON struct {
	Title       string   `json:"title"`
	Description []string `json:"description"`
	Command     string   `json:"command"`
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	stateJSON := r.URL.Query().Get("state")

	selectedIdx := 0
	if stateJSON != "" {
		var state stateResponse
		if err := json.Unmarshal([]byte(stateJSON), &state); err == nil {
			selectedIdx = state.SelectedIdx
		}
	}

	filtered := search.Match(s.Entries, query)

	if selectedIdx < 0 {
		selectedIdx = 0
	}
	if selectedIdx >= len(filtered) {
		selectedIdx = len(filtered) - 1
		if selectedIdx < 0 {
			selectedIdx = 0
		}
	}

	if isCLIRequest(r) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if len(filtered) == 0 {
			fmt.Fprintf(w, "# No commands found for: %s\n", query)
			return
		}
		w.Write([]byte(formatEntriesText(filtered)))
		return
	}

	entriesJSON := make([]entryJSON, len(filtered))
	for i, e := range filtered {
		desc := e.Description
		if desc == nil {
			desc = []string{}
		}
		entriesJSON[i] = entryJSON{
			Title:       e.Title,
			Description: desc,
			Command:     e.Command,
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintf(w, `<div class="entry-list-inner">`)

	if len(filtered) == 0 {
		fmt.Fprint(w, `<div class="empty"><div class="empty-title">No commands found</div></div>`)
	} else {
		for i, entry := range filtered {
			selectedClass := ""
			if i == selectedIdx {
				selectedClass = " selected"
			}
			fmt.Fprintf(w, `<div class="entry%s" data-idx="%d">`, selectedClass, i)
			if entry.Title != "" && entry.Title != entry.Command {
				fmt.Fprintf(w, `<div class="entry-title">%s</div>`, escapeHTML(entry.Title))
			}
			if len(entry.Description) > 0 {
				fmt.Fprintf(w, `<div class="entry-desc">%s</div>`, escapeHTML(joinLines(entry.Description)))
			}
			fmt.Fprintf(w, `<div class="entry-cmd">%s</div>`, escapeHTML(entry.Command))
			fmt.Fprint(w, `</div>`)
		}
	}

	fmt.Fprint(w, `</div>`)

	stateOut := stateResponse{
		SelectedIdx: selectedIdx,
		Entries:     entriesJSON,
	}
	stateBytes, _ := json.Marshal(stateOut)
	fmt.Fprintf(w, `<input type="hidden" name="state" value="%s">`, escapeHTML(string(stateBytes)))
}

func escapeHTML(s string) string {
	s = replaceAll(s, "&", "&amp;")
	s = replaceAll(s, "<", "&lt;")
	s = replaceAll(s, ">", "&gt;")
	s = replaceAll(s, `"`, "&quot;")
	s = replaceAll(s, "'", "&#39;")
	return s
}

func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); i++ {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old) - 1
		} else {
			result += string(s[i])
		}
	}
	return result
}

func joinLines(lines []string) string {
	result := ""
	for i, line := range lines {
		if i > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}
