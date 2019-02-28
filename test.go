package main

import (
	"fmt"

	"crypto/tls"
	"net/http"

	"./skype"
	"./utils"
)

const (
	actionHookPath     string = "/skype/actionhook"
	address                   = ":9443"
	someOtherStuffPath string = "/"
)
var (
	authorizationBearerToken string
)

func main() {
	// SSkype := &skype.SkypeService{}
	// SSkype.InitClient(utils.MIRCOSOFT_APP_ID, utils.MIRCOSOFT_APP_PASSWORD)
	// er := SSkype.Authenticate(utils.AUTHENTICATE_URL)
	// if er != nil {
	// 	fmt.Println("Error: Can not login to bot service", er.Error())
	// 	panic(er.Error())
	// }
	// fmt.Println(SSkype.Authorization)
	// conver := skype.Conversation{
	// 	Bot:     skype.ChannelAccount{ID: "ea228964-bb55-42c7-9bcc-ea90bbe53b92", Name: "ureka_dmp_bot"},
	// 	IsGroup: false,
	// 	Members: []skype.ChannelAccount{
	// 		skype.ChannelAccount{ID: "ea228964-bb55-42c7-9bcc-ea90bbe53b92", Name: "ureka_dmp_bot"}},
	// 	TopicName: "Start conversation"}
	// converresp, err := SSkype.StartConversation(utils.MESSAGE_TRANSFER_URL+utils.START_CONVERSATION, conver)
	// fmt.Println(converresp, err)

	// message := skype.Messages{Type: "non-reply", Text: "hello world"}
	// SSkype.SendMessage(message)
	startCustomServerEndpoint()
}

// this function handles our skype activity
func handleActivity(activity *skype.Activity) {
	fmt.Println("say hello", activity)
	if activity.Type == "message" {
		// hard coding an auth token is no good practice! I am just doing this to make this example more simple.
		if err := skype.SendReplyMessage(activity, "Good evening. Nice to meet you!", authorizationBearerToken); err != nil {
			panic(err.Error())
		} else {
			fmt.Println("Successfully sent response message to skype user: " + activity.From.Name)
		}
	}
}

// our custom application handler function
func handleMainPath(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("This content is hilarious."))
}

func startCustomServerEndpoint() {
	fmt.Println("start custom server endpoint", actionHookPath)
	authorizationBearer, _ := skype.RequestAccessToken(utils.MIRCOSOFT_APP_ID, utils.MIRCOSOFT_APP_PASSWORD)
	mux := http.NewServeMux()
	// here we setup an own activity handler which listens to the path "/skype/actionhook"
	authorizationBearerToken = authorizationBearer.TokenType + " " + authorizationBearer.AccessToken
	fmt.Println("Token, ", authorizationBearerToken)
	mux.Handle(actionHookPath, skype.NewEndpointHandler(handleActivity, authorizationBearerToken, utils.MIRCOSOFT_APP_ID))
	// here we could probably just handle our main application
	mux.HandleFunc(someOtherStuffPath, handleMainPath)
	// here you could provide your own TLS configuration
	customTlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	// custom server setup
	srv := &http.Server{
		Addr:         address,
		Handler:      mux,
		TLSConfig:    customTlsConfig,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	// finally we just use the default method to start the server
	panic(srv.ListenAndServe())
}
