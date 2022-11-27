package common

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	Id        string `json:"id"`
	UserName  string `json:"userName"`
	UserGroup string `json:"userGroup"`
	jwt.RegisteredClaims
}

func CreateToken(data Token) (string, error) {
	var err error
	exp := GetTokenExp()
	x := Token{
		data.Id,
		data.UserName,
		data.UserGroup,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(exp))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	atClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, x)
	token, err := atClaims.SignedString([]byte(GetTokenSecret()))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseToken(payload, secret string) (*Token, error) {
	var claims *Token
	token, err := jwt.ParseWithClaims(payload, &Token{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		} else if method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}

		// Use Default if not defined
		if secret != "" {
			return []byte(secret), nil
		} else {
			return []byte(GetTokenSecret()), nil
		}
	})
	if err != nil {
		log.Print(err)
		return claims, fmt.Errorf("failed parsing token")
	}

	if claims, ok := token.Claims.(*Token); ok && token.Valid {
		return claims, nil
	} else {
		log.Print("invalid token")
		return claims, fmt.Errorf("invalid token")
	}
}

func IsStrongPassword(pass string) bool {
	// Match 1 <= x Numbers
	match, err := regexp.MatchString(`[0-9]{1}`, pass)

	// Match 1 <= x Lowercase Characters
	if err == nil && match {
		match, err = regexp.MatchString(`[a-z]{1}`, pass)
	}

	// Match 1 <= x Uppercase Characters
	if err == nil && match {
		match, err = regexp.MatchString(`[A-Z]{1}`, pass)
	}

	// Match 1 <= x Symbols
	if err == nil && match {
		match, err = regexp.MatchString(`[^a-zA-Z0-9]{1}`, pass)
	}

	// Match Length 8 <= x <= 30
	if err == nil && 8 <= len(pass) && len(pass) <= 30 && match {
		match = true
	} else {
		match = false
	}

	return match
}

func HashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Print(err)
	} // GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Print(err)
		return false
	}

	return true
}
