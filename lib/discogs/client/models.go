package client

// Identifier ...
type Identifier struct {
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	Value       string `json:"value"`
}

type Artist struct {
	Anv         string `json:"anv"`
	ID          int    `json:"id"`
	Join        string `json:"join"`
	Name        string `json:"name"`
	ResourceURL string `json:"resource_url"`
	Role        string `json:"role"`
	Tracks      string `json:"tracks"`
}

type Track struct {
	Duration string `json:"duration"`
	Position string `json:"position"`
	Title    string `json:"title"`
	Type     string `json:"type_"`
}
