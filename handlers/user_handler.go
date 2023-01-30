package handlers

import (
	"net/http"

	"remood/models"
	"remood/pkg/auth"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecievedUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(ctx *gin.Context) {
	var receivedUser RecievedUser
	if ctx.ShouldBindJSON(&receivedUser) != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Message: "Fail to read user information",
			Error:   true,
		})
		return
	}

	var user models.User
	hashedPassword, err := auth.HashPassword(receivedUser.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Message: "Fail to generate hash password",
			Error:   true,
		})
		return
	}

	err = user.Create(receivedUser.Username, receivedUser.Email, hashedPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Message: err.Error(),
			Error:   true,
		})
		return
	}

	// Hide password
	user.Password = "*"

	ctx.JSON(http.StatusCreated, models.Response{
		Message: "Create User Successfully",
		Error:   false,
		Data: gin.H{
			"user": user,
		},
	})
}

func Login(ctx *gin.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")

	var user models.User
	if err := auth.ValidateUsername(&user, username); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "User not found",
			"error":   true,
		})
		return
	}

	if err := auth.ValidatePassword(user, password); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Wrong password",
			"error":   true,
		})
		return
	}

	// expireTime := time.Now().Add(time.Hour * 24)

	claims := auth.Claims{
		ID: user.ID,
	}

	tokenString, err := auth.GenerateTokenString(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail to generate token string",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successfully",
		"error":   false,
		"data": gin.H{
			"token": tokenString,
			"user":  user,
		},
	})
}

func GoogleLogin(ctx *gin.Context) {
	code := ctx.Query("code")

	googleUser, err := auth.GetGoogleUserInfo(code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Fail to get user info",
			"error":   true,
		})
		return
	}

	var user models.User
	err = user.GetOne("google_id", googleUser.ID)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Fail to check if user is existing",
				"error":   true,
			})
			return
		}

		err = user.Create(googleUser.Name, googleUser.Email, "")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Fail to create user",
				"error":   true,
			})
		}

		user.GoogleID = googleUser.ID
		user.Picture = googleUser.Picture

		err = user.Update(user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Fail to create user",
				"error":   true,
			})
		}
	}

	claims := auth.Claims{
		ID: user.ID,
	}

	tokenString, err := auth.GenerateTokenString(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fail to generate token string",
			"error":   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login by google successfully!",
		"error":   false,
		"data": gin.H{
			"token": tokenString,
			"user":  user,
		},
	})
}

func GetUser(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Message: "Invalid token",
			Error:   true,
		})
		return
	}

	user := models.User{}
	err = user.GetOne("_id", claims.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, models.Response{
			Message: "Fail to get user",
			Error:   true,
		})
	}

	// Hide password
	user.Password = "*"

	ctx.JSON(http.StatusOK, models.Response{
		Message: "Get user successfully",
		Error:   false,
		Data: gin.H{
			"user": user,
		},
	})
}

func UpdateUser(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Message: "Invalid token",
			Error:   true,
		})
		return
	}

	var newUser models.User

	if err := ctx.BindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Message: "Fail to read new user",
			Error:   true,
		})
		return
	}

	if newUser.ID != claims.ID {
		ctx.JSON(http.StatusForbidden, models.Response{
			Message: "Can not update this user",
			Error:   true,
		})
		return
	}

	if err := newUser.Update(newUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Message: "Can not update the user",
			Error:   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Message: "Update user successfully",
		Error:   false,
	})
}

func UpdatePassword(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Message: "Invalid token",
			Error:   true,
		})
		return
	}

	var user models.User

	if err := user.GetOne("_id", claims.ID); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Message: "User not found",
			Error:   true,
		})
		return
	}

	oldPassword := ctx.Query("old-password")
	newPassword := ctx.Query("new-password")

	err = auth.ValidatePassword(user, oldPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Message: "Wrong Password",
			Error:   true,
		})
		return
	}

	newHashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Message: "Fail to generate hashed password",
			Error:   true,
		})
		return
	}

	err = user.UpdatePassword(string(newHashedPassword))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Message: "Fail to update password",
			Error:   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Message: "Update Password Successfully",
		Error:   true,
	})
}

func DeleteUser(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Message: "Invalid token",
			Error:   true,
		})
		return
	}

	var user models.User
	user.ID = claims.ID

	if err := user.Delete(); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Message: "Can not delete the user",
			Error:   true,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Message: "Delete user successfully",
		Error:   false,
	})
}