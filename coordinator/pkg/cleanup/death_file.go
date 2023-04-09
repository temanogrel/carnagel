package cleanup

import (
	"encoding/csv"
	"io"
	"os"

	"context"
	"time"

	"strconv"
	"sync"

	"strings"

	"sync/atomic"

	"git.misc.vee.bz/carnagel/coordinator/pkg"
	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type upstoreFile struct {
	url       string
	filename  string
	uploaded  time.Time
	downloads uint64
}

type deathFile struct {
	app *coordinator.Application
	log logrus.FieldLogger

	report *coordinator.DeathFileReport
}

func NewDeathFileService(app *coordinator.Application) coordinator.DeathFileService {
	return &deathFile{
		app: app,
		log: app.Logger.WithField("component", "death-file"),
	}
}

func (df *deathFile) GetReport() *coordinator.DeathFileReport {
	return df.report
}

func (df *deathFile) Process(ctx context.Context, routines int, source string) error {

	df.report = &coordinator.DeathFileReport{
		MissingInHermes:    new(uint64),
		MissingInAphrodite: new(uint64),
		Skipped:            new(uint64),
		Found:              new(uint64),
		Deleted:            new(uint64),
	}

	sourceFile, err := os.Open(source)
	if err != nil {
		return errors.Wrap(err, "Failed to open source")
	}

	defer sourceFile.Close()

	wg := &sync.WaitGroup{}

	processingQueue := make(chan *upstoreFile)

	// Start a bunch of goroutines for processing
	for i := 0; i <= routines; i++ {
		wg.Add(1)

		go df.processUpstoreFile(ctx, wg, processingQueue)
	}

	// Parse the source file
	r := csv.NewReader(sourceFile)
	r.Comma = ';'
	r.TrimLeadingSpace = true
	r.FieldsPerRecord = 4

RuntimeLoop:
	for {
		select {
		case <-ctx.Done():
			return nil

		default:
			record, err := r.Read()
			if err == io.EOF {
				df.log.Info("Received EOF")

				close(processingQueue)
				break RuntimeLoop

			} else if err == csv.ErrFieldCount {
				df.log.WithField("record", record).Warn("Record found with unexpected number of fields")
				continue
			}

			if err != nil {
				df.log.WithError(err).Error("Received an error reading from file")

				return err
			}

			log := df.log.WithField("record", record)

			uploadedAt, err := time.Parse("2006-01-02 15:04:05", record[2])
			if err != nil {
				log.WithError(err).Error("Failed to parse uploaded at")
				continue
			}

			downloads, err := strconv.ParseUint(record[3], 10, 32)
			if err != nil {
				log.WithError(err).Error("Failed to parse downloads")
				continue
			}

			if downloads > 0 {
				atomic.AddUint64(df.report.Skipped, 1)
				continue
			}

			processingQueue <- &upstoreFile{
				url:       strings.Replace(record[0], "http://upsto.re/", "https://upstore.net/", 1),
				filename:  record[1],
				uploaded:  uploadedAt,
				downloads: downloads,
			}
		}
	}

	wg.Wait()

	df.log.WithField("report", df.report).Info("Finished processing the death file")

	return nil
}

// processUpstoreFile will process the incoming channel of upstore files and if the deletion is successful writes to
// the out channel
func (df *deathFile) processUpstoreFile(ctx context.Context, wg *sync.WaitGroup, queue chan *upstoreFile) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case file, ok := <-queue:
			if !ok {
				return
			}

			log := df.log.WithField("upstoreFile", file)

			hermesUrl, err := df.app.HermesClient.GetByOriginalUrl(file.url)
			if err == ecosystem.HermesUrlNotFoundErr {
				atomic.AddUint64(df.report.MissingInHermes, 1)
				log.Warn("Failed to locate file queue hermes")
				continue
			}

			if err != nil {
				log.WithError(err).Error("Unexpected error occurred retrieving from hermes")
				continue
			}

			recording, err := df.app.AphroditeClient.GetRecordingByVideoUrl(hermesUrl.ShortURL)
			if err == ecosystem.RecordingNotFoundErr {
				atomic.AddUint64(df.report.MissingInAphrodite, 1)

				log.
					WithField("url", hermesUrl.ShortURL).
					Warn("Failed to locate recording queue aphrodite")

				continue
			}

			if err != nil {
				log.WithError(err).Error("Unexpected error occurred retrieving from aphrodite")
				continue
			}

			atomic.AddUint64(df.report.Found, 1)

			if err := df.app.CleanupService.Delete(recording, ""); err != nil {
				// Check if infinity gave us an error indicating that the recording shouldn't be deleted due to:
				// 1. views
				// 2. is part of a premium user's collection
				if err == coordinator.RecordingHasInfinityViewsErr || err == coordinator.RecordingInInfinityPremiumUserCollectionErr {
					log.
						WithField("recording", recording).
						Info("Not deleting recording because infinity reported views or part of premium user collection")

					atomic.AddUint64(df.report.Skipped, 1)
				} else {
					log.WithError(err).Error("Failed to delete recording")
				}

				continue
			}

			atomic.AddUint64(df.report.Deleted, 1)
			log.WithField("recordingId", recording.Id).Debug("Deleted recording")
		}
	}
}
