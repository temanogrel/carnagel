package recording

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/sasha-s/go-deadlock"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
)

const RecordingsPerPage = 90
const RedisKeyPageCount = "infinity:page-count"

var getRedisKeyCacheEntryOfPage = func(page uint64) string {
	return fmt.Sprintf("infinity:page-%d", page)
}

type cacheEntry struct {
	Mutex              deadlock.RWMutex `json:"-"`
	StartTime          time.Time        `json:"startTime"`
	EndTime            time.Time        `json:"endTime"`
	RecordingCount     uint             `json:"recordingCount"`
	OriginSectionCount map[string]uint  `json:"originSectionCount"`
	OriginServiceCount map[string]uint  `json:"originServiceCount"`
}

func (entry *cacheEntry) isFull() bool {
	return entry.RecordingCount == RecordingsPerPage
}

type repositoryCache struct {
	app *infinity.Application

	repository *repository

	// cache, maps page number to a start and end time
	mtx        deadlock.RWMutex
	pageCount  uint64
	cache      map[uint64]*cacheEntry
	cacheBuilt bool
}

func NewRepository(app *infinity.Application) infinity.RecordingRepository {
	cache := &repositoryCache{
		app:        app,
		repository: NewWrappedRepository(app),

		mtx:        deadlock.RWMutex{},
		cache:      make(map[uint64]*cacheEntry),
		cacheBuilt: false,
	}

	return cache
}

func (cache *repositoryCache) LoadPageCacheFromRedis() {
	log := cache.app.Logger.WithField("operation", "LoadPageCacheFromRedis")
	log.Info("Loading the recording page cache from redis")

	pageCount, err := cache.readRedisPageCount()
	if err != nil {
		cache.RebuildPageCache(1)
		return
	}

	criteria := &infinity.RecordingRepositoryCriteria{
		SortMode: infinity.SortModeOldest,
		Limit:    RecordingsPerPage,
	}

	// If entire cache is not loaded, rebuild from offset
	if _, count, _ := cache.repository.Matching(criteria, true); uint64(count) > pageCount*RecordingsPerPage {
		cache.RebuildPageCache(pageCount)
		return
	}

	cache.pageCount = pageCount

	for i := 1; uint64(i) <= pageCount; i++ {
		cacheEntry, err := cache.readRedisEntryForCacheEntry(uint64(i))
		if err != nil {
			cache.RebuildPageCache(uint64(i))
			return
		}

		cache.mtx.Lock()
		cache.cache[uint64(i)] = cacheEntry
		cache.mtx.Unlock()
	}

	cache.mtx.Lock()
	defer cache.mtx.Unlock()
	log.Info("Finished loading page cache from redis")
	cache.cacheBuilt = true
}

func (cache *repositoryCache) RebuildPageCache(offsetPage uint64) {
	cache.cacheBuilt = false

	log := cache.app.Logger.
		WithFields(logrus.Fields{"operation": "RebuildPageCache", "offsetPage": offsetPage})
	log.Info("Rebuilding recordings page cache")

	criteria := &infinity.RecordingRepositoryCriteria{
		SortMode: infinity.SortModeOldest,
		Limit:    RecordingsPerPage,
	}

	var lastSeenDate time.Time

	for i := offsetPage; true; i++ {
		if lastSeenDate.IsZero() {
			criteria.Offset = int((offsetPage - 1) * RecordingsPerPage)
		} else {
			criteria.CreatedAfter = lastSeenDate
			criteria.Offset = 0
		}

		log.WithField("criteria", criteria).Debug("Retrieving from database with criteria")
		recordings, _, err := cache.repository.Matching(criteria, false)
		if err != nil {
			log.WithField("criteria", criteria).WithError(err).Error("Failed to retrieve records for criteria")
			return
		}

		log.WithField("criteria", criteria).Debug("Received recordings from database")

		if len(recordings) == 0 {
			cache.cacheBuilt = true
			log.Info("Finished rebuilding page cache")
			return
		}

		cacheEntry := cacheEntry{
			StartTime:          recordings[0].CreatedAt,
			EndTime:            recordings[len(recordings)-1].CreatedAt,
			RecordingCount:     uint(len(recordings)),
			Mutex:              deadlock.RWMutex{},
			OriginSectionCount: map[string]uint{},
			OriginServiceCount: map[string]uint{},
		}

		// Add tracking of each origin service/section
		for _, recording := range recordings {
			if _, ok := cacheEntry.OriginServiceCount[recording.Performer.OriginService]; !ok {
				cacheEntry.OriginServiceCount[recording.Performer.OriginService] = 1
			} else {
				cacheEntry.OriginServiceCount[recording.Performer.OriginService] += 1
			}

			if recording.Performer.OriginSection.Valid {
				if _, ok := cacheEntry.OriginSectionCount[recording.Performer.OriginSection.String]; !ok {
					cacheEntry.OriginSectionCount[recording.Performer.OriginSection.String] = 1
				} else {
					cacheEntry.OriginSectionCount[recording.Performer.OriginSection.String] += 1
				}
			}
		}

		if err = cache.updateRedisEntryForCacheEntry(uint64(i), &cacheEntry); err != nil {
			log.WithError(err).Error("Failed to update redis page cache entry")
			return
		}

		if err = cache.updateRedisPageCount(uint64(i)); err != nil {
			log.WithError(err).Error("Failed to update redis page count")
			return
		}

		cache.mtx.Lock()
		cache.pageCount = uint64(i)
		cache.cache[uint64(i)] = &cacheEntry
		cache.mtx.Unlock()

		if len(recordings) < RecordingsPerPage {
			cache.cacheBuilt = true
			log.Info("Finished rebuilding page cache")
			return
		}

		lastSeenDate = cacheEntry.EndTime
	}
}

