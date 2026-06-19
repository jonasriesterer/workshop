package model

// Auto repräsentiert ein Fahrzeug, das einem Autohaus zugeordnet ist.
type Auto struct {
	ID          uint   `json:"id"`
	Kennzeichen string `json:"kennzeichen" validate:"required,min=3,max=10"`
	Marke       string `json:"marke"       validate:"required,min=2,max=50"`
	Modell      string `json:"modell"      validate:"required,min=1,max=50"`
	Baujahr     int    `json:"baujahr"     validate:"required,min=1886,max=2100"`
	AutohausID  uint   `json:"autohausId"`
}
