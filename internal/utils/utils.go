package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type User struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
	Flags         int    `json:"flags"`
	Banner        string `json:"banner"`
	BannerColor   string `json:"banner_color"`
	AccentColor   string `json:"accent_color"`
	Locale        string `json:"locale"`
	MFA           bool   `json:"mfa_enabled"`
	Premium       int    `json:"premium_type"`
	Email         string `json:"email"`
	Verified      bool   `json:"verified"`
}

const (
	userEndpoint string = "https://discord.com/api/users/@me"
)

func GetEnvKey(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading .env file")
	}

	return os.Getenv(key)
}

func ValidateApiKey(key string) bool {
	if GetEnvKey("API_KEY") != key {

		log.Printf("failed api key match %s", key)
		return false
	}
	return true
}

func GetDiscordInfo(code string) *User {

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", userEndpoint, nil)
	if err != nil {
		log.Print(err)
	}

	log.Print(code)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+code)
	resp, err := client.Do(req)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		log.Print(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Non 200 response code received %d", resp.StatusCode)
	}

	bits, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}

	r := bytes.NewReader(bits)
	decoder := json.NewDecoder(r)

	val := &User{}
	err = decoder.Decode(val)
	if err != nil {
		log.Print(err)
	}

	return val
}