func (cache *repositoryCache) readRedisEntryForCacheEntry(page uint64) (*cacheEntry, error) {
	log := cache.app.Logger.WithField("operation", "readRedisEntryForCacheEntry")

	key := getRedisKeyCacheEntryOfPage(page)
	redisEntry, err := cache.app.Redis.Get(key).Result()
	if err != nil || redisEntry == "" {
		log.
			WithField("redisKey", key).
			Error("Could not load cache entry from redis or it does not exist in redis")

		return nil, err
	}

	cacheEntry := &cacheEntry{}
	if err := json.Unmarshal([]byte(redisEntry), cacheEntry); err != nil {
		log.
			WithFields(logrus.Fields{
				"redisKey":   key,
				"redisEntry": redisEntry,
			}).WithError(err).Error("Failed to json unmarshal redis cache entry")

		return nil, err
	}

	// Ensure it's consistent with latest cache entry structure
	if cacheEntry.StartTime.IsZero() ||
		cacheEntry.EndTime.IsZero() ||
		cacheEntry.RecordingCount == 0 ||
		cacheEntry.OriginServiceCount == nil ||
		cacheEntry.OriginSectionCount == nil {
		log.
			WithFields(logrus.Fields{
				"redisKey":   key,
				"redisEntry": redisEntry,
			}).Warn("Cache entry structure from redis is not consistent with cache struct, rebuilding page cache")

		return nil, errors.New("Redis entry is not consistent with cache struct")
	}

	cacheEntry.Mutex = deadlock.RWMutex{}

	return cacheEntry, nil
}

func (cache *repositoryCache) updateRedisEntryForCacheEntry(page uint64, cacheEntry *cacheEntry) error {
	key := getRedisKeyCacheEntryOfPage(page)

	log := cache.app.Logger.WithField("operation", "updateRedisEntryForCacheEntry")

	redisEntry, err := json.Marshal(cacheEntry)
	if err != nil {
		log.WithError(err).Error("Failed to json marshal the cache entry so it can be placed in redis")

		return err
	}

	if _, err := cache.app.Redis.Set(key, redisEntry, 0).Result(); err != nil {
		log.WithFields(logrus.Fields{
			"redisKey":   key,
			"redisEntry": redisEntry,
		}).WithError(err).Error("Failed to put cache entry in redis")

		return err
	}

	return nil
}

func (cache *repositoryCache) readRedisPageCount() (uint64, error) {
	log := cache.app.Logger.WithField("operation", "readRedisPageCount")

	pageCount, err := cache.app.Redis.Get(RedisKeyPageCount).Uint64()
	if err != nil || pageCount == 0 {
		log.
			WithField("redisKey", RedisKeyPageCount).
			Error("Could not load the page count key from redis or it does not exist in redis")

		return 0, err
	}

	return pageCount, nil
}

func (cache *repositoryCache) updateRedisPageCount(pageCount uint64) error {
	log := cache.app.Logger.WithField("operation", "updateRedisPageCount")

	if _, err := cache.app.Redis.Set(RedisKeyPageCount, pageCount, 0).Result(); err != nil {
		log.WithFields(logrus.Fields{
			"redisKey":   RedisKeyPageCount,
			"redisEntry": pageCount,
		}).WithError(err).Error("Failed to put page count entry into redis")

		return err
	}

	return nil
}

