// Package seed befüllt die Datenbank mit Testdaten aus CSV-Dateien.
package seed

import (
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Load seeds the database from CSV files located in csvDir.
// It uses OVERRIDING SYSTEM VALUE to preserve the IDs from the CSV files
// and resets the identity sequences afterwards.
func Load(ctx context.Context, pool *pgxpool.Pool, csvDir string) error {
	slog.Info("seeding database", "dir", csvDir)

	if err := seedAutohaus(ctx, pool, csvDir); err != nil {
		return fmt.Errorf("seed autohaus: %w", err)
	}
	if err := seedAdresse(ctx, pool, csvDir); err != nil {
		return fmt.Errorf("seed adresse: %w", err)
	}
	if err := seedAuto(ctx, pool, csvDir); err != nil {
		return fmt.Errorf("seed auto: %w", err)
	}

	if err := resetSequences(ctx, pool); err != nil {
		return fmt.Errorf("reset sequences: %w", err)
	}

	slog.Info("seeding complete")
	return nil
}

func readCSV(path string) ([]map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	r := csv.NewReader(f)
	r.Comma = ';'

	headers, err := r.Read()
	if err != nil {
		return nil, err
	}

	var rows []map[string]string
	for {
		record, err := r.Read()
		if err != nil {
			break
		}
		row := make(map[string]string, len(headers))
		for i, h := range headers {
			if i < len(record) {
				row[strings.TrimSpace(h)] = strings.TrimSpace(record[i])
			}
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func seedAutohaus(ctx context.Context, pool *pgxpool.Pool, dir string) error {
	rows, err := readCSV(dir + "/autohaus.csv")
	if err != nil {
		return err
	}
	for _, r := range rows {
		id, _ := strconv.Atoi(r["id"])
		version, _ := strconv.Atoi(r["version"])
		anzahl, _ := strconv.Atoi(r["anzahl_fahrzeuge"])
		gruendung, err := time.Parse("2006-01-02", r["gruendungsdatum"])
		if err != nil {
			return fmt.Errorf("parse gruendungsdatum: %w", err)
		}
		erzeugt, _ := time.Parse("2006-01-02 15:04:05", r["erzeugt"])
		aktualisiert, _ := time.Parse("2006-01-02 15:04:05", r["aktualisiert"])

		_, err = pool.Exec(ctx,
			`INSERT INTO autohaus (id, version, name, username, email, anzahl_fahrzeuge, gruendungsdatum, homepage, telefonnummer, erzeugt, aktualisiert)
			 OVERRIDING SYSTEM VALUE
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
			id, version, r["name"], r["username"], r["email"],
			anzahl, gruendung, r["homepage"], r["telefonnummer"],
			erzeugt, aktualisiert,
		)
		if err != nil {
			return fmt.Errorf("insert autohaus id=%d: %w", id, err)
		}
	}
	return nil
}

func seedAdresse(ctx context.Context, pool *pgxpool.Pool, dir string) error {
	rows, err := readCSV(dir + "/adresse.csv")
	if err != nil {
		return err
	}
	for _, r := range rows {
		id, _ := strconv.Atoi(r["id"])
		autohausID, _ := strconv.Atoi(r["autohaus_id"])

		_, err = pool.Exec(ctx,
			`INSERT INTO adresse (id, plz, ort, land, autohaus_id)
			 OVERRIDING SYSTEM VALUE
			 VALUES ($1,$2,$3,$4,$5)`,
			id, r["plz"], r["ort"], r["land"], autohausID,
		)
		if err != nil {
			return fmt.Errorf("insert adresse id=%d: %w", id, err)
		}
	}
	return nil
}

func seedAuto(ctx context.Context, pool *pgxpool.Pool, dir string) error {
	rows, err := readCSV(dir + "/auto.csv")
	if err != nil {
		return err
	}
	for _, r := range rows {
		id, _ := strconv.Atoi(r["id"])
		baujahr, _ := strconv.Atoi(r["baujahr"])
		autohausID, _ := strconv.Atoi(r["autohaus_id"])

		_, err = pool.Exec(ctx,
			`INSERT INTO auto (id, kennzeichen, marke, modell, baujahr, autohaus_id)
			 OVERRIDING SYSTEM VALUE
			 VALUES ($1,$2,$3,$4,$5,$6)`,
			id, r["kennzeichen"], r["marke"], r["modell"], baujahr, autohausID,
		)
		if err != nil {
			return fmt.Errorf("insert auto id=%d: %w", id, err)
		}
	}
	return nil
}

func resetSequences(ctx context.Context, pool *pgxpool.Pool) error {
	tables := []string{"autohaus", "adresse", "auto"}
	for _, t := range tables {
		query := fmt.Sprintf(
			`SELECT setval(pg_get_serial_sequence('%s', 'id'), COALESCE(MAX(id), 999) + 1, false) FROM %s`,
			t, t,
		)
		if _, err := pool.Exec(ctx, query); err != nil {
			return fmt.Errorf("reset sequence for %s: %w", t, err)
		}
	}
	return nil
}
