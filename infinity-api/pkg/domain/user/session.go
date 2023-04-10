package user

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"encoding/json"
	ecosystem "git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type dailySessionsForRemoteAddr struct {
	GuestSession uuid.UUID `json:"guestSession"`

	// Map user uuid to a session
	Sessions map[uuid.UUID]uuid.UUID `json:"sessions"`

	// Track which session ids are black listed
	BlackListedSessions map[uuid.UUID]bool `json:"blackListedSession"`
}

type sessionService struct {
	app *infinity.Application
	log logrus.FieldLogger
}

func NewSessionService(app *infinity.Application) infinity.UserSessionService {
	return &sessionService{
		app: app,
		log: app.Logger.WithField("component", "user_session"),
	}
}

// New generates a new sessionService based on the remote addr
func (auth *sessionService) New(remoteAddr string) (string, error) {
	return auth.getTokenForSessionOfRemoteAddr(remoteAddr, nil)
}

func (auth *sessionService) Renew(userUuid uuid.UUID, session uuid.UUID, outOfBasicPlans bool) (string, error) {

	user, err := auth.app.UserRepository.GetByUuid(userUuid)
	if err != nil {
		return "", err
	}

	return auth.createToken(session, user, outOfBasicPlans)
}

func (auth *sessionService) Authenticate(identity, password, remoteAddr string) (string, *infinity.User, error) {

	user, err := auth.app.UserRepository.GetByUsernameOrEmail(identity)
	if err == infinity.UserNotFoundErr {
		return "", nil, infinity.InvalidCredentialsErr
	}

	if err != nil {
		return "", nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return "", nil, infinity.InvalidCredentialsErr

		default:
			return "", nil, errors.Wrap(err, "Failed to compare passwords")
		}
	}

	token, err := auth.getTokenForSessionOfRemoteAddr(remoteAddr, user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (auth *sessionService) ParseToken(token string) (*ecosystem.JwtClaims, error) {
	claims := &ecosystem.JwtClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(auth.app.SessionSignKey), nil
	})

	if err != nil {
		// Check if error is because of validation expired
		if vErr, ok := err.(*jwt.ValidationError); ok && vErr.Errors&jwt.ValidationErrorExpired == jwt.ValidationErrorExpired {
			return nil, ecosystem.SessionExpiredErr
		} else {
			return nil, errors.Wrapf(err, "Failed to parse token")
		}
	}

	return claims, nil
}

func (auth *sessionService) getTokenForSessionOfRemoteAddr(remoteAddr string, user *infinity.User) (string, error) {

	if strings.Contains(remoteAddr, ":") {
		remoteAddr = strings.Split(remoteAddr, ":")[0]
	}

	// Generate a byte array checksum of the provided remoteAddrs
	checkSum := md5.Sum([]byte(fmt.Sprintf("%s", remoteAddr)))

	// Not exactly sure about the conversion here
	identifier := hex.EncodeToString(checkSum[:])

	data, err := auth.getSessionDataFromRedis(identifier)

	switch err {
	case nil:
		session := uuid.NewV4()
		var storedSession uuid.UUID
		var ok bool
		var sessionBlackListed bool

		if user != nil {
			if data.Sessions == nil {
				data.Sessions = map[uuid.UUID]uuid.UUID{}
			}

			if data.BlackListedSessions == nil {
				data.BlackListedSessions = map[uuid.UUID]bool{}
			}

			storedSession, ok = data.Sessions[user.Uuid]
			if !ok || storedSession == uuid.Nil {
				numberOfBasicPlansUsed, err := auth.getNumberOfBasicPlansUsedTodayForRemoteAddr(data)
				if err != nil {
					return "", err
				}

				basicPlan, err := auth.app.PaymentPlanRepository.GetBasicPlan()
				if err != nil {
					return "", err
				}

				if numberOfBasicPlansUsed >= infinity.NumberOfBasicPlansUsablePerDayForRemoteAddr && user.PaymentPlanUuid == basicPlan.Uuid {
					auth.log.WithField("session", session).Debug("Black listing session")
					auth.app.BandwidthConsumptionCollector.AddBlackListedSession(session)

					data.BlackListedSessions[session] = true
					sessionBlackListed = true
				}

				data.Sessions[user.Uuid] = session

				if err = auth.persistSessionDataToRedis(identifier, data); err != nil {
					return "", err
				}

			} else {
				sessionBlackListed, _ = data.BlackListedSessions[storedSession]
				session = storedSession
			}
		} else {
			if data.GuestSession != uuid.Nil {
				session = data.GuestSession
			} else {
				data.GuestSession = session

				if err = auth.persistSessionDataToRedis(identifier, data); err != nil {
					return "", err
				}
			}
		}

		return auth.createToken(session, user, sessionBlackListed)

	case infinity.FailedToParseSessionDataFromRedisErr:
		fallthrough
	case redis.Nil:
		auth.log.Debug("Did not find stored session data for remote address in redis")

		session := uuid.NewV4()

		data := &dailySessionsForRemoteAddr{
			Sessions:            map[uuid.UUID]uuid.UUID{},
			BlackListedSessions: map[uuid.UUID]bool{},
		}

		if user != nil {
			data.Sessions[user.Uuid] = session
		} else {
			data.GuestSession = session
		}

		if err = auth.persistSessionDataToRedis(identifier, data); err != nil {
			return "", err
		}

		return auth.createToken(session, user, false)

	default:
		auth.app.Logger.WithError(err).Error("Error when retrieving session data from redis")

		return "", errors.New("Unknown redis error")
	}
}

