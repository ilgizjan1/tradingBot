package telegramBot

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"

	"trade-bot/pkg/client/models"
	"trade-bot/pkg/client/service"
	"trade-bot/pkg/telegramBot/utils"
)

var (
	ErrCouldNotSendMessage            = errors.New("could not send message")
	ErrExitFromSignUpInput            = errors.New("exited from sign up input")
	ErrExitFromSignInInput            = errors.New("exited from sign in input")
	ErrExitFromSendOrderInput         = errors.New("exited from send order input")
	ErrExitFromStartTradingCommand    = errors.New("exited from start trading input")
	ErrUnableToReadFromUpdatesChannel = errors.New("unable to read from updates channel")
	ErrUserAlreadyLoggedIn            = errors.New("user already logged in")
)

const (
	startCommand                = "/start"
	helpCommand                 = "/help"
	signUpCommand               = "/sign_up"
	exitFromSignUpCommand       = "/exit_from_sign_up"
	signInCommand               = "/sign_in"
	exitFromSignInCommand       = "/exit_from_sign_in"
	sendOrderCommand            = "/send_order"
	exitFromSendOrderCommand    = "/exit_from_send_order"
	startTradingCommand         = "/start_trading"
	exitFromStartTradingCommand = "/exit_from_start_trading"
	getUserOrdersCommand        = "/get_user_orders"
	logoutCommand               = "/logout"
)

type BotMan struct {
	bot              *tgbotapi.BotAPI
	tradeBotServices *service.Service
	usersJWT         map[string]string
}

func NewBotMan(bot *tgbotapi.BotAPI, tradeBotServices *service.Service) *BotMan {
	return &BotMan{bot: bot, tradeBotServices: tradeBotServices, usersJWT: map[string]string{}}
}

func (b *BotMan) ServeTelegram() {
	updates := b.bot.ListenForWebhook("/")

	for update := range updates {
		if update.Message != nil {
			if !update.Message.IsCommand() {
				continue
			}

			log.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)

			chatID := update.Message.Chat.ID

			switch update.Message.Text {
			case startCommand:
				message := tgbotapi.NewMessage(chatID, utils.StartMessage)
				message.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				b.sendMessage(chatID, message)

			case helpCommand:
				message := tgbotapi.NewMessage(chatID, utils.HelpMessage)
				message.ReplyToMessageID = update.Message.MessageID
				b.sendMessage(chatID, message)

			case signUpCommand:
				if _, ok := b.usersJWT[update.Message.From.UserName]; ok {
					log.Warn(ErrUserAlreadyLoggedIn)
					errMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", utils.SignUpErrMessage, ErrUserAlreadyLoggedIn))
					b.sendMessage(chatID, errMessage)
					continue
				}

				message := tgbotapi.NewMessage(chatID, utils.SignUpMessage)
				message.ReplyToMessageID = update.Message.MessageID
				b.sendMessage(chatID, message)

				if err := b.executeSignUp(updates); err != nil {
					log.Warn(err)
					errMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", utils.SignUpErrMessage, err.Error()))
					b.sendMessage(chatID, errMessage)
				} else {
					successMessage := tgbotapi.NewMessage(chatID, utils.SignUpSuccessMessage)
					b.sendMessage(chatID, successMessage)
				}

			case signInCommand:
				if _, ok := b.usersJWT[update.Message.From.UserName]; ok {
					log.Warn(ErrUserAlreadyLoggedIn)
					errMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", utils.SignInErrMessage, ErrUserAlreadyLoggedIn))
					b.sendMessage(chatID, errMessage)
					continue
				}

				message := tgbotapi.NewMessage(chatID, utils.SignInMessage)
				message.ReplyToMessageID = update.Message.MessageID
				b.sendMessage(chatID, message)

				token, err := b.executeSignIn(updates)
				if err != nil {
					log.Warn(err)
					errMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", utils.SignInErrMessage, err.Error()))
					b.sendMessage(chatID, errMessage)
					continue
				}

				b.usersJWT[update.Message.From.UserName] = token
				successMessage := tgbotapi.NewMessage(chatID, utils.SignInSuccessMessage)
				b.sendMessage(chatID, successMessage)

			case logoutCommand:
				token, err := b.userIdentity(update.Message.From.UserName)
				if err != nil {
					log.Warn(err)
					errMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", utils.LogoutErrMessage, err.Error()))
					b.sendMessage(chatID, errMessage)
					continue
				}

				_, err = b.tradeBotServices.Authorization.Logout(models.LogoutInput{JWTToken: token})
				if err != nil {
					log.Warn(err)
					errMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", utils.LogoutErrMessage, err.Error()))
					b.sendMessage(chatID, errMessage)
					continue
				}

				delete(b.usersJWT, update.Message.From.UserName)
				successMessage := tgbotapi.NewMessage(chatID, utils.LogoutSuccessMessage)
				b.sendMessage(chatID, successMessage)

			case sendOrderCommand:
				token, err := b.userIdentity(update.Message.From.UserName)
				if err != nil {
					log.Warn(err)
					errMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", utils.SendOrderErrMessage, err.Error()))
					b.sendMessage(chatID, errMessage)
					continue
				}

				message := tgbotapi.NewMessage(chatID, utils.SendOrderMessage)
				b.sendMessage(chatID, message)

				resp, err := b.executeSendOrder(updates, token)
				if err != nil {
					log.Warn(err)
					errMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", utils.SendOrderErrMessage, err.Error()))
					b.sendMessage(chatID, errMessage)
					continue
				}

				successMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s\n%s", utils.SendOrderSuccessMessage, resp.String()))
				b.sendMessage(chatID, successMessage)

			case getUserOrdersCommand:
				token, err := b.userIdentity(update.Message.From.UserName)
				if err != nil {
					log.Warn(err)
					errMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", utils.GetUserOrdersErrMessage, err.Error()))
					b.sendMessage(chatID, errMessage)
					continue
				}

				resp, err := b.executeGetUserOrders(token)
				if err != nil {
					log.Warn(err)
					errMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", utils.GetUserOrdersErrMessage, err.Error()))
					b.sendMessage(chatID, errMessage)
					continue
				}

				successMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s\n%s", utils.SendOrderSuccessMessage, resp.String()))
				b.sendMessage(chatID, successMessage)

			case startTradingCommand:
				token, err := b.userIdentity(update.Message.From.UserName)
				if err != nil {
					log.Warn(err)
					errMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", utils.StartTradingErrMessage, err.Error()))
					b.sendMessage(chatID, errMessage)
					continue
				}

				message := tgbotapi.NewMessage(chatID, utils.StartTradingMessage)
				b.sendMessage(chatID, message)

				if err := b.executeStartTrading(chatID, updates, token); err != nil {
					log.Warn(err)
					errMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %s", utils.StartTradingErrMessage, err.Error()))
					b.sendMessage(chatID, errMessage)
					continue
				}

				message = tgbotapi.NewMessage(chatID, utils.StartTradingWillNotifyMessage)
				b.sendMessage(chatID, message)

			default:
				message := tgbotapi.NewMessage(chatID, utils.InvalidCommandMessage)
				b.sendMessage(chatID, message)
			}
		}
	}
}

