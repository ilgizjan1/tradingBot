package redisRepo

import (
	"strconv"
	"testing"
	"trade-bot/pkg/utils"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestJWTRedis_CreateJWT(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mr.Close()

	c := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	r := NewJWTRedis(c)

	type args struct {
		userID int
		td     utils.TokenDetails
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				userID: 1,
				td: utils.TokenDetails{
					AccessToken: "token",
					AccessUUID:  "accessUUID",
					AtExpires:   1,
				},
			},
			want: "token",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := r.CreateJWT(test.args.userID, test.args.td)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
			mr.FlushAll()
		})
	}
}

func TestJWTRedis_GetJWTUserID(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mr.Close()

	c := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	r := NewJWTRedis(c)

	tests := []struct {
		name    string
		prepare func()
		ad      utils.AccessDetails
		want    int
		wantErr bool
	}{
		{
			name: "OK",
			prepare: func() {
				userID := 1
				td := utils.TokenDetails{
					AccessToken: "token",
					AccessUUID:  "accessUUID",
				}

				errAccess := r.client.Set(context.Background(), td.AccessUUID, strconv.Itoa(userID), 0).Err()
				if errAccess != nil {
					t.Errorf("unable to set value in redis: (%v)", errAccess)
				}
			},
			ad: utils.AccessDetails{
				AccessUUID: "accessUUID",
			},
			want: 1,
		},
		{
			name:    "Key does not exist",
			prepare: func() {},
			ad: utils.AccessDetails{
				AccessUUID: "accessUUID",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.prepare()

			got, err := r.GetJWTUserID(test.ad)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
			mr.FlushAll()
		})
	}
}

func TestJWTRedis_DeleteJWT(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mr.Close()

	c := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	r := NewJWTRedis(c)

	tests := []struct {
		name    string
		prepare func()
		ad      utils.AccessDetails
		wantErr bool
	}{
		{
			name: "OK",
			prepare: func() {
				userID := 1
				td := utils.TokenDetails{
					AccessToken: "token",
					AccessUUID:  "accessUUID",
				}

				errAccess := r.client.Set(context.Background(), td.AccessUUID, strconv.Itoa(userID), 0).Err()
				if errAccess != nil {
					t.Errorf("unable to set value in redis: (%v)", errAccess)
				}
			},
			ad: utils.AccessDetails{
				AccessUUID: "accessUUID",
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.prepare()

			err := r.DeleteJWT(test.ad)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mr.FlushAll()
		})
	}
}