func (auth *sessionService) getNumberOfBasicPlansUsedTodayForRemoteAddr(data *dailySessionsForRemoteAddr) (uint, error) {
	var numberOfBasicPlansUsedByRemoteAddrToday uint

	basicPaymentPlan, err := auth.app.PaymentPlanRepository.GetBasicPlan()
	if err != nil {
		return 0, err
	}

	for userUuid := range data.Sessions {
		u, err := auth.app.UserRepository.GetByUuid(userUuid)
		if err != nil {
			return 0, err
		}

		if u.PaymentPlanUuid == basicPaymentPlan.Uuid {
			numberOfBasicPlansUsedByRemoteAddrToday++
		}
	}

	return numberOfBasicPlansUsedByRemoteAddrToday, nil
}

func (auth *sessionService) persistSessionDataToRedis(identifier string, data *dailySessionsForRemoteAddr) error {
	redisData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Midnight in UTC we reset the session
	midnight := time.Now().
		Truncate(time.Hour * 24).
		Add(time.Hour * 24).
		UTC()

	if err := auth.app.Redis.Set(identifier, redisData, time.Until(midnight)).Err(); err != nil {
		auth.app.Logger.WithError(err).Error("Failed to write to redis")

		return errors.Wrap(err, "Failed to cache session details for identifier in redis")
	}

	return nil
}

func (auth *sessionService) getSessionDataFromRedis(identifier string) (*dailySessionsForRemoteAddr, error) {
	// Check redis if the provided remoteAddr and user agent exist
	// in redis. This way it does not matter if they clear their browser cache/cookies
	redisData, err := auth.app.Redis.Get(identifier).Result()
	if err == nil {
		data := &dailySessionsForRemoteAddr{}
		if err = json.Unmarshal([]byte(redisData), data); err != nil {
			return nil, infinity.FailedToParseSessionDataFromRedisErr
		}

		return data, nil
	}

	return nil, err
}

func (auth *sessionService) createToken(session uuid.UUID, user *infinity.User, blackListedToday bool) (string, error) {

	if user == nil {
		paymentPlan, err := auth.app.PaymentPlanRepository.GetGuestPlan()
		if err != nil {
			return "", errors.Wrapf(err, "Failed to retrieve guest plan")
		}

		return jwt.NewWithClaims(auth.app.SessionSignMethod, ecosystem.JwtClaims{
			StandardClaims: jwt.StandardClaims{
				Issuer:    "camtube.co",
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(auth.app.SessionTokenDuration).Unix(),
			},

			Role:             infinity.RoleGuest,
			Session:          session,
			PaymentPlan:      paymentPlan.Uuid,
			BlackListedToday: false,
		}).SignedString(auth.app.SessionSignKey)
	}

	return jwt.NewWithClaims(auth.app.SessionSignMethod, ecosystem.JwtClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "camtube.co",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(auth.app.SessionTokenDuration).Unix(),
		},

		Role:             user.Role,
		User:             user.Uuid,
		Session:          session,
		PaymentPlan:      user.PaymentPlanUuid,
		BlackListedToday: blackListedToday,
	}).SignedString(auth.app.SessionSignKey)
}
