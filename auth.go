package transferrer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	emailContextKey = "email"
	secret          = "verseAPI"
)

// Error encapsulates the type and the message of an error
type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Error prints the error in the desired format
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Login defines the required data to acces the app
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Claim is the type for add information inside the claims of the JWT
type Claim struct {
	*jwt.StandardClaims
	Login
}

// Token contains the three parts of the token
type Token struct {
	Type      string `json:"type"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

// GenerateJWT returns a signed token with the email and password provided
func GenerateJWT(c *gin.Context) {
	var login Login
	err := json.NewDecoder(c.Request.Body).Decode(&login)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	expiresAt := time.Now().Add(time.Minute * 60).Unix()

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = &Claim{
		&jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
		Login{login.Email,
			login.Password,
		},
	}

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, Token{
		Token:     tokenStr,
		Type:      "Bearer",
		ExpiresAt: expiresAt,
	})
}

// ValidateToken checks the JWT in the Authorization header and validate it. If it's correct, adds the info to the context.
func ValidateToken(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if len(authHeader) == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte(secret), nil
	})

	mssg := "Invalid token"
	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			validationError := err.(*jwt.ValidationError)
			switch validationError.Errors {
			case jwt.ValidationErrorExpired:
				mssg = "Token expired"
			default:
				mssg = "Invalid token"
			}
		}
	}
	err = &Error{
		Type:    "JWTError",
		Message: mssg,
	}
	if token == nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	if token.Valid {
		c.Set(emailContextKey, token.Claims.(jwt.MapClaims)["email"])
		c.Next()
		return
	}

	c.JSON(http.StatusForbidden, err)
	return
}
