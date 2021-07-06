package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Define command line options
	configPtr := flag.String("config", "./IsUserInApple.json", "Configuration file")
	usernamePtr := flag.String("username", "", "Username to find (in quotes)")

	flag.Parse()

	var userName string = *usernamePtr
	var ConfigFileName string = *configPtr

	if len(userName) == 0 {
		fmt.Println("Please specify an email address to match (in quotes)")
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

	fmt.Println("Looking for " + userName)

	config, err := ReadConfig(ConfigFileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	CheckUserList(config, userName)
}
