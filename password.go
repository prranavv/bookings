package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "vk"
	hashedpass, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	fmt.Println(string(hashedpass))
}
