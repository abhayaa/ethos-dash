package keygen

import (
	"crypto/sha256"
	"encoding/base32"
	b64 "encoding/base64"
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
		arr[i] = rand.Intn(65)
	}

	firstHash := hashedString[arr[0] : arr[2]+6]
	secondHash := hashedString[arr[3] : arr[3]+4]
	thirdHash := hashedString[arr[1] : arr[1]+6]
	fourthHash := hashedString[arr[2] : arr[2]+4]

	finalHash := thirdHash + "-" + secondHash + "-" + firstHash + "-" + fourthHash

	return finalHash

}
