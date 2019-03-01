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

func LoadContacts(path string) (Contacts, error) {
	var contacts = Contacts{}

	//File not exist in path, implement crawl proxy
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return contacts, err
	}
	contact_file, err := os.Open(path)
	if err != nil {
		return contacts, err
	}
	defer contact_file.Close()

	err = json.NewDecoder(contact_file).Decode(&contacts)
	return contacts, err
}

func SaveContacts(path string, contacts Contacts) error {
	contactsJson, err := json.Marshal(contacts)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, contactsJson, 0644)
}

func CheckContactExist(contacts Contacts, id string) bool {
	for _, contact := range contacts.Contacts {
		if contact.Id == id {
			return true
		}
	}
	return false
}
