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

// struct generated via https://mholt.github.io/json-to-go/
type AppConnectErrors struct {
	Errors []struct {
		Status string `json:"status"`
		Code   string `json:"code"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
	} `json:"errors"`
}

type User struct {
	UserName  string
	FirstName string
	LastName  string
	Roles     string
}

func CheckUsers(config *ConfigSettings, users []string) {
	var FoundMatch = false
	var Status = ""
	teamMembers, err := GetUserList(config)

	if err == nil {
		// Walk through the users list
		FoundMatch = false
		var lastUser = ""
		for _, user := range users {
			if lastUser != user {
				lastUser = user
				var foundMember User
				Status = "NOMATCH"
				for _, teamMember := range teamMembers {
					FoundMatch = strings.EqualFold(teamMember.UserName, user)

					if FoundMatch {
						Status = "MATCH"
						foundMember = teamMember
						fmt.Printf("%s, %s %s %s, %s\n",
							Status,
							foundMember.UserName,
							foundMember.FirstName,
							foundMember.LastName,
							foundMember.Roles)
					} else {
						foundMember.UserName = user
					}
				}
			}
		}
	}
}

func GetUserList(config *ConfigSettings) ([]User, error) {

	var users []User
	var err error

	token, err := CreateAppleJWT(config)
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}

	var nextUrl string = "https://api.appstoreconnect.apple.com/v1/users?limit=100"

	for {
		req, err := http.NewRequest("GET", nextUrl, nil)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(3)
		}

		req.Header.Add("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
			os.Exit(3)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error while reading the response bytes:", err)
			os.Exit(3)
		}

		// Check for the AppConnect API returning an error first
		var appConnectErrors AppConnectErrors
		err = json.Unmarshal(body, &appConnectErrors)
		if err != nil {
			log.Println("Error while deserializing the response bytes:", err)
			os.Exit(3)
		}

		// If there is an error object in the body, print it and exit
		if len(appConnectErrors.Errors) > 0 {
			firstError := appConnectErrors.Errors[0]
			log.Println("Status:", firstError.Status)
			log.Println("Error accessing API:", firstError.Detail)
			os.Exit(4)
		}

		// Otherwise keep going
		var appConnectUsers AppConnectUsers

		err = json.Unmarshal(body, &appConnectUsers)
		if err != nil {
			log.Println("Error while deserializing the response bytes:", err)
			os.Exit(3)
		}

		for _, s := range appConnectUsers.Data {

			var user User
			user.UserName = s.Attributes.Username
			user.FirstName = s.Attributes.FirstName
			user.LastName = s.Attributes.LastName
			user.Roles = strings.Join(s.Attributes.Roles, ", ")
			users = append(users, user)
		}

		nextUrl = appConnectUsers.Links.Next

		if len(nextUrl) == 0 {
			break
		}
	}

	return users, err
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
			os.Exit(3)
		}

		req.Header.Add("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
			os.Exit(3)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error while reading the response bytes:", err)
			os.Exit(3)
		}

		// Check for the AppConnect API returning an error first
		var appConnectErrors AppConnectErrors
		err = json.Unmarshal(body, &appConnectErrors)
		if err != nil {
			log.Println("Error while deserializing the response bytes:", err)
			os.Exit(3)
		}

		// If there is an error object in the body, print it and exit
		if len(appConnectErrors.Errors) > 0 {
			firstError := appConnectErrors.Errors[0]
			log.Println("Status:", firstError.Status)
			log.Println("Error accessing API:", firstError.Detail)
			os.Exit(4)
		}

		// Otherwise keep going
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
