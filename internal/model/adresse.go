package model

type Adresse struct {
	ID         uint   `json:"id"`
	PLZ        string `json:"plz"  validate:"required,len=5,numeric"`
	Ort        string `json:"ort"  validate:"required,min=2,max=100"`
	Land       string `json:"land" validate:"required,min=2,max=100"`
	AutohausID uint   `json:"autohausId"`
}
