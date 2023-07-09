package models

type Member struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Role     string   `json:"role,omitempty"`
	Duration int      `json:"duration,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}
