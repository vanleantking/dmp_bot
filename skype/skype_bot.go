package skype

import "net/http"

type Bot struct {
	BotEndpoint  *Endpoint
	ConnectorAPI *SkypeService
	AppID        string
	AppPassword  string
}

func InitBot(address string, appId string, appPassword string) *Bot {
	sSkype := SkypeService{AppID: appId, AppPassword: appPassword}
	sSkype.InitClient()

	endPoint := NewEndpoint(address)
	return &Bot{AppID: appId,
		AppPassword:  appPassword,
		ConnectorAPI: &sSkype,
		BotEndpoint:  endPoint}
}

func (bot *Bot) SetupEndpointHandler(handleBotActivity func(activity *Activity)) *http.Server {

	return bot.BotEndpoint.SetupServer(*bot.BotEndpoint.NewEndpointHandler(
		handleBotActivity,
		bot.ConnectorAPI.Authorization,
		bot.AppID))
}
