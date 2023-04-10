package bandwidth

import (
	"time"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type consumptionRepository struct {
	app *infinity.Application
	log logrus.FieldLogger
}

func NewConsumptionRepository(app *infinity.Application) infinity.BandwidthConsumptionRepository {
	return &consumptionRepository{
		app: app,
		log: app.Logger.WithField("component", "consumption_repository"),
	}
}

func (repository *consumptionRepository) GetTodaysEntries() ([]*infinity.BandwidthConsumption, error) {
	entries := []*infinity.BandwidthConsumption{}

	err := repository.app.DB.Model(&entries).Where("created_at::date = current_date").Select()
	if err == nil || err == pg.ErrNoRows {
		return entries, nil
	}

	return nil, errors.Wrap(err, "Failed to retrieve today's entries")
}

func (repository *consumptionRepository) AddConsumption(consumption *infinity.BandwidthConsumption) error {
	if consumption.CreatedAt.IsZero() {
		consumption.CreatedAt = time.Now()
	}

	_, err := repository.app.DB.
		Model(consumption).
		OnConflict("(session_uuid, recording_uuid, CAST(created_at as DATE)) DO UPDATE SET bytes = ?", consumption.Bytes).
		Insert()

	return err
}

func (repository *consumptionRepository) GetConsumptionToday(session uuid.UUID) (uint64, error) {

	var consumed uint64

	err := repository.app.DB.
		Model(&infinity.BandwidthConsumption{}).
		ColumnExpr("SUM(bytes)").
		Where("session_uuid = ?", session.String()).
		Select(&consumed)

	return consumed, err
}
