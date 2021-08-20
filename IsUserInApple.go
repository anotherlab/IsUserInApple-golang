package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	// Define command line options
	configPtr := flag.String("config", "./IsUserInApple.json", "Configuration file")
	usernamePtr := flag.String("username", "", "Username to find (in quotes)")
	userFileListPtr := flag.String("userlist", "", "list of usernames")

	flag.Parse()

	var userName string = *usernamePtr
	var userFileList string = *userFileListPtr
	var ConfigFileName string = *configPtr

	if (len(userName) == 0) && (len(userFileList) == 0) {
		fmt.Println("Please specify an username or a list of users to match")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// If a config file was not specified from the commmand line, assume it's
	// in the same folder as the executable
	if len(ConfigFileName) == 0 {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		ConfigFileName = filepath.Join(filepath.Dir(ex), "IsUserInApple.json")
	}

	if _, err := os.Stat(ConfigFileName); os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(2)
	}

	var users []string
	var err error

	if len(userFileList) > 0 {
		if _, err := os.Stat(userFileList); os.IsNotExist(err) {
			fmt.Println(err)
			os.Exit(2)
		}
		users, err = readLines(userFileList)

		if err != nil {
			log.Fatalf("readLines: %s", err)
		}

		// Sort the list so that if there are duplicates, we can skip over the dups
		sort.Strings(users)
	} else {
		// Treat a single user name as an array of 1 element
		fmt.Println("Looking for " + userName)
		users = append(users, userName)
	}

	config, err := ReadConfig(ConfigFileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	CheckUsers(config, users)
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	// Close the file when the readLines function returns
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, strings.ToLower(scanner.Text()))
	}

	return lines, scanner.Err()
}
