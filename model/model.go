package model


import (
	"time"
	"database/sql"
	"simpledex/utils"
)

// email validation??
type User struct {
	Id 			  int64 `db:"id"`
	Email         string `db:"email"`
	Password      string `db:"password"`
	Referral_link string `db:"referral_link"`
	ReferredBy string `db:"referred_by"`
	Created_at    time.Time `db:"created_at"`
	Is_confirmed  bool `db:"is_confirmed"`
	Activation_link string `db:"activation_link"`
}


type db interface {
	SelectUsers() ([]*User, error)
	SelectOneUser(email string) ([]*User, error)
	GetUserById(id int64) (*User, error)
	GetUserByLink(activationLink string) (*User, error)
	UpdateActivateUser(id int64) (sql.Result, error)
	InsertUser(
		email string,
		hashedPassword []byte,
		isConfirmed bool,
		referralLink string,
		referredBy string,
		activationLink string) (sql.Result, error)
}

type Model struct {
	db
}

func New(db db) *Model {
	return &Model{
		db: db,
	}
}

func (m *Model) Users() ([]*User, error) {
	return m.SelectUsers()
}

func (m *Model) OneUser(email string) ([]*User, error) {
	return m.SelectOneUser(email)
}

func (m *Model) GetUserById(id int64) ([]*User, error) {
	return m.GetUserById(id)
}

func (m *Model) GetUserByActivationLink(activationLink string) (*User, error) {
	return m.GetUserByLink(activationLink)
}

func (m *Model) ActivateUser(id int64) (sql.Result, error) {
	return m.UpdateActivateUser(id)
}

func (m *Model) CreateUser(email string, password string, referredBy string, isConfirmed bool) (sql.Result, error) {
	referraLink := utils.RandStringRunes(12)
	activationLink := utils.RandStringRunes(15)
	hashedPassword := utils.EncryptPassword(password)
	return m.InsertUser(email, hashedPassword, isConfirmed, referraLink, referredBy, activationLink)
}