func (b *BotMan) userIdentity(username string) (string, error) {
	val, ok := b.usersJWT[username]
	if !ok {
		return "", fmt.Errorf("user not logged in")
	}
	return val, nil
}

func (b *BotMan) sendMessage(chatID int64, message tgbotapi.MessageConfig) {
	if _, err := b.bot.Send(message); err != nil {
		log.Warnf("%s: [chatID] - %d", ErrCouldNotSendMessage, chatID)
	}
}

func (b *BotMan) executeGetUserOrders(token string) (models.GetUserOrdersResponse, error) {
	input := models.GetUserOrdersInput{JWTToken: token}

	resp, err := b.tradeBotServices.GetUserOrders(input)
	if err != nil {
		return models.GetUserOrdersResponse{}, err
	}

	return resp, nil
}

func (b *BotMan) executeStartTrading(chatID int64, updates tgbotapi.UpdatesChannel, token string) error {
	input, err := b.getStartTradingInput(updates)
	if err != nil {
		return err
	}

	input.JWTToken = token

	startTradingResp, errCh, err := b.tradeBotServices.OrdersManager.StartTrading(input)
	if err != nil {
		return err
	}

	go func() {
		for err := range errCh {
			log.Error(err)
		}
	}()

	go func(chatID int64) {
		for val := range startTradingResp {
			message := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s\nOrder: %s", utils.StartTradingSuccessMessage, val.String()))
			b.sendMessage(chatID, message)
		}
	}(chatID)

	return nil
}

func (b *BotMan) getStartTradingInput(updates tgbotapi.UpdatesChannel) (models.StartTradingInput, error) {
	for update := range updates {
		if update.Message == nil {
			return models.StartTradingInput{}, nil
		}

		switch update.Message.Text {
		case exitFromStartTradingCommand:
			return models.StartTradingInput{}, ErrExitFromStartTradingCommand
		default:
			inputValues := strings.FieldsFunc(update.Message.Text, split)
			if len(inputValues) != 5 {
				return models.StartTradingInput{}, fmt.Errorf("invalid count of arguments")
			}
			if inputValues[1] != "buy" && inputValues[1] != "sell" {
				return models.StartTradingInput{}, fmt.Errorf("invalid strat trading Side argument")
			}
			amount, err := strconv.ParseUint(inputValues[2], 10, 64)
			if err != nil {
				return models.StartTradingInput{}, fmt.Errorf("invalid start trading Size argument")
			}
			stopLoss, err := strconv.ParseFloat(inputValues[3], 64)
			if err != nil {
				return models.StartTradingInput{}, fmt.Errorf("invalid start trading Stop loss argument")
			}
			takeProfit, err := strconv.ParseFloat(inputValues[4], 64)
			if err != nil {
				return models.StartTradingInput{}, fmt.Errorf("invalid start trading Take profit argument")
			}

			return models.StartTradingInput{
				Event: "start_trading",
				TradingDetails: models.StartTradingDetails{
					SendOrderInput: models.SendOrderInput{
						OrderType: "mkt",
						Symbol:    inputValues[0],
						Side:      inputValues[1],
						Size:      uint(amount),
					},
					StopLossBorder:   uint(stopLoss),
					TakeProfitBorder: uint(takeProfit),
				},
			}, nil
		}
	}

	return models.StartTradingInput{}, ErrUnableToReadFromUpdatesChannel
}

