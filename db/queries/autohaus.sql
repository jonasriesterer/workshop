-- name: GetAutohausById :one
SELECT id, version, name, username, email, anzahl_fahrzeuge, gruendungsdatum, homepage, telefonnummer, erzeugt, aktualisiert
FROM autohaus
WHERE id = $1;

-- name: GetAdresseByAutohausId :one
SELECT id, plz, ort, land, autohaus_id
FROM adresse
WHERE autohaus_id = $1;

-- name: GetAutosByAutohausId :many
SELECT id, kennzeichen, marke, modell, baujahr, autohaus_id
FROM auto
WHERE autohaus_id = $1;

-- name: CreateAutohaus :one
INSERT INTO autohaus (version, name, username, email, anzahl_fahrzeuge, gruendungsdatum, homepage, telefonnummer, erzeugt, aktualisiert)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
RETURNING id;

-- name: CreateAdresse :exec
INSERT INTO adresse (plz, ort, land, autohaus_id)
VALUES ($1, $2, $3, $4);

-- name: CreateAuto :exec
INSERT INTO auto (kennzeichen, marke, modell, baujahr, autohaus_id)
VALUES ($1, $2, $3, $4, $5);
