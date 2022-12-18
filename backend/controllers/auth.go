package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/ip-05/quizzus/config"
	"github.com/ip-05/quizzus/middleware"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type UserInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
}

type AuthController struct{}

var cfg = config.GetConfig()

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  cfg.Frontend.Base,
	ClientID:     cfg.Google.ClientId,
	ClientSecret: cfg.Google.ClientSecret,
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

func (a AuthController) GoogleLogin(c *gin.Context) {
	var expiration = int(20 * time.Minute)
	b := make([]byte, 16)
	rand.Read(b)
	oauthState := base64.URLEncoding.EncodeToString(b)
	c.SetCookie("oauthstate", oauthState, expiration, "*", cfg.Server.Domain, cfg.Server.Secure, false)

	url := googleOauthConfig.AuthCodeURL(oauthState)

	c.JSON(http.StatusOK, gin.H{"redirectUrl": url})
}

func (a AuthController) GoogleCallback(c *gin.Context) {
	oauthState, err := c.Cookie("oauthstate")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cookie"})
		return
	}

	if c.Request.FormValue("state") != oauthState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while verifying auth token"})
		return
	}

	token, err := googleOauthConfig.Exchange(context.Background(), c.Request.FormValue("code"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while verifying auth token"})
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while verifying auth token"})
		return
	}

	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while verifying auth token"})
		return
	}

	var userInfo UserInfo

	err = json.Unmarshal(contents, &userInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while parsing user info"})
		return
	}

	secretKey := []byte(cfg.Secrets.Jwt)
	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":             userInfo.Id,
		"name":           userInfo.GivenName,
		"email":          userInfo.Email,
		"profilePicture": userInfo.Picture,
		"exp":            time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	tokenString, err := tokenJWT.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while signing JWT token"})
		return
	}

	c.SetCookie("token", tokenString, 7*24*60*60, "/", cfg.Server.Domain, cfg.Server.Secure, false)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user", "token": tokenString})
}

func (a AuthController) Me(c *gin.Context) {
	authedUser, _ := c.Get("authedUser")
	user := authedUser.(middleware.AuthedUser)
	c.JSON(http.StatusOK, user)
}
