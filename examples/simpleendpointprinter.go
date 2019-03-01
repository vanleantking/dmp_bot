/*
MIT License

Copyright (c) 2017 MichiVIP

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package main

import (
	"encoding/json"
	"fmt"

	"../skype"
	"../utils"
)

var authorizationBearerToken string

const (
	actionHookPath string = "/skype/actionhook"
)

/*
In this example I am setting up a basic skype APi endpoint and print activity objects.
We will just use the given http.Server to listen to incoming requests.
*/

func startSimpleEndpointPrinter() {
	// bad practice. In real production you should better request the token via skypeapi.RequestAccessToken
	// WARNING: when using a static authorization token it could expire. In future the will be an automatic refresher
	// authorizationBearerToken := "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ii1zeE1KTUxDSURXTVRQdlp5SjZ0eC1DRHh3MCIsImtpZCI6Ii1zeE1KTUxDSURXTVRQdlp5SjZ0eC1DRHh3MCJ9.eyJhdWQiOiJodHRwczovL2FwaS5ib3RmcmFtZXdvcmsuY29tIiwiaXNzIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvZDZkNDk0MjAtZjM5Yi00ZGY3LWExZGMtZDU5YTkzNTg3MWRiLyIsImlhdCI6MTU1MTMzNDgyOSwibmJmIjoxNTUxMzM0ODI5LCJleHAiOjE1NTEzMzg3MjksImFpbyI6IjQySmdZQ2c5T1VHMnl0WGhnT1QwVFBkYjZ2cHpBUT09IiwiYXBwaWQiOiI2MzJjZmRkOC1lMmRiLTQzZWYtOTY1My04MjA1NzI5YTEwZjkiLCJhcHBpZGFjciI6IjEiLCJpZHAiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC9kNmQ0OTQyMC1mMzliLTRkZjctYTFkYy1kNTlhOTM1ODcxZGIvIiwidGlkIjoiZDZkNDk0MjAtZjM5Yi00ZGY3LWExZGMtZDU5YTkzNTg3MWRiIiwidXRpIjoiQUI4dHU0UG5KVWl5VjN6Z0psTVpBQSIsInZlciI6IjEuMCJ9.K26B4JBQAaF0m42kYR_9J_4B1KmJT09b8D_aGq7Wzen3-HrA6W1OOWAIqvoT6MuDr9CCA0qF5Z-XMjCWZuQkhUXNt7iYgsThWDK1W8ceBhrgRIO7sHAN53fV7f1coU2TaaIjzolZnX364bTJm_nQ0ULylW9qQvzf9lKMrq_4lUMnmigydehd2NGpzJTHyHmzG_96KaKjcpufz-InXoauMKucEBhoFiwkcy7Eqy0XTESOvMx_UlVBKX-_Van_Wgpvq5YeCiwQUcTrBOGb_f_GLJaJo75cZi97Rv5gOo-g1X1iqm1snEB8r-rkUm8IJtF1ekYh9O7Yc2wDNOSbAWORFw"
	fmt.Println("start custom server endpoint", actionHookPath)
	authorizationBearer, _ := skype.RequestAccessToken(utils.MIRCOSOFT_APP_ID, utils.MIRCOSOFT_APP_PASSWORD)
	authorizationBearerToken = authorizationBearer.TokenType + " " + authorizationBearer.AccessToken
	// Endpoint is going to listen on 0.0.0.0:8080
	endpoint := skype.NewEndpoint(":9443")

	// we define our own handle function
	srv := endpoint.SetupServer(*skype.NewEndpointHandler(func(activity *skype.Activity) {
		bytes, _ := json.MarshalIndent(activity, "", "  ")
		fmt.Println(string(bytes))
	}, authorizationBearerToken, utils.MIRCOSOFT_APP_ID))
	// finally we just use the default method to start the server
	panic(srv.ListenAndServe())
}

func main() {
	startSimpleEndpointPrinter()
}