func (cache *repositoryCache) injectPageCacheEntry(criteria *infinity.RecordingRepositoryCriteria) bool {
	cache.mtx.RLock()
	defer cache.mtx.RUnlock()

	// Inject cached stuff if possible
	if cache.cacheBuilt && cache.cacheUsableWithCriteria(criteria) {
		var page uint64

		// Newest to oldest sort mode
		if criteria.SortMode == "" {
			page = cache.pageCount - uint64(math.Floor(float64(criteria.Offset/RecordingsPerPage)))
		} else if criteria.SortMode == infinity.SortModeOldest {
			page = uint64(math.Floor(float64(criteria.Offset/RecordingsPerPage))) + 1
		}

		// Cannot retrieve a negative page, but since it's uint it becomes large number instead
		if page > cache.pageCount {
			page = cache.pageCount
		}

		log := cache.app.Logger.WithFields(logrus.Fields{"page": page, "operation": "injectPageCacheEntry"})
		log.Debug("Searching page cache")

		pageCacheEntry, ok := cache.cache[page]
		if ok {
			pageCacheEntry.Mutex.RLock()
			defer pageCacheEntry.Mutex.RUnlock()

			// Handle special case when we requested the newest page and it's not complete
			if !pageCacheEntry.isFull() {
				secondCacheEntry, ok := cache.cache[page-1]
				// If cache is not properly generated we might not have 90 entries cached, just don't do anything then
				if !ok {
					return false
				}

				secondCacheEntry.Mutex.RLock()
				defer secondCacheEntry.Mutex.RUnlock()

				// Use data from two latest cache entries
				criteria.CreatedAfter = secondCacheEntry.StartTime
				criteria.CreatedBefore = pageCacheEntry.EndTime
				criteria.Offset = 0
				return true
			}

			criteria.CreatedAfter = pageCacheEntry.StartTime
			criteria.CreatedBefore = pageCacheEntry.EndTime
			criteria.Offset = 0
			return true
		}

		log.Warn("Page cache miss")
	}

	return false
}

func (cache *repositoryCache) cacheUsableWithCriteria(criteria *infinity.RecordingRepositoryCriteria) bool {
	// Doesn't work with any kind of filtering
	if criteria.PerformerId != uuid.Nil ||
		criteria.FavoritesOnly ||
		criteria.OriginSection != "" ||
		criteria.OriginService != "" ||
		criteria.StageName != "" {
		return false
	}

	// Newest to oldest sort mode
	if criteria.SortMode == "" || criteria.SortMode == infinity.SortModeOldest {
		return true
	} else {
		// Can't use this cache for most popular / most viewed
		return false
	}
}

func (cache *repositoryCache) addNewRecordingToPageCache(recording *infinity.Recording, performer *infinity.Performer) {
	log := cache.app.Logger.WithField("operation", "addNewRecordingToPageCache")

	cache.mtx.RLock()
	// Cache is not built yet
	if !cache.cacheBuilt {
		cache.mtx.RUnlock()
		log.Warn("Unable to update page cache since cache is not built yet")
		return
	}

	cache.mtx.RUnlock()
	cache.mtx.Lock()
	defer cache.mtx.Unlock()

	lastCacheEntry, ok := cache.cache[cache.pageCount]
	if !ok {
		cache.app.Logger.WithField("page", cache.pageCount).Error("Tried to access undefined cache entry")
		//return
	}

	lastCacheEntry.Mutex.Lock()
	defer lastCacheEntry.Mutex.Unlock()

	if !lastCacheEntry.isFull() {
		log.WithField("pageCount", cache.pageCount).Debug("Last cache entry is not full, adding recording")

		// Track origin section/service
		if _, ok := lastCacheEntry.OriginServiceCount[performer.OriginService]; !ok {
			lastCacheEntry.OriginServiceCount[performer.OriginService] = 1
		} else {
			lastCacheEntry.OriginServiceCount[performer.OriginService] += 1
		}

		if performer.OriginSection.Valid {
			if _, ok := lastCacheEntry.OriginSectionCount[performer.OriginSection.String]; !ok {
				lastCacheEntry.OriginSectionCount[performer.OriginSection.String] = 1
			} else {
				lastCacheEntry.OriginSectionCount[performer.OriginSection.String] += 1
			}
		}

		lastCacheEntry.RecordingCount++
		lastCacheEntry.EndTime = recording.CreatedAt

		cache.updateRedisEntryForCacheEntry(cache.pageCount, lastCacheEntry)

	} else {
		log.WithField("pageCount", cache.pageCount).Debug("Last cache entry is full, adding new cache entry")

		newCacheEntry := &cacheEntry{
			StartTime:          recording.CreatedAt,
			EndTime:            recording.CreatedAt,
			OriginSectionCount: map[string]uint{},
			OriginServiceCount: map[string]uint{
				performer.OriginService: 1,
			},
			RecordingCount: 1,
		}

		if performer.OriginSection.Valid {
			newCacheEntry.OriginSectionCount[performer.OriginSection.String] = 1
		}

		cache.updateRedisEntryForCacheEntry(cache.pageCount, lastCacheEntry)
		cache.updateRedisPageCount(cache.pageCount + 1)

		cache.pageCount++
		cache.cache[cache.pageCount] = newCacheEntry
	}
}

