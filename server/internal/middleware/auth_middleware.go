package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Pranavp37/magic_movie_stream/internal/config"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	User_id string `json:"user_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

type JwtCustomClaimsRespose struct {
	User_id string `json:"user_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

type JwtTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var jwtSecretKey = config.LoadConfig().JWT_SECRET_KEY
var tokenKey_byte = []byte(config.LoadConfig().JWT_SECRET_KEY)

// GenerateToken generates a JWT token
func GenerateToken(user_id, email, name, role string, expireAt time.Time, secret_key []byte) (string, error) {

	claims := &JwtCustomClaims{
		User_id: user_id,
		Name:    name,
		Email:   email,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret_key)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateAccessAndRefreshToken generates both access and refresh tokens
func GenerateAccessAndRefreshToken(user_id, email, role, name string) (jwtTokens JwtTokens, err error) {

	expirationTime := time.Now().Add(15 * time.Minute)

	accessToken, err := GenerateToken(user_id, email, name, role, expirationTime, tokenKey_byte)

	if err != nil {
		return JwtTokens{}, err
	}

	expirationTime = time.Now().Add(7 * 24 * time.Hour)

	refreshToken, err := GenerateToken(user_id, email, name, role, expirationTime, tokenKey_byte)

	if err != nil {
		return JwtTokens{}, err
	}

	return JwtTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// JwtTokenFromHeader extracts the JWT token from the Authorization header
func JwtTokenFromHeader(c *gin.Context) string {

	vals := c.GetHeader("Authorization")

	// splitedtoken := strings.Split(vals, "Bearer")

	// if len(splitedtoken) != 2 {
	// 	return ""
	// }

	splitedtoken := strings.SplitN(vals, " ", 2)

	if len(splitedtoken) != 2 {
		return ""
	}

	return splitedtoken[1]

}

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := JwtTokenFromHeader(c)

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found in header"})
			c.Abort()
			return
		}

		keyFunc := func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("token is invalid")
			}
			return []byte(jwtSecretKey), nil
		}

		token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, keyFunc)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(*JwtCustomClaims)

		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("user", claims)

		c.Next()

	}
}

func GetTokenData(tokenString string) (*JwtCustomClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("token is invalid")
		}
		return []byte(jwtSecretKey), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, keyFunc)

	if err != nil {
		return nil, err
	}

	clien := token.Claims.(*JwtCustomClaims)
	return clien, nil

}
