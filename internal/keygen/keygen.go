package keygen

import (
	"crypto/sha256"
	"encoding/base32"
	b64 "encoding/base64"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func Keygen(username string) string {
	stringGen := ""

	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

	stringGen = username + timestamp
	hash1 := firstHash([]byte(stringGen))
	hash2 := secondHash(hash1)
	finalHash := thirdHash(string(hash2))
	log.Print(finalHash)

	return finalHash

}

func firstHash(data []byte) []byte {
	hash := sha256.Sum256(data)

	return hash[:]
}

func secondHash(data []byte) []byte {

	hash := strings.ReplaceAll(
		b64.StdEncoding.EncodeToString(([]byte(data))), "+", "-")

	return []byte(hash)

}

func thirdHash(data string) string {

	hashedString := base32.StdEncoding.EncodeToString([]byte(data))

	arr := make([]int, 5)

	for i := range arr {
		arr[i] = rand.Intn(60)
	}
	log.Print("arr: ")
	log.Print(arr)

	firstHash := hashedString[arr[0] : arr[0]+5]
	log.Print("first " + firstHash)
	secondHash := hashedString[arr[3] : arr[3]+5]
	log.Print("second " + secondHash)
	thirdHash := hashedString[arr[1] : arr[1]+5]
	log.Print("third " + thirdHash)
	fourthHash := hashedString[arr[2] : arr[2]+5]
	log.Print("fourth " + fourthHash)

	finalHash := thirdHash + "-" + secondHash + "-" + firstHash + "-" + fourthHash

	return finalHash

}
