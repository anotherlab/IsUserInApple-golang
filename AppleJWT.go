package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

type ConfigSettings struct {
	PrivateKeyFile string `json:"PrivateKeyFile"`
	KeyID          string `json:"KeyID"`
	IssuerID       string `json:"IssuerID"`
}

func ReadConfig(ConfigFileName string) (*ConfigSettings, error) {
	file, err := ioutil.ReadFile(ConfigFileName)

	if err != nil {
		return nil, err
	}

	config := new(ConfigSettings)

	err = json.Unmarshal([]byte(file), &config)

	return config, err
}

func CreateAppleJWT(settings *ConfigSettings) (string, error) {
	bytes, err := ioutil.ReadFile(settings.PrivateKeyFile)

	if err != nil {
		fmt.Println(err)
	}

	x509Encoded, _ := pem.Decode(bytes)

	parsedKey, err := x509.ParsePKCS8PrivateKey(x509Encoded.Bytes)

	if err != nil {
		log.Fatal(err)
	}

	ecdsaPrivateKey, ok := parsedKey.(*ecdsa.PrivateKey)

	if !ok {
		panic("not ecdsa private key")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": settings.IssuerID,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
		"aud": "appstoreconnect-v1",
	})

	token.Header["kid"] = settings.KeyID

	tokenString, err := token.SignedString(ecdsaPrivateKey)

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}
