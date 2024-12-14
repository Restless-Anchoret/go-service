package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type signInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type signInResponse struct {
	User userDto `json:"user"`
}

type signUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type signUpResponse struct {
	User userDto `json:"user"`
}

type userDto struct {
	ID    string `json:"id"`
	Login string `json:"login"`
}

type user struct {
	ID           string
	Login        string
	PasswordHash string
}

var users []user

func main() {
	router := gin.Default()
	router.POST("/api/v1/sign-in", signIn)
	router.POST("/api/v1/sign-up", signUp)

	router.Run("localhost:8080")
}

func signIn(c *gin.Context) {
	var request signInRequest

	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request"})
		return
	}

	passwordHash := getPasswordHash(request.Password)

	for _, u := range users {
		if u.Login == request.Login && u.PasswordHash == passwordHash {
			userDto := userDto{ID: u.ID, Login: u.Login}
			response := signInResponse{User: userDto}
			c.IndentedJSON(http.StatusOK, response)
			return
		}
	}

	c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid login or password"})
}

func signUp(c *gin.Context) {
	var request signUpRequest

	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request"})
		return
	}

	for _, u := range users {
		if u.Login == request.Login {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Login already exists"})
			return
		}
	}

	userId := fmt.Sprintf("%v", len(users)+1)
	passwordHash := getPasswordHash(request.Password)
	user := user{ID: userId, Login: request.Login, PasswordHash: passwordHash}
	users = append(users, user)

	userDto := userDto{ID: user.ID, Login: user.Login}
	response := signUpResponse{User: userDto}
	c.IndentedJSON(http.StatusOK, response)
}

func getPasswordHash(password string) string {
	passwordHashBytes := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%v", passwordHashBytes)
}
