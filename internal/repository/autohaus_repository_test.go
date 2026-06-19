package repository_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myorg/myservice/internal/model"
	"github.com/myorg/myservice/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	ctr, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("autohaus"),
		tcpostgres.WithUsername("autohaus"),
		tcpostgres.WithPassword("p"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = ctr.Terminate(ctx) })

	connStr, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	require.NoError(t, pool.Ping(ctx))

	applyMigrations(t, ctx, pool)
	return pool
}

func applyMigrations(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "..", "..")

	// Use the clean schema file (no tablespace requirement) for testcontainers.
	schemaFile := filepath.Join(root, "db", "schema", "schema.sql")
	content, err := os.ReadFile(schemaFile)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, string(content))
	require.NoError(t, err, "schema migration failed")
}

func TestFindByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	repo := repository.New(pool)

	_, err := repo.FindByID(context.Background(), 999999)

	require.Error(t, err)
	assert.True(t, errors.Is(err, repository.ErrNotFound))
}

func TestCreate_AndFindByID(t *testing.T) {
	pool := setupTestDB(t)
	repo := repository.New(pool)
	ctx := context.Background()

	input := &model.Autohaus{
		Name:            "Autohaus Test",
		Username:        "autohaus_test",
		Email:           "test@autohaus-test.de",
		AnzahlFahrzeuge: 3,
		Gruendungsdatum: time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
		Homepage:        "https://autohaus-test.de",
		Adresse: model.Adresse{
			PLZ:  "76131",
			Ort:  "Karlsruhe",
			Land: "Deutschland",
		},
		Autos: []model.Auto{
			{Kennzeichen: "KA-TS-001", Marke: "BMW", Modell: "X3", Baujahr: 2022},
		},
	}

	id, err := repo.Create(ctx, input)
	require.NoError(t, err)
	assert.Greater(t, id, uint(0))

	got, err := repo.FindByID(ctx, id)
	require.NoError(t, err)

	assert.Equal(t, id, got.ID)
	assert.Equal(t, "Autohaus Test", got.Name)
	assert.Equal(t, "autohaus_test", got.Username)
	assert.Equal(t, "76131", got.Adresse.PLZ)
	assert.Equal(t, "Karlsruhe", got.Adresse.Ort)
	require.Len(t, got.Autos, 1)
	assert.Equal(t, "KA-TS-001", got.Autos[0].Kennzeichen)
}
