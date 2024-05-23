package bcrypt

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// Generate a salted hash for the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword compares a password with its hashed version
func ComparePassword(password, hashedPassword string) error {
	// Compare the hashed password with the input password
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
