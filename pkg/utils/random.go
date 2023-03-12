package utils

import (
	"log"
	"strings"
	"math/rand"
)

func RandomPassword(length int) string {
	template := "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	var builder strings.Builder
	for i := 0; i < length; i++ {
		char := template[rand.Intn(len(template))]
		builder.WriteByte(char)
	}
	log.Println(builder.String())
	return builder.String()
}