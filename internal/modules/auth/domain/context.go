package domain

type UserContext struct {
	ID    string   `json:"id"`
	Roles []string `json:"roles"`
}
