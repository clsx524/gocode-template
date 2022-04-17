package model

type Company struct {
	ID   string `json:"id,required"`
	Name string `json:"name,omitempty"`
}
