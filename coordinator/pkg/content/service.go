package content

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"git.misc.vee.bz/carnagel/coordinator/pkg"
	"git.misc.vee.bz/carnagel/coordinator/pkg/aphrodite/repository"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain/wordpress"
	"git.misc.vee.bz/carnagel/go-ecosystem/protobuf/common"
	"git.misc.vee.bz/carnagel/minerva/pkg"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func NewContentService(application *coordinator.Application) coordinator.ContentService {
	return &contentService{
		app:    application,
		logger: application.Logger.WithField("component", "content"),
	}
}

type contentService struct {
	app    *coordinator.Application
	logger logrus.FieldLogger
}

func (service *contentService) Publish(recording *ecosystem.Recording) error {

	performer, err := service.app.AphroditeClient.GetPerformer(recording.PerformerId)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve performer")
	}

	if recording.ImageUrls.Large == "" {
		if err := service.uploadWordpressCollage(recording); err != nil {
			if err == ecosystem.FailedToCallApiErr {
				return ecosystem.FailedToCallApiErr
			}

			return errors.Wrap(err, "Failed to upload collage to camgirl.gallery")
		}
	}

	if recording.VideoUrl == "" {
		if recording.UpstoreHash == "" {
			return coordinator.UpstoreHashMissingErr
		}

		if err := service.setVideoUrl(recording.UpstoreHash, recording); err != nil {
			return errors.Wrap(err, "Failed to set hermes url")
		}
	}

	if err := service.publishToUltron(recording, performer); err != nil {
		return errors.Wrap(err, "Failed to publish to ultron")
	}

	if err := service.publishToWordpress(recording); err != nil {
		return errors.Wrap(err, "Failed to publish to wordpress")
	}

	if err := service.publishToInfinity(recording, performer); err != nil {
		return errors.Wrap(err, "Failed to publish to infinity")
	}

	if err := service.app.MinervaClient.RequestDeletion(recording.VideoMp4Uuid); err != nil {
		return errors.Wrap(err, "Failed to delete mp4 file")
	}

	recording.State = "published"
	recording.VideoMp4Uuid = uuid.Nil
	recording.LastCheckedAt = time.Now().Truncate(time.Second)

	if err := service.app.AphroditeClient.UpdateRecording(recording); err != nil {
		return errors.Wrap(err, "Failed to update recording in aphrodite")
	}

	return nil
}

func (service *contentService) removeFromInfinity(recording *ecosystem.Recording) error {
	resp, err := service.app.InfinityContentClient.DeleteRecording(context.Background(), &common.RecordingIdentifier{
		Id: recording.Id,
	})

	if err != nil {
		return err
	}

	switch resp.Status {
	case common.StatusCode_Ok:
		return nil

	case common.StatusCode_RecordingHasViews:
		return coordinator.RecordingHasInfinityViewsErr

	case common.StatusCode_RecordingInPremiumUserCollection:
		return coordinator.RecordingInInfinityPremiumUserCollectionErr

	default:
		return errors.New("Unknown error deleting recording from infinity")
	}
}

func (service *contentService) removeFromWordpress(recording *ecosystem.Recording) error {
	sites, err := service.app.AphroditeClient.GetSites()
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve sites from aphrodite")
	}

	wpClient := wordpress.NewWordpressClient()

	for _, site := range sites {
		if postId, ok := recording.IsPublishedOnSite(site); ok {
			if err := wpClient.Delete(site, postId); err != nil {
				return errors.Wrapf(err, "Failed to delete from wordpress site: %s", site.Name)
			}
		}
	}

	if err := service.app.UltronClient.Delete(recording); err != nil {
		return errors.Wrap(err, "Failed to delete from ultron")
	}

	return nil
}

func (service *contentService) Remove(recording *ecosystem.Recording, system string) error {
	switch system {
	case "":
		// If no system is specified, remove from both
		if err := service.removeFromInfinity(recording); err != nil {
			return err
		}

		return service.removeFromWordpress(recording)

	case ecosystem.SystemWordpress:
		return service.removeFromWordpress(recording)

	case ecosystem.SystemInfinity:
		return service.removeFromInfinity(recording)

	default:
		return nil
	}
}

func (service *contentService) publishToUltron(recording *ecosystem.Recording, performer *ecosystem.Performer) error {

	err := service.app.UltronClient.Create(&ecosystem.UltronPost{
		Recording: recording,
		Performer: performer,
	})

	if err != nil && err != ecosystem.UltronEntryAlreadyExistsErr {
		return errors.Wrap(err, "Failed to publish to the new ultron site")
	}

	return nil
}

func (service *contentService) publishToInfinity(recording *ecosystem.Recording, performer *ecosystem.Performer) error {
	resp, err := service.app.InfinityContentClient.UpsertRecording(context.TODO(), &common.RecordingIdentifier{
		Id: recording.Id,
	})

	if err != nil {
		return errors.Wrap(err, "Failed to call infinity content server")
	}

	if resp.Status != common.StatusCode_Ok {
		return errors.Wrapf(err, "Unexpected response code %d received from infinity content server", resp.Status)
	}

	return nil
}

