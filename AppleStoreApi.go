package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type AppConnectUsers = struct {
	Data  []Datum              `json:"data"`
	Links AppConnectUsersLinks `json:"links"`
}

type Datum = struct {
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

type Attributes = struct {
	Username  string   `json:"username"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Roles     []string `json:"roles"`
}

type AppConnectUsersLinks = struct {
	Self string `json:"self"`
	Next string `json:"next"`
}

func CheckUserList(config *ConfigSettings, Username string) {
	token, err := CreateAppleJWT(config)
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}

	var nextUrl string = "https://api.appstoreconnect.apple.com/v1/users?limit=100"

	var FoundMatch = false

	for {
		req, err := http.NewRequest("GET", nextUrl, nil)
		if err != nil {
			fmt.Print(err.Error())
			break
		}

		req.Header.Add("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error while reading the response bytes:", err)
			break
		}

		var appConnectUsers AppConnectUsers

		err = json.Unmarshal(body, &appConnectUsers)
		if err != nil {
			log.Println("Error while deserializing the response bytes:", err)
			os.Exit(3)
		}

		for _, s := range appConnectUsers.Data {
			FoundMatch = strings.EqualFold(s.Attributes.Username, Username)

			if FoundMatch {
				fmt.Printf("Found %s, %s %s, %s\n",
					s.Attributes.Username,
					s.Attributes.FirstName,
					s.Attributes.LastName,
					strings.Join(s.Attributes.Roles, ", "))
				break
			}
		}

		if FoundMatch {
			break
		}

		nextUrl = appConnectUsers.Links.Next

		if len(nextUrl) == 0 {
			break
		}
	}

	if !FoundMatch {
		fmt.Println("No match for " + Username)
	}

}
