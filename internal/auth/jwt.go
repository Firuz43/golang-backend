package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 1. The Claims: This is the "Data" inside your VIP pass.
// We use Capital letters (UserID, Email) so the JWT library can "export"
// them into the JSON string that lives inside the token.
type Claims struct {
	UserID               string `json:"user_id"`
	Email                string `json:"email"`
	jwt.RegisteredClaims        // This adds standard fields like 'exp' (expiration)
}

// 2. The Secret Key: This pulls the "Wax Seal" from your .env file.
// In Go, the JWT library needs the secret as a []byte (byte slice).
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Fallback for local development if you forget to set the .env
		return []byte("my_temporary_local_secret")
	}
	return []byte(secret)
}

// 3. The Token Factory: This function takes user info and spits out the signed JWT.//
func GenerateToken(userID string, email string) (string, error) {
	// Step A: Fill out the "ID Card" (Claims)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			// The pass expires 24 hours from now
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ecommerce-api", // Useful for identifying where the token came from
		},
	}

	// Step B: Choose the "Lock" (HS256 is the standard signing algorithm)//g
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Step C: Apply the "Wax Seal" (Sign the token with our Secret Key)
	tokenString, err := token.SignedString(GetJWTSecret())
	if err != nil {
		return "", err // If something goes wrong (like a bad key), return the error
	}

	return tokenString, nil
}

// ValidateToken parses the JWT string and returns the claims if valid
func ValidateToken(tokenString string) (*Claims, error) {
	// 1. Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the "Signing Method" is what we expect (HS256)
		return GetJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	// 2. Extract the data (Claims)
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
