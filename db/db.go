package db

import (
	"simpledex/model"

	"database/sql"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type Config struct {
	ConnectString string
}

type pgDb struct {
	dbConn           *sqlx.DB
	sqlSelectUsers   *sqlx.Stmt
	sqlSelectUserByEmail *sqlx.Stmt
	sqlGetUserById *sqlx.Stmt
	sqlInsertUser *sqlx.Stmt
	sqlGetUserByActivationLink *sqlx.Stmt
	sqlActivateUser *sqlx.Stmt
}

func InitDb(cfg Config) (*pgDb, error) {
	if dbConn, err := sqlx.Connect("postgres", cfg.ConnectString); err != nil {
		return nil, err
	} else {
		p := &pgDb{dbConn: dbConn}
		if err := p.dbConn.Ping(); err != nil {
			return nil, err
		}
		if err := p.createTablesIfNotExist(); err != nil {
			return nil, err
		}
		if err := p.prepareSqlStatements(); err != nil {
			return nil, err
		}
		return p, nil
	}
}

func (p *pgDb) createTablesIfNotExist() error {
	createSQL := `
		CREATE TABLE IF NOT EXISTS users (
       		id SERIAL NOT NULL PRIMARY KEY,
       		email TEXT NOT NULL,
       		password TEXT NOT NULL,
       		created_at TIMESTAMP not null default CURRENT_TIMESTAMP,
       		referral_link TEXT NOT NULL,
       		referred_by TEXT,
       		activation_link TEXT NOT NULL,
       		is_confirmed BOOLEAN DEFAULT FALSE
       	);
	`
	if _, err := p.dbConn.Exec(createSQL); err != nil {
		return err
	}
	return nil
}

func (p *pgDb) prepareSqlStatements() (err error) {

	if p.sqlSelectUsers, err = p.dbConn.Preparex(
		"SELECT * FROM USERS",
	); err != nil {
		return err
	}
	if p.sqlSelectUserByEmail, err = p.dbConn.Preparex(
		"SELECT * FROM USERS  WHERE email = $1",
	); err != nil {
		return err
	}
	if p.sqlInsertUser, err = p.dbConn.Preparex(
		"INSERT INTO Users (email, password, is_confirmed, REFERRAL_LINK, referred_by, activation_link) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
	); err != nil {
		return err
	}
	if p.sqlGetUserById, err = p.dbConn.Preparex(
		"SELECT ID, EMAIL, CREATED_AT, PASSWORD, REFERRAL_LINK, IS_CONFIRMED FROM USERS WHERE id = $1",
	); err != nil {
		return err
	}
	if p.sqlGetUserByActivationLink, err = p.dbConn.Preparex(
		"SELECT * FROM USERS WHERE ACTIVATION_LINK = $1",
	); err != nil {
		return err
	}
	if p.sqlActivateUser, err = p.dbConn.Preparex(
		"UPDATE USERS SET IS_CONFIRMED = TRUE WHERE id = $1 ",
	); err != nil {
		return err
	}

	return nil
}

func (p *pgDb) SelectUsers() ([]*model.User, error) {
	user := make([]*model.User, 0)
	if err := p.sqlSelectUsers.Select(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (p *pgDb) SelectOneUser(email string) ([]*model.User, error) {
	user := make([]*model.User, 0)
	if err := p.sqlSelectUserByEmail.Select(&user, email); err != nil {
		return nil, err
	}
	return user, nil
}

func (p *pgDb) GetUserById(id int64) (*model.User, error) {
	var user *model.User
	if err := p.sqlGetUserById.Get(&user, id); err != nil {
		return nil, err
	}
	return user, nil
}

func (p *pgDb) GetUserByLink(activationLink string) (*model.User, error) {
	var user []*model.User
	if err := p.sqlGetUserByActivationLink.Select(&user, activationLink); err != nil {
		return nil, err
	}
	if len(user) > 0 {
		return user[0], nil
	}
	return nil, nil
}


func (p *pgDb) UpdateActivateUser(id int64) (sql.Result, error) {
	result, err := p.sqlActivateUser.Exec(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}



func (p *pgDb) InsertUser(email string, hashedPassword []byte, isConfirmed bool, referraLink string, referredBy string, activationLink string) (sql.Result, error) {
	result, err := p.sqlInsertUser.Exec(email, hashedPassword, isConfirmed, referraLink, referredBy, activationLink)
	if err != nil {
		return nil, err
	}
	return result, nil

}
