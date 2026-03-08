package symbols

type Symbol struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	FilePath    string `json:"file_path"`
	WorkspaceID string `json:"workspace_id"`
	StartLine   int    `json:"start_line"`
	EndLine     int    `json:"end_line"`
	Metadata    string `json:"metadata"`
}
