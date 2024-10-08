package auth

import (
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHashPassword(t *testing.T) {
	password := "banana13"
	hash, err := HashPassword(password)
	if hash == password {
		t.Fatalf("HashPassword returned an unhashed password")
	}
	if err != nil {
		t.Fatalf("HashedPassword is throwing errors: %v", err)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "helloworld"
	hash, _ := HashPassword(password)
	err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash failed to match passwords: %v", err)
	}
}