func (b *BotMan) executeSendOrder(updates tgbotapi.UpdatesChannel, token string) (models.SendOrderResponse, error) {
	input, err := b.getSendOrderInput(updates)
	if err != nil {
		return models.SendOrderResponse{}, err
	}
	input.JWTToken = token

	resp, err := b.tradeBotServices.OrdersManager.SendOrder(input)
	if err != nil {
		return models.SendOrderResponse{}, err
	}

	return resp, nil
}

func (b *BotMan) getSendOrderInput(updates tgbotapi.UpdatesChannel) (models.SendOrderInput, error) {
	for update := range updates {
		if update.Message == nil {
			return models.SendOrderInput{}, nil
		}

		switch update.Message.Text {
		case exitFromSendOrderCommand:
			return models.SendOrderInput{}, ErrExitFromSendOrderInput
		default:
			inputValues := strings.FieldsFunc(update.Message.Text, split)
			if len(inputValues) != 3 {
				return models.SendOrderInput{}, fmt.Errorf("invalid count of arguments")
			}
			if inputValues[1] != "buy" && inputValues[1] != "sell" {
				return models.SendOrderInput{}, fmt.Errorf("invalid send order Side argument")
			}
			amount, err := strconv.ParseUint(inputValues[2], 10, 64)
			if err != nil {
				return models.SendOrderInput{}, fmt.Errorf("invalid send order Size argument")
			}
			return models.SendOrderInput{
				OrderType: "mkt",
				Symbol:    inputValues[0],
				Side:      inputValues[1],
				Size:      uint(amount),
			}, nil
		}
	}

	return models.SendOrderInput{}, ErrUnableToReadFromUpdatesChannel
}

func (b *BotMan) executeSignIn(updates tgbotapi.UpdatesChannel) (string, error) {
	input, err := b.getSignInInput(updates)
	if err != nil {
		return "", err
	}

	resp, err := b.tradeBotServices.Authorization.SignIn(input)
	if err != nil {
		return "", err
	}
	return resp.AccessToken, nil
}

func (b *BotMan) getSignInInput(updates tgbotapi.UpdatesChannel) (models.SignInInput, error) {
	for update := range updates {
		if update.Message == nil {
			return models.SignInInput{}, nil
		}

		switch update.Message.Text {
		case exitFromSignInCommand:
			return models.SignInInput{}, ErrExitFromSignInInput
		default:
			inputValues := strings.FieldsFunc(update.Message.Text, split)
			if len(inputValues) != 2 {
				return models.SignInInput{}, fmt.Errorf("invalid count of arguments")
			}
			return models.SignInInput{
				Username: inputValues[0],
				Password: inputValues[1],
			}, nil
		}
	}

	return models.SignInInput{}, ErrUnableToReadFromUpdatesChannel
}

func (b *BotMan) executeSignUp(updates tgbotapi.UpdatesChannel) error {
	input, err := b.getSignUpInput(updates)
	if err != nil {
		return err
	}

	_, err = b.tradeBotServices.Authorization.SignUp(input)
	return err
}

func (b *BotMan) getSignUpInput(updates tgbotapi.UpdatesChannel) (models.SignUpInput, error) {
	for update := range updates {
		if update.Message == nil {
			return models.SignUpInput{}, nil
		}

		switch update.Message.Text {
		case exitFromSignUpCommand:
			return models.SignUpInput{}, ErrExitFromSignUpInput
		default:
			inputValues := strings.FieldsFunc(update.Message.Text, split)
			if len(inputValues) != 5 {
				return models.SignUpInput{}, fmt.Errorf("invalid count of arguments")
			}
			return models.SignUpInput{
				Name:          inputValues[0],
				Username:      inputValues[1],
				Password:      inputValues[2],
				PublicAPIKey:  inputValues[3],
				PrivateAPIKey: inputValues[4],
			}, nil
		}
	}

	return models.SignUpInput{}, ErrUnableToReadFromUpdatesChannel
}

func split(r rune) bool {
	return r == ' ' || r == '\n'
}