func (service *contentService) publishToWordpress(recording *ecosystem.Recording) error {
	sites, err := service.app.AphroditeClient.GetSites()
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve sites from aphrodite")
	}

	for _, site := range sites {
		if !site.Enabled {
			continue
		}

		if !recording.IsViableForSite(site) {
			continue
		}

		post := &ecosystem.WordpressPost{
			Title:      recording.GetPostTitle(),
			ContentRaw: recording.GetWordpressContent(),

			Status:        "publish",
			CommentStatus: "closed",
		}

		client := wordpress.NewWordpressClient()

		if postId, ok := recording.IsPublishedOnSite(site); ok {
			post.Id = postId

			if err := client.Update(site, post); err != nil {
				return errors.Wrapf(err, "Failed to update existing post on %s with id %d", site.Name, post.Id)
			}

		} else {

			if err := client.Create(site, post); err != nil {
				return errors.Wrapf(err, "Failed to publish to %s with base api uri %s", site.Name, site.ApiUri)
			}
		}

		if err := service.app.AphroditeClient.CreatePostAssociation(recording, uint64(site.Id), post.Id); err != nil {
			return errors.Wrapf(err, "Failed to create post association with %s and id %d", site.Name, post.Id)
		}
	}

	return nil
}

func (service *contentService) uploadWordpressCollage(recording *ecosystem.Recording) error {

	service.logger.WithField("recording", recording).Debug("Uploading wordpress collage")

	path, err := service.app.MinervaClient.Download(recording.WordpressCollageUuid)
	if err != nil {
		if err == minerva.FileNotFoundErr {
			if err := service.notifyMissingCollage(recording); err != nil {
				service.logger.WithError(err).Error("Failed notify about missing collage")
			}
		}

		return err
	}

	defer os.Remove(path)

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	urls, err := service.app.CamgirlGalleryClient.Upload(recording.GetPublishedFilename("jpg"), file)
	if err != nil {
		return err
	}

	recording.ImageUrls = urls

	// Update ecosystem after every state change
	if err := service.app.AphroditeClient.UpdateRecording(recording); err != nil {
		return err
	}

	return nil
}

func (service *contentService) notifyMissingCollage(recording *ecosystem.Recording) error {
	ch, err := ecosystem.GetAmqpChannel(service.app.Amqp, 0)
	if err != nil {
		return err
	}

	defer ch.Close()

	body, err := json.Marshal(ecosystem.SystemScopedRecordingPayload{
		RecordingId: recording.Id,
		System:      ecosystem.SystemWordpress,
	})

	if err != nil {
		return errors.Wrap(err, "Failed to marshal to json")
	}

	ch.Publish("", ecosystem.AmqpQueueImageRegeneration, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})

	return nil
}

func (service *contentService) setVideoUrl(upstoreHash string, recording *ecosystem.Recording) error {
	if len(upstoreHash) < 4 {
		return errors.Errorf("Upstore hash %s is too short to be valid", upstoreHash)
	}

	hermesUrl, err := service.app.HermesClient.Create(fmt.Sprintf("https://upstore.net/%s", upstoreHash))
	if err != nil {
		return errors.Wrap(err, "Failed to generate hermes url")
	}

	recording.VideoUrl = hermesUrl.ShortURL

	// Update ecosystem after every state change
	if err := service.app.AphroditeClient.UpdateRecording(recording); err != nil {
		return err
	}

	return nil
}

func (service *contentService) RebuildUltron(lastSeenId uint64) error {

	ctx := context.Background()
	repo := repository.NewRecordingRepository(service.app)

	recordings, err := repo.GetPublishedRecordingIds(ctx, lastSeenId)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve recordings from ecosystem")
	}

	wg := sync.WaitGroup{}

	publish := func() {
		wg.Add(1)
		defer wg.Done()

		for recordingId := range recordings {

			logger := service.app.Logger.WithFields(logrus.Fields{
				"recordingId": recordingId,
			})

			recording, err := service.app.AphroditeClient.GetRecording(recordingId)
			if err != nil {
				logger.WithError(err).Warn("Failed to retrieve recording")
				continue
			}

			performer, err := service.app.AphroditeClient.GetPerformer(recording.PerformerId)
			if err != nil {
				logger.WithError(err).Warn("Failed to retrieve performer")
				continue
			}

			if err := service.publishToUltron(recording, performer); err != nil {
				logger.WithError(err).Error("Failed to publish recording")
				continue
			}

			logger.Debug("Successfully added recording")
		}
	}

	for i := 0; i <= 150; i++ {
		go publish()

		// We do this to stop huge requests at the same time
		time.Sleep(10 * time.Millisecond)
	}

	wg.Wait()

	return nil
}
