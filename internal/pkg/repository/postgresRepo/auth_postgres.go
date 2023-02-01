package postgresRepo

import (
	"github.com/jmoiron/sqlx"

	"trade-bot/internal/pkg/models"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

const insertUserQuery = `
	INSERT INTO users
    (name, username, password_hash, public_api_key, private_api_key) values ($1, $2, $3, $4, $5)
    RETURNING id`

func (r *AuthPostgres) CreateUser(user models.User) (int, error) {
	var id int
	row := r.db.QueryRow(insertUserQuery, user.Name, user.Username, user.Password, user.PublicAPIKey, user.PrivateAPIKey)
	err := row.Scan(&id)
	return id, err
}

const getUserQuery = "SELECT * FROM users WHERE username=$1"

func (r *AuthPostgres) GetUser(username string) (models.User, error) {
	var user models.User
	err := r.db.Get(&user, getUserQuery, username)
	return user, err
}

const getUserAPIKeysQuery = "SELECT * FROM users WHERE id=$1"

func (r *AuthPostgres) GetUserAPIKeys(userID int) (string, string, error) {
	var user models.User
	err := r.db.Get(&user, getUserAPIKeysQuery, userID)
	return user.PublicAPIKey, user.PrivateAPIKey, err
}
