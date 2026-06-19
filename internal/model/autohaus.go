package model

import "time"

type Autohaus struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"            validate:"required,min=2,max=100"`
	Username        string    `json:"username"        validate:"required,min=3,max=50,excludesall= "`
	Email           string    `json:"email"           validate:"required,email"`
	AnzahlFahrzeuge int       `json:"anzahlFahrzeuge" validate:"min=0"`
	Gruendungsdatum time.Time `json:"gruendungsdatum" validate:"required"`
	Homepage        string    `json:"homepage"        validate:"omitempty,url"`
	Telefonnr       string    `json:"telefonnr"       validate:"omitempty,e164"`
	// Version wird für optimistisches Locking verwendet
	Version      uint      `json:"version"`
	Erzeugt      time.Time `json:"erzeugt"`
	Aktualisiert time.Time `json:"aktualisiert"`

	Adresse Adresse `json:"adresse"`
	Autos   []Auto  `json:"autos"`
}

