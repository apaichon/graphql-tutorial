package auth

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	
)

func GenerateRandomSalt(length int) string {
	// Define the character set from which to generate the salt
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Seed the random number generator
	//rand.Seed(time.Now().UnixNano())

	// Create a byte slice to store the salt characters
	salt := make([]byte, length)

	// Fill the byte slice with random characters from the charset
	for i := 0; i < length; i++ {
		salt[i] = charset[rand.Intn(len(charset))]
	}

	// Convert the byte slice to a string and return it
	return string(salt)
}

func HashString(input string) string {

	// Calculate the MD5 hash of the concatenated string
	hasher := md5.New()
	hasher.Write([]byte(input))
	hash := hasher.Sum(nil)

	// Convert the hash to a hexadecimal string
	hashString := hex.EncodeToString(hash)

	return hashString
}
