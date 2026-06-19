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

-- name: UpdateAutohaus :execrows
UPDATE autohaus
SET name             = $1,
    username         = $2,
    email            = $3,
    anzahl_fahrzeuge = $4,
    gruendungsdatum  = $5,
    homepage         = $6,
    telefonnummer    = $7,
    version          = version + 1,
    aktualisiert     = NOW()
WHERE id      = $8
  AND version = $9;

-- name: UpdateAdresse :exec
UPDATE adresse
SET plz  = $1,
    ort  = $2,
    land = $3
WHERE autohaus_id = $4;

-- name: DeleteAutohaus :execrows
DELETE FROM autohaus WHERE id = $1;