func (cache *repositoryCache) MarkAsValidated(ids []uuid.UUID) error {
	return cache.repository.MarkAsValidated(ids)
}

func (cache *repositoryCache) GetByUuid(id uuid.UUID) (*infinity.Recording, error) {
	return cache.repository.GetByUuid(id)
}

func (cache *repositoryCache) GetUuidByExternalId(externalId uint64) (uuid.UUID, error) {
	return cache.repository.GetUuidByExternalId(externalId)
}

func (cache *repositoryCache) GetByUuidWithUserContext(recordingId uuid.UUID, userId uuid.UUID) (*infinity.RecordingWithUserData, error) {
	return cache.repository.GetByUuidWithUserContext(recordingId, userId)
}

func (cache *repositoryCache) GetBySlug(slug string) (*infinity.Recording, error) {
	return cache.repository.GetBySlug(slug)
}

func (cache *repositoryCache) GetBySlugWithUserContext(slug string, user uuid.UUID) (*infinity.RecordingWithUserData, error) {
	return cache.repository.GetBySlugWithUserContext(slug, user)
}

func (cache *repositoryCache) GetAllMissingSlug() ([]infinity.Recording, error) {
	return cache.repository.GetAllMissingSlug()
}

func (cache *repositoryCache) AddView(recording *infinity.Recording, userHash string) (bool, error) {
	return cache.repository.AddView(recording, userHash)
}

func (cache *repositoryCache) Create(recording *infinity.Recording, performer *infinity.Performer) error {
	log := cache.app.Logger.WithField("operation", "createRecording")
	log.Info("Creating recording")

	if err := cache.repository.Create(recording, performer); err != nil {
		log.WithError(err).Error("Failed to create recording")
		return err
	}

	log.WithField("cacheBuilt", cache.cacheBuilt).Debug("Checking if cache is built")

	// Update page cache in another go routine
	go cache.addNewRecordingToPageCache(recording, performer)

	return nil
}

func (cache *repositoryCache) GetByExternalId(id uint64) (*infinity.Recording, error) {
	return cache.repository.GetByExternalId(id)
}

func (cache *repositoryCache) CanRemove(id uuid.UUID) error {
	return cache.repository.CanRemove(id)
}

func (cache *repositoryCache) Remove(id uuid.UUID) (bool, error) {
	return cache.repository.Remove(id)
}

func (cache *repositoryCache) Update(recording *infinity.Recording) error {
	return cache.repository.Update(recording)
}

func (cache *repositoryCache) GetRecordingsCreatedBefore(time time.Time, limit int) ([]infinity.Recording, error) {
	return cache.repository.GetRecordingsCreatedBefore(time, limit)
}

func (cache *repositoryCache) Matching(criteria *infinity.RecordingRepositoryCriteria) ([]infinity.Recording, int, error) {
	injectedCachedData := cache.injectPageCacheEntry(criteria)

	recordings, count, err := cache.repository.Matching(criteria, !injectedCachedData)
	if err != nil || !injectedCachedData {
		return recordings, count, err
	}

	cache.mtx.RLock()

	lastCacheEntry, ok := cache.cache[cache.pageCount]
	if !ok {
		cache.mtx.RUnlock()
		return recordings, count, errors.New("Could not locate last cache entry")
	}

	cache.mtx.RUnlock()
	lastCacheEntry.Mutex.RLock()
	defer lastCacheEntry.Mutex.RUnlock()

	return recordings, int(uint(RecordingsPerPage*(cache.pageCount-1)) + lastCacheEntry.RecordingCount), nil
}

func (cache *repositoryCache) MatchingWithUserContext(criteria *infinity.RecordingRepositoryCriteria, user uuid.UUID) ([]infinity.RecordingWithUserData, int, error) {
	injectedCachedData := cache.injectPageCacheEntry(criteria)

	recordings, count, err := cache.repository.MatchingWithUserContext(criteria, user, !injectedCachedData)
	if err != nil || !injectedCachedData {
		return recordings, count, err
	}

	cache.mtx.RLock()

	lastCacheEntry, ok := cache.cache[cache.pageCount]
	if !ok {
		cache.mtx.RUnlock()
		return recordings, count, errors.New("Could not locate last cache entry")
	}

	cache.mtx.RUnlock()
	lastCacheEntry.Mutex.RLock()
	defer lastCacheEntry.Mutex.RUnlock()

	return recordings, int(uint(RecordingsPerPage*(cache.pageCount-1)) + lastCacheEntry.RecordingCount), nil
}
