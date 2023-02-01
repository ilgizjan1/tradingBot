package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"trade-bot/configs"
	"trade-bot/internal/app"
	"trade-bot/internal/pkg/handler"
	"trade-bot/internal/pkg/repository"
	"trade-bot/internal/pkg/repository/postgresRepo"
	"trade-bot/internal/pkg/repository/redisRepo"
	"trade-bot/internal/pkg/service"
	"trade-bot/internal/pkg/tradeAlgorithm"
	"trade-bot/internal/pkg/web"
	"trade-bot/pkg/krakenFuturesSDK"
	"trade-bot/pkg/krakenFuturesWSSDK"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	ErrUnableToInitConfig           = errors.New("unable to init config files")
	ErrReadConfig                   = errors.New("read config")
	ErrRunServer                    = errors.New("run server")
	ErrUnableToConnectToDB          = errors.New("unable to connect to database")
	ErrUnableToConnectToJWTDB       = errors.New("unable to connect to jwt databased")
	ErrUnableToLoadEnvVariables     = errors.New("unable to load enviroment variables")
	ErrCouldNotShutdownServer       = errors.New("could not shut down server normally")
	ErrCouldNotCloseDBConnection    = errors.New("could not close db connection normally")
	ErrCouldNotCloseRedisConnection = errors.New("could not close redis connection normally")
)

const (
	publicAPIKey  = "PUBLIC_API_KEY"
	privateAPIKey = "PRIVATE_API_KEY"
)

// @title Trade-bot API
// @version 1.0
// @description API Server for Trade-bot Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	config, err := initConfig()
	if err != nil {
		log.Panicf("%s: %s", ErrUnableToInitConfig, err)
	}

	db, err := postgresRepo.NewPostgresDB(config.PostgreDatabase)
	if err != nil {
		log.Panicf("%s: %s", ErrUnableToConnectToDB, err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Panicf("%s: %s", ErrCouldNotCloseDBConnection, err)
		}
	}()

	redisClient, err := redisRepo.NewRedisClient(config.RedisDatabase)
	if err != nil {
		log.Panicf("%s: %s", ErrUnableToConnectToJWTDB, err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Panicf("%s: %s", ErrCouldNotCloseRedisConnection, err)
		}
	}()

	krakenAPI := krakenFuturesSDK.NewAPI(os.Getenv(publicAPIKey), os.Getenv(privateAPIKey), config.Kraken.APIURL)
	krakenWSAPI := krakenFuturesWSSDK.NewWSAPI(config.KrakenWS)

	repo := repository.NewRepository(db, redisClient)
	newWeb := web.NewWeb(krakenAPI, krakenWSAPI)
	newTrader := tradeAlgorithm.NewTradeAlgorithm(newWeb)

	validate := validator.New()
	upgrader := websocket.Upgrader{
		WriteBufferSize: config.Server.Websocket.WriteBufferSize,
		ReadBufferSize:  config.Server.Websocket.ReadBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return config.Server.Websocket.CheckOrigin
		},
	}

	services := service.NewService(repo, newWeb, newTrader)
	handlers := handler.NewHandler(services, validate, &upgrader)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	srv := new(app.Server)
	go func() {
		if err := srv.Run(config.Server.Port, handlers.InitRoutes()); err != nil && err != http.ErrServerClosed {
			log.Panicf("%s: %s", ErrRunServer, err)
		}
	}()

	log.Info("Trade bot server started")

	<-interrupt

	log.Info("interrupt signal caught")
	log.Info("Trade bot server shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Panicf("%s: %s", ErrCouldNotShutdownServer, err)
	}

	log.Info("Trade bot server shut down")
}

func initConfig() (configs.Configuration, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("configs")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal(fmt.Errorf("%s: %s", ErrReadConfig, err))
		}
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal(fmt.Errorf("%s: %s", ErrUnableToLoadEnvVariables, err))
	}

	var c configs.Configuration
	err := viper.Unmarshal(&c)
	c.PostgreDatabase.Password = os.Getenv("DB_PASSWORD")
	return c, err
}
