package main

import (
	"fmt"

	"./skype"
	"./utils"
)

func main() {
	SSkype := &skype.SkypeService{}
	SSkype.InitClient(utils.MIRCOSOFT_APP_ID, utils.MIRCOSOFT_APP_PASSWORD)
	er := SSkype.Authenticate(utils.AUTHENTICATE_URL)
	if er != nil {
		fmt.Println("Error: Can not login to bot service", er.Error())
		panic(er.Error())
	}
	fmt.Println(SSkype.Authorization)
}
