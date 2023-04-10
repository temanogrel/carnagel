package bandwidth

import (
	"context"
	"github.com/sasha-s/go-deadlock"
	"time"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// todo: figure when what a good interval to persist to the database is
const synchronizationInterval = time.Hour

type recordingByteCount map[uuid.UUID]uint64
type userSessionsMap map[uuid.UUID]uuid.UUID
type sessionConsumptionMap map[uuid.UUID]recordingByteCount
type blackListedSessionsMap map[uuid.UUID]bool

type consumptionCollector struct {
	app *infinity.Application
	log logrus.FieldLogger

	mtx                 deadlock.RWMutex
	sessionConsumption  sessionConsumptionMap
	userSessions        userSessionsMap
	blackListedSessions blackListedSessionsMap
}

func NewConsumptionCollector(app *infinity.Application) infinity.BandwidthConsumptionCollector {
	collector := &consumptionCollector{
		app: app,
		log: app.Logger.WithField("component", "consumption_collector"),
		mtx: deadlock.RWMutex{},

		userSessions:        make(userSessionsMap),
		sessionConsumption:  make(sessionConsumptionMap),
		blackListedSessions: make(blackListedSessionsMap),
	}

	go collector.RebuildFromToday()

	return collector
}

func (service *consumptionCollector) RebuildFromToday() error {
	log := service.app.Logger.WithField("operation", "rebuildBandwidthConsumption")
	log.Debug("Rebuilding bandwidth consumption from today")

	entries, err := service.app.BandwidthConsumptionRepository.GetTodaysEntries()
	if err != nil {
		log.WithError(err).Error("Failed to rebuild bandwidth collector")

		return errors.Wrap(err, "Failed to rebuild collector")
	}

	for _, entry := range entries {
		service.AddConsumption(entry)
	}

	log.Debug("Finished rebuilding bandwidth consumption from today")

	return nil
}

func (service *consumptionCollector) Synchronizer(ctx context.Context) {

	// Midnight in UTC we reset the bandwidth
	midnight := time.Now().
		Truncate(time.Hour * 24).
		Add(time.Hour * 24).
		UTC()

	syncTimer := time.NewTimer(synchronizationInterval)
	resetTimer := time.NewTimer(time.Until(midnight))

	for {
		select {
		case <-ctx.Done():
			syncTimer.Stop()

			// save the current state to database
			service.sendToDatabase()

			return

		case <-syncTimer.C:
			service.sendToDatabase()

			syncTimer.Reset(synchronizationInterval)

		case <-resetTimer.C:
			service.log.Info("Resetting user & session bandwidth consumption")

			service.sendToDatabase()
			service.reset()

			resetTimer.Reset(time.Hour * 24)
		}
	}
}

// reset will clear the user sessions map and session consumption
func (service *consumptionCollector) reset() {
	service.mtx.Lock()
	service.userSessions = make(userSessionsMap)
	service.sessionConsumption = make(sessionConsumptionMap)
	service.blackListedSessions = make(blackListedSessionsMap)
	service.mtx.Unlock()
}

func (service *consumptionCollector) sendToDatabase() {
	startedAt := time.Now()

	service.mtx.RLock()
	defer service.mtx.RUnlock()

	for session, recordings := range service.sessionConsumption {
		for recordingUuid, bytes := range recordings {
			consumption := &infinity.BandwidthConsumption{
				SessionUuid:   session,
				RecordingUuid: recordingUuid,
				Bytes:         bytes,
			}

			if userUuid, ok := service.userSessions[session]; ok {
				consumption.UserUuid = uuid.NullUUID{
					Valid: true,
					UUID:  userUuid,
				}
			}

			if err := service.app.BandwidthConsumptionRepository.AddConsumption(consumption); err != nil {
				service.log.WithError(err).Errorf("Failed to persist consumption to database")
			}
		}
	}

	service.log.
		WithField("entries", len(service.sessionConsumption)).
		WithField("duration", time.Since(startedAt).Seconds()).
		Debug("Persisted bandwidth usage to database")
}

func (service *consumptionCollector) AddConsumption(consumption *infinity.BandwidthConsumption) (uint64, uint64) {

	service.mtx.Lock()

	recordingConsumption, ok := service.sessionConsumption[consumption.SessionUuid]
	if !ok {
		recordingConsumption = make(recordingByteCount)
		service.sessionConsumption[consumption.SessionUuid] = recordingConsumption
	}

	onRecording, _ := recordingConsumption[consumption.RecordingUuid]
	onRecording += consumption.Bytes

	recordingConsumption[consumption.RecordingUuid] = onRecording

	// associate this session to a user
	if consumption.UserUuid.Valid {
		service.userSessions[consumption.SessionUuid] = consumption.UserUuid.UUID
	}

	// release write lock before calling get total consumption today
	service.mtx.Unlock()

	total, _ := service.GetTotalConsumptionToday(consumption.SessionUuid)

	return total, onRecording
}

func (service *consumptionCollector) GetTotalConsumptionToday(session uuid.UUID) (uint64, error) {
	service.mtx.RLock()
	defer service.mtx.RUnlock()

	if _, ok := service.blackListedSessions[session]; ok {
		return 0, infinity.BlackListedSessionErr
	}

	recordingMap, ok := service.sessionConsumption[session]
	if !ok {
		return 0, nil
	}

	var total uint64

	for _, bytes := range recordingMap {
		total += bytes
	}

	return total, nil
}

func (service *consumptionCollector) AddBlackListedSession(session uuid.UUID) {
	service.mtx.Lock()
	defer service.mtx.Unlock()

	service.blackListedSessions[session] = true
}

func (service *consumptionCollector) GetTotalConsumptionTodayOnRecording(session uuid.UUID, recordingId uuid.UUID) (uint64, error) {
	service.mtx.RLock()
	defer service.mtx.RUnlock()

	if _, ok := service.blackListedSessions[session]; ok {
		return 0, infinity.BlackListedSessionErr
	}

	recordingMap, ok := service.sessionConsumption[session]
	if !ok {
		return 0, nil
	}

	onRecording, ok := recordingMap[recordingId]
	if !ok {
		return 0, nil
	}

	return onRecording, nil
}
