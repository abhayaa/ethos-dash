package keygen

import (
	"crypto/sha256"
	"encoding/base32"
	b64 "encoding/base64"
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

	//hash := b64.StdEncoding.EncodeToString(([]byte(NewSHA256(data))))
	hash := strings.ReplaceAll(
		b64.StdEncoding.EncodeToString(([]byte(data))), "+", "-")

	return []byte(hash)

}

func thirdHash(data string) string {

	hashedString := base32.StdEncoding.EncodeToString([]byte(data))

	firstHash := hashedString[5:11]
	secondHash := hashedString[21:27]
	thirdHash := hashedString[49:55]
	fourthHash := hashedString[64:70]

	finalHash := thirdHash + "-" + secondHash + "-" + firstHash + "-" + fourthHash

	return finalHash

}
