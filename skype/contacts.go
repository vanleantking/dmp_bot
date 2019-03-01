package skype

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Contacts struct {
	Contacts []Contact `json:"contacts"`
}
type Contact struct {
	Id         string `json:"id"`
	IsGroup    bool   `json:"is_group"`
	ServiceURL string `json:"serviceUrl"`
	Name       string `json:"name"`
}

func LoadContacts(path string) ([]Contact, error) {
	var contactsType = Contacts{}
	var contacts = []Contact{}

	//File not exist in path, implement crawl proxy
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return contacts, err
	}
	contact_file, err := os.Open(path)
	if err != nil {
		return contacts, err
	}
	defer contact_file.Close()

	err = json.NewDecoder(contact_file).Decode(&contactsType)
	if err != nil {
		return contacts, err
	}
	for _, contact := range contactsType.Contacts {
		contacts = append(contacts, contact)
	}
	return contacts, err
}

func SaveContacts(path string, contacts []Contact) error {
	contactType := Contacts{Contacts: contacts}
	contactsJson, err := json.Marshal(contactType)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, contactsJson, 0644)
}

func CheckContactExist(contacts []Contact, id string) (bool, int) {
	for index, contact := range contacts {
		if contact.Id == id {
			return true, index
		}
	}
	return false, 0
}

func CheckInChannelList(contacts []ChannelAccount, id string) (bool, int) {
	for index, contact := range contacts {
		if contact.ID == id {
			return true, index
		}
	}
	return false, 0
}
