package performer

import (
	"fmt"

	"database/sql"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
)

type service struct {
	app *infinity.Application
}

func NewService(app *infinity.Application) infinity.PerformerService {
	service := &service{
		app: app,
	}

	go service.updatePerformersWithoutSlug()

	return service
}

func (service *service) ImportFromAphrodite(raw *ecosystem.Performer) (*infinity.Performer, error) {

	performer, err := service.app.PerformerRepository.GetByExternalId(raw.Id)
	switch err {
	case nil:
		if err := service.updatePerformerWithSlug(performer); err != nil {
			service.app.Logger.WithError(err).Error("Failed to update performer with slug")

			return nil, err
		}

		return performer, service.updateFromAphrodite(raw, performer)

	case infinity.PerformerNotFoundErr:
		return service.createPerformer(raw)

	default:
		return nil, errors.Wrap(err, "Failed to import from aphrodite due to database error")
	}
}

func (service *service) createPerformer(raw *ecosystem.Performer) (*infinity.Performer, error) {
	performer := &infinity.Performer{
		ExternalId: raw.Id,

		OriginService: raw.Service,
		OriginSection: sql.NullString{
			Valid:  raw.ServiceSection != "",
			String: raw.ServiceSection,
		},

		StageName: raw.StageName,
		Aliases:   raw.Aliases,
	}

	if err := service.updatePerformerWithSlug(performer); err != nil {
		service.app.Logger.WithError(err).Error("Failed to update performer with slug")

		return nil, err
	}

	if err := service.app.PerformerRepository.Create(performer); err != nil {
		return performer, errors.Wrap(err, "Failed to create performer")
	}

	if err := service.app.PerformerSearchService.AddPerformerAlias(performer.Uuid, performer.StageName); err != nil {
		return performer, errors.Wrap(err, "Failed to add performer alias")
	}

	return performer, nil
}

func (service *service) updateFromAphrodite(raw *ecosystem.Performer, performer *infinity.Performer) error {

	performer.StageName = raw.StageName
	performer.Aliases = raw.Aliases

	if err := service.app.PerformerSearchService.AddPerformerAlias(performer.Uuid, performer.StageName); err != nil {
		return errors.Wrap(err, "Failed to add performer alias")
	}

	return service.app.PerformerRepository.Update(performer)
}

func (service *service) updatePerformersWithoutSlug() {
	log := service.app.Logger.WithField("operation", "updatePerformersWithoutSlug")
	log.Info("Updating performers missing slug")

	performers, err := service.app.PerformerRepository.GetAllMissingSlug()
	if err != nil {
		log.WithError(err).Error("Failed to retrieve performers without slug")
		return
	}

	for _, performer := range performers {
		if err := service.updatePerformerWithSlug(&performer); err != nil {
			log.WithError(err).Error("Failed to generate slug for performer")
		}

		if err := service.app.PerformerRepository.Update(&performer); err != nil {
			log.WithError(err).Error("Failed to update performer with slug")
		}
	}

	log.Info("Finished updating performers missing slug")
}

// Code is basically ported straight from ultron
func (service *service) updatePerformerWithSlug(performer *infinity.Performer) error {
	serviceName, err := performer.GetFullOriginServiceName()
	if err != nil {
		return err
	}

	formatString := fmt.Sprintf("%s %s", serviceName, performer.StageName)
	generatedSlug := slug.Make(formatString)

	if performer.Slug == generatedSlug {
		return nil
	}

	for i := 1; true; i++ {
		if _, err := service.app.PerformerRepository.GetBySlug(generatedSlug); err != nil {
			if err == infinity.PerformerNotFoundErr {
				performer.Slug = generatedSlug
				return nil
			}

			return err
		}

		generatedSlug = slug.Make(fmt.Sprintf("%s %d", formatString, i))
	}

	return nil
}
