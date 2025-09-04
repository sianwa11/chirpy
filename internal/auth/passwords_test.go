package auth

import "testing"

func TestHashPassword(t *testing.T) {
	passwords := []string{"password123", "helloworld", "intID=11;"}
	hashedPasswords := []string{}
	for _, pass := range passwords {
		hashedPass, err := HashPassword(pass)
		if err != nil {
			t.Errorf("failed to hash password %v: %v\n", hashedPass, err)
		}
		hashedPasswords = append(hashedPasswords, hashedPass)
	}

	for i, hashedPass := range hashedPasswords {
		err := CheckPasswordHash(passwords[i], hashedPass)
		if err != nil {
			t.Errorf("password does not match hash: %v", err)
		}
	}

}