package auth

import (
	"fmt"
	"mist/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func LoginHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		type LoginRequest struct {
			Username string
			Password string
		}
		var req LoginRequest
		req.Username = c.PostForm("username")
		req.Password = c.PostForm("password")
		fmt.Println("Login attempt for user:", req.Username)
		userID, err := AuthenticateUser(req.Username, req.Password)
		if err != nil {
			c.String(200, "Invalid username or password")
			return
		}
		fmt.Println("User authenticated, generating token")
		jwtToken, err := utils.SignJWT(userID)
		if err != nil {
			c.String(200, "Failed to generate token")
			return
		}
		fmt.Println("Token generated:", jwtToken)
		err = RefreshSession(userID, jwtToken, utils.Pointer(time.Now().Add(7*24*time.Hour)))
		if err != nil {
			fmt.Println("Failed to refresh session:", err)
			c.String(200, "Failed to refresh session, try again later")
			return
		}
		fmt.Println("Session refreshed for user:", userID)
		c.SetCookie("mist-token", jwtToken, 3600*24*7, "/", "", false, true)
		c.Header("HX-Redirect", "/home")
		c.Status(200)
	}
}

func RegisterHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		type RegisterRequest struct {
			FullName        string
			Username        string
			Email           string
			Password        string
			ConfirmPassword string
		}
		var req RegisterRequest
		req.FullName = c.PostForm("fullName")
		req.Username = c.PostForm("username")
		req.Email = c.PostForm("email")
		req.Password = c.PostForm("password")
		req.ConfirmPassword = c.PostForm("passwordConfirm")
		if req.Password != req.ConfirmPassword {
			c.String(200, "Passwords do not match")
			return
		}
		isFirstUser, err := IsFirstUser()
		if err != nil {
			c.String(200, "Internal server error, try again later")
			return
		}
		if !isFirstUser {
			return
		}
		passwordHash, err := utils.HashPassword(req.Password)
		if err != nil {
			c.String(200, "Failed to hash password")
			return
		}
		req.Password = passwordHash
		userID, err := RegisterUser(req.FullName, req.Username, req.Email, req.Password)
		if err != nil {
			c.String(200, "Failed to register user, try again later")
			return
		}
		jwtToken, err := utils.SignJWT(userID)
		if err != nil {
			c.String(200, "Failed to generate token")
			return
		}
		err = CreateSession(userID, jwtToken, utils.Pointer(time.Now().Add(7*24*time.Hour)))
		if err != nil {
			c.String(200, "Failed to create session, try again later")
			return
		}
		c.SetCookie("mist-token", jwtToken, 3600*24*7, "/", "", false, true)
		c.Header("HX-Redirect", "/home")
		c.Status(200)
	}
}
