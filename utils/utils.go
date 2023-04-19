package utils

import "math/rand"

func GenerateRandomString(n int) string {
	output := make([]rune, n)
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	for i := range output {
		output[i] = chars[rand.Intn(len(chars))]
	}

	return string(output)
}
