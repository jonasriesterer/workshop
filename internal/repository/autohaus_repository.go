package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myorg/myservice/internal/db"
	"github.com/myorg/myservice/internal/model"
)

// ErrNotFound is returned when the requested entity does not exist.
var ErrNotFound = errors.New("not found")

// AutohausRepository defines the persistence operations for Autohaus.
type AutohausRepository interface {
	FindByID(ctx context.Context, id uint) (*model.Autohaus, error)
	Create(ctx context.Context, a *model.Autohaus) (uint, error)
}

type autohausRepository struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

// New creates a new AutohausRepository backed by the given connection pool.
func New(pool *pgxpool.Pool) AutohausRepository {
	return &autohausRepository{
		queries: db.New(pool),
		pool:    pool,
	}
}

func (r *autohausRepository) FindByID(ctx context.Context, id uint) (*model.Autohaus, error) {
	row, err := r.queries.GetAutohausById(ctx, int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	adresseRow, err := r.queries.GetAdresseByAutohausId(ctx, row.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	autoRows, err := r.queries.GetAutosByAutohausId(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	return toModel(row, adresseRow, autoRows), nil
}

func (r *autohausRepository) Create(ctx context.Context, a *model.Autohaus) (uint, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	q := r.queries.WithTx(tx)

	var gruendungsdatum pgtype.Date
	_ = gruendungsdatum.Scan(a.Gruendungsdatum)

	var homepage pgtype.Text
	if a.Homepage != "" {
		homepage = pgtype.Text{String: a.Homepage, Valid: true}
	}

	var telefonnr pgtype.Text
	if a.Telefonnr != "" {
		telefonnr = pgtype.Text{String: a.Telefonnr, Valid: true}
	}

	newID, err := q.CreateAutohaus(ctx, db.CreateAutohausParams{
		Version:         int32(a.Version),
		Name:            a.Name,
		Username:        a.Username,
		Email:           a.Email,
		AnzahlFahrzeuge: int32(a.AnzahlFahrzeuge),
		Gruendungsdatum: gruendungsdatum,
		Homepage:        homepage,
		Telefonnummer:   telefonnr,
	})
	if err != nil {
		return 0, err
	}

	if err := q.CreateAdresse(ctx, db.CreateAdresseParams{
		Plz:        a.Adresse.PLZ,
		Ort:        a.Adresse.Ort,
		Land:       a.Adresse.Land,
		AutohausID: newID,
	}); err != nil {
		return 0, err
	}

	for _, auto := range a.Autos {
		if err := q.CreateAuto(ctx, db.CreateAutoParams{
			Kennzeichen: auto.Kennzeichen,
			Marke:       auto.Marke,
			Modell:      auto.Modell,
			Baujahr:     int32(auto.Baujahr),
			AutohausID:  newID,
		}); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return uint(newID), nil
}

func toModel(row db.Autohau, adresse db.Adresse, autos []db.Auto) *model.Autohaus {
	ah := &model.Autohaus{
		ID:              uint(row.ID),
		Version:         uint(row.Version),
		Name:            row.Name,
		Username:        row.Username,
		Email:           row.Email,
		AnzahlFahrzeuge: int(row.AnzahlFahrzeuge),
		Adresse: model.Adresse{
			ID:         uint(adresse.ID),
			PLZ:        adresse.Plz,
			Ort:        adresse.Ort,
			Land:       adresse.Land,
			AutohausID: uint(adresse.AutohausID),
		},
	}

	if row.Gruendungsdatum.Valid {
		ah.Gruendungsdatum = row.Gruendungsdatum.Time
	}
	if row.Homepage.Valid {
		ah.Homepage = row.Homepage.String
	}
	if row.Telefonnummer.Valid {
		ah.Telefonnr = row.Telefonnummer.String
	}
	if row.Erzeugt.Valid {
		ah.Erzeugt = row.Erzeugt.Time
	} else {
		ah.Erzeugt = time.Now()
	}
	if row.Aktualisiert.Valid {
		ah.Aktualisiert = row.Aktualisiert.Time
	} else {
		ah.Aktualisiert = time.Now()
	}

	for _, a := range autos {
		ah.Autos = append(ah.Autos, model.Auto{
			ID:          uint(a.ID),
			Kennzeichen: a.Kennzeichen,
			Marke:       a.Marke,
			Modell:      a.Modell,
			Baujahr:     int(a.Baujahr),
			AutohausID:  uint(a.AutohausID),
		})
	}

	return ah
}
