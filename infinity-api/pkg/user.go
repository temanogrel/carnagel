package infinity

import (
	"errors"
	"time"

	"git.misc.vee.bz/carnagel/go-ecosystem/domain"
	"github.com/satori/go.uuid"
)

var (
	UserNotFoundErr                      = errors.New("The user was not found")
	InvalidCredentialsErr                = errors.New("Invalid credentials were provided")
	FailedToParseSessionDataFromRedisErr = errors.New("The session data in redis was not valid json")
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
	RoleGuest = "guest"
)

type User struct {
	TableName struct{} `sql:"users,alias:u" json:"-"`

	Uuid     uuid.UUID `sql:",pk" json:"uuid"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Password string    `json:"-"`
	Role     string    `json:"role"`

	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	PaymentPlan             *PaymentPlan `json:"-"`
	PaymentPlanUuid         uuid.UUID    `json:"paymentPlanUuid"`
	PaymentPlanSubscribedAt *time.Time   `json:"paymentPlanSubscribedAt"`
	PaymentPlanEndsAt       *time.Time   `json:"paymentPlanEndsAt"`
}

type UserLike struct {
	TableName     struct{}  `sql:"user_recording_likes,alias:url" json:"-"`
	UserUuid      uuid.UUID `sql:",pk"`
	RecordingUuid uuid.UUID `sql:",pk"`
	CreatedAt     time.Time
}

type UserFavorite struct {
	TableName     struct{}  `sql:"user_recording_favorites,alias:urf" json:"-"`
	UserUuid      uuid.UUID `sql:",pk"`
	RecordingUuid uuid.UUID `sql:",pk"`
}

type UserRepositoryCriteria struct {
	CurrentlyPremium bool
	ExPremium        bool
	NeverPremium     bool

	CreatedAfter time.Time

	Offset int
	Limit  int

	Sorting map[string]string
}

func NewUserRepositoryCriteria() *UserRepositoryCriteria {
	return &UserRepositoryCriteria{
		Limit:   20,
		Sorting: map[string]string{"created_at": "desc"},
	}
}

type UserRepository interface {
	GetByUuid(uuid.UUID) (*User, error)
	GetByEmail(string) (*User, error)
	GetByUsername(string) (*User, error)
	GetByUsernameOrEmail(string) (*User, error)
	GetUsersWithExpiredPaymentPlan() ([]*User, error)
	GetUsersWithExpiringPaymentPlan(daysBeforeExpiration uint) ([]*User, error)

	Matching(*UserRepositoryCriteria) ([]User, int, error)

	Create(*User) error
	Update(*User) error
	RemoveById(uuid.UUID) (bool, error)
	ToggleLike(uuid.UUID, uuid.UUID) (bool, error)
	ToggleFavorite(uuid.UUID, uuid.UUID) (bool, error)
}

type UserSessionService interface {
	New(remoteAddr string) (string, error)
	Renew(uuid.UUID, uuid.UUID, bool) (string, error)
	Authenticate(username, password, remoteAddr string) (string, *User, error)
	ParseToken(token string) (*ecosystem.JwtClaims, error)
}
