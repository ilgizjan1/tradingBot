package main

import (
	"fmt"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"trade-bot/configs"
	"trade-bot/pkg/client/app"
	"trade-bot/pkg/client/service"
	"trade-bot/pkg/telegramBot"
)

var (
	ErrReadConfig                = errors.New("read config")
	ErrUnableToCreateTelegramBot = errors.New("unable to create telegram bot")
	ErrUnableToCreateClient      = errors.New("unable to create client")
	ErrSetupBot                  = errors.New("setup bot")
)

func setBot(configuration configs.TelegramBotConfiguration) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(configuration.APIToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrSetupBot, err)
	}

	log.Infof("Authorized on account: %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(configuration.WebhookURL))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrSetupBot, err)
	}

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal("server stopped: ", err)
		}
	}()

	return bot, nil
}

func main() {
	config, err := initConfig()
	if err != nil {
		log.Fatalf("init config: %s", err)
	}

	bot, err := setBot(config.Telegram)
	if err != nil {
		log.Fatalf("%s: %s", ErrUnableToCreateTelegramBot, err)
	}

	client, err := app.NewClient(config.Client)
	if err != nil {
		log.Fatalf("%s: %s", ErrUnableToCreateClient, err)
	}
	s := service.NewService(client)

	botman := telegramBot.NewBotMan(bot, s)
	botman.ServeTelegram()
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

	var c configs.Configuration
	err := viper.Unmarshal(&c)
	return c, err
}
