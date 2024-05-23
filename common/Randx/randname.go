package Randx

import (
	"math/rand"
	"time"
)


func RandName()string{
	rand.Seed(time.Now().UnixNano())

    // Set of allowed characters
    chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

    // Generate random text of length 10
    text := ""
    for i := 0; i < 10; i++ {
        text += string(chars[rand.Intn(len(chars))])
    }
	return text
}