package models

// Company defines the company class
type Company struct {
	ID   string `json:"id,required"`
	Name string `json:"name,omitempty"`
}
