package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/myorg/myservice/internal/model"
	"github.com/myorg/myservice/internal/repository"
	"github.com/myorg/myservice/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRepo is a test double for AutohausRepository.
type mockRepo struct {
	findByIDFn func(ctx context.Context, id uint) (*model.Autohaus, error)
	createFn   func(ctx context.Context, a *model.Autohaus) (uint, error)
}

func (m *mockRepo) FindByID(ctx context.Context, id uint) (*model.Autohaus, error) {
	return m.findByIDFn(ctx, id)
}

func (m *mockRepo) Create(ctx context.Context, a *model.Autohaus) (uint, error) {
	return m.createFn(ctx, a)
}

func TestGetByID_ReturnsAutohaus(t *testing.T) {
	expected := &model.Autohaus{
		ID:       100,
		Name:     "Autohaus Karlsruhe",
		Username: "autohaus_karlsruhe",
		Email:    "kontakt@autohaus-karlsruhe.de",
		Gruendungsdatum: time.Date(1998, 5, 12, 0, 0, 0, 0, time.UTC),
		Adresse: model.Adresse{PLZ: "76131", Ort: "Karlsruhe", Land: "Deutschland"},
	}

	repo := &mockRepo{
		findByIDFn: func(_ context.Context, id uint) (*model.Autohaus, error) {
			assert.Equal(t, uint(100), id)
			return expected, nil
		},
	}

	svc := service.New(repo)
	got, err := svc.GetByID(context.Background(), 100)

	require.NoError(t, err)
	assert.Equal(t, expected.ID, got.ID)
	assert.Equal(t, expected.Name, got.Name)
	assert.Equal(t, expected.Adresse.PLZ, got.Adresse.PLZ)
}

func TestGetByID_NotFound(t *testing.T) {
	repo := &mockRepo{
		findByIDFn: func(_ context.Context, _ uint) (*model.Autohaus, error) {
			return nil, repository.ErrNotFound
		},
	}

	svc := service.New(repo)
	got, err := svc.GetByID(context.Background(), 999999)

	assert.Nil(t, got)
	require.Error(t, err)
	assert.True(t, errors.Is(err, repository.ErrNotFound))
}

func TestCreate_ReturnsNewID(t *testing.T) {
	repo := &mockRepo{
		createFn: func(_ context.Context, a *model.Autohaus) (uint, error) {
			assert.Equal(t, "Autohaus Test", a.Name)
			return 1001, nil
		},
	}

	svc := service.New(repo)
	input := &model.Autohaus{
		Name:     "Autohaus Test",
		Username: "autohaus_test",
		Email:    "test@example.de",
	}

	id, err := svc.Create(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, uint(1001), id)
}

func TestCreate_PropagatesError(t *testing.T) {
	repoErr := errors.New("db constraint violation")
	repo := &mockRepo{
		createFn: func(_ context.Context, _ *model.Autohaus) (uint, error) {
			return 0, repoErr
		},
	}

	svc := service.New(repo)
	id, err := svc.Create(context.Background(), &model.Autohaus{})

	assert.Equal(t, uint(0), id)
	require.Error(t, err)
}
