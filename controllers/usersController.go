package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/NurymGM/jwt-token/initializers"
	"github.com/NurymGM/jwt-token/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	// get email/pass from req body
	body := models.User{}
	if c.Bind(&body) != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to read body"})
		return
	}

	// hash the password
	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
		return
	}

	// create the user
	user := models.User{Email: body.Email, Password: string(hashed)}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	// respond
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "User created!"})
}

func LogIn(c *gin.Context) {
	// get email/pass from req body
	body := models.User{}
	if c.Bind(&body) != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to read body"})
		return
	}

	// look up requested user
	user := models.User{}
	initializers.DB.First(&user, "email = ?", body.Email)
	if user.ID == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password"})
		return
	}

	// compare passwords
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password"})
		return
	}

	// generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))		// dont forget the []byte here
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create token"})
		return
	}

	// send it back as a cookie (or you can directly respond with it)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600 * 24 * 30, "", "", true, true)

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Logged in! (jwt token is at cookies)"})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.IndentedJSON(http.StatusOK, gin.H{"message (validate)": user})
}