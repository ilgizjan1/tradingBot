package utils

const StartMessage = `
Welcome to kraken futures trading bot! 📈

This bot supports trading on simple indicator calling "stop loss & take profit"

Here is bot commands:
	💁 /help - list descriptions for commands
`

const HelpMessage = `
Commands 📟:
	🔵 /help - list descriptions for commands
	🔵 /sign_up - register you in trading bot system
	🔵 /exit_from_sign_up - stop getting input data to register you in the bot
	🔵 /sign_in - login you in trading bot system and allow to trade on kraken futures
	🔵 /exit_from_sign_in - stop getting input data to login you in the bot
	🔵 /send_order - allow to send market order with symbol, side and amount arguments to kraken futures
	🔵 /exit_from_send_order - stop getting input data to send order to kraken futures
	🔵 /logout - logout you from trading bot system on every telegram device associated with your username
`

const InvalidCommandMessage = `
⛔ No such command
`

const SignUpMessage = `
🔳 Enter message in format:

Name
Username
Password
Your public api key
Your private api key

🔳 Example:

Ivan ivan password key key
`

const SignUpErrMessage = `
⛔ Unable to continue further execution of sign up due to
`

const SignUpSuccessMessage = `
✅ User successfully registered!
`

const SignInMessage = `
🔳 Enter message in format:

Username
Password

🔳 Example:

ivan password
`

const SignInErrMessage = `
⛔ Unable to continue further execution of sign in due to
`

const SignInSuccessMessage = `
✅ User successfully logged in!
`

const LogoutErrMessage = `
⛔ Unable to continue further execution of logout due to
`

const LogoutSuccessMessage = `
✅ User successfully logged out!
`

const SendOrderMessage = `
🔳 Enter message in format:

Symbol (one of symbols on kraken futures)
Side   (buy or sell)      
Size   (integer up to 25000)      

🔳 Example:

PI_XBTUSD buy 10000
`

const SendOrderErrMessage = `
⛔ Unable to continue further execution of send order due to
`

const SendOrderSuccessMessage = `
✅ Successfully send order!
`

const StartTradingMessage = `
🔳 Enter message in format:

Symbol (one of symbols on kraken futures)
Side   (buy or sell)      
Size   (integer up to 25000)
Take profit border (the value of the delta above which the order will be closed 📈)
Stop loss border (the value of the delta below which the order will be closed 📉)

🔳 Example:

PI_XBTUSD buy 10000 1000 1000
`

const StartTradingErrMessage = `
⛔ Unable to continue further execution of start trading due to
`

const StartTradingWillNotifyMessage = `
⌛ Bot will notify you when trading will stop
`

const StartTradingSuccessMessage = `
✅ Order have been successfully traded!
`

const GetUserOrdersErrMessage = `
⛔ Unable to continue further execution of get user orders due to
`
