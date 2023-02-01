package postgresRepo

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"testing"
	"trade-bot/internal/pkg/models"
)

func TestAuthPostgres_CreateUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewAuthPostgres(sqlxDB)

	tests := []struct {
		name    string
		mock    func()
		input   models.User
		want    int
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("name", "username", "password", "key", "key").WillReturnRows(rows)
			},
			input: models.User{
				Name:          "name",
				Username:      "username",
				Password:      "password",
				PublicAPIKey:  "key",
				PrivateAPIKey: "key",
			},
			want: 1,
		},
		{
			name: "Empty Fields",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("name", "username", "", "key", "key").WillReturnRows(rows)
			},
			input: models.User{
				Name:          "name",
				Username:      "username",
				Password:      "",
				PublicAPIKey:  "key",
				PrivateAPIKey: "key",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.CreateUser(test.input)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthPostgres_GetUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewAuthPostgres(sqlxDB)

	tests := []struct {
		name     string
		mock     func()
		username string
		want     models.User
		wantErr  bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "password_hash",
					"public_api_key", "private_api_key"}).
					AddRow(1, "name", "username", "password", "key", "key")
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("username").WillReturnRows(rows)
			},
			username: "username",
			want: models.User{
				ID:            1,
				Name:          "name",
				Username:      "username",
				Password:      "password",
				PublicAPIKey:  "key",
				PrivateAPIKey: "key",
			},
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "password_hash",
					"public_api_key", "private_api_key"})
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("username").WillReturnRows(rows)
			},
			username: "username",
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.GetUser(test.username)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthPostgres_GetUserAPIKeys(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	r := NewAuthPostgres(sqlxDB)

	type wantArgs struct {
		publicAPIKey  string
		privateAPIKey string
	}

	tests := []struct {
		name    string
		mock    func()
		userID  int
		want    wantArgs
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "password_hash",
					"public_api_key", "private_api_key"}).
					AddRow(1, "name", "username", "password", "key", "key")
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(1).WillReturnRows(rows)
			},
			userID: 1,
			want: wantArgs{
				publicAPIKey:  "key",
				privateAPIKey: "key",
			},
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "password_hash",
					"public_api_key", "private_api_key"})
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs(1).WillReturnRows(rows)
			},
			userID:  1,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			publicAPIKeyGot, privateAPIKeyGot, err := r.GetUserAPIKeys(test.userID)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want.privateAPIKey, publicAPIKeyGot)
				assert.Equal(t, test.want.privateAPIKey, privateAPIKeyGot)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
