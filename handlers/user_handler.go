package handlers

import (
	"net/http"
	"log"

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
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Fail to read user information"))
		return
	}

	var user models.User
	hashedPassword, err := auth.HashPassword(receivedUser.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to generate hash password"))
		return
	}

	err = user.Create(receivedUser.Username, receivedUser.Email, hashedPassword)
	if err != nil {
		log.Println(err.Error())
		var message string
		if (err.Error() == "username is existed" || err.Error() == "email is existed") {
			message = err.Error()
		} else {
			message = "Fail to create user"
		}
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse(message))
		return
	}

	// Hide password
	user.Password = "*"

	ctx.JSON(http.StatusOK, 	
		models.SuccessResponse("Create User Successfully", gin.H{"user": user}))
}

func Login(ctx *gin.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")

	var user models.User
	if err := user.GetOne("username", username); err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("User not found"))
		return
	}

	if err := auth.ValidatePassword(user.Password, password); err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Wrong password"))
		return
	}

	// expireTime := time.Now().Add(time.Hour * 24)

	claims := auth.Claims{
		ID: user.ID,
	}

	tokenString, err := auth.GenerateTokenString(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to generate token string"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Login successfully", gin.H{
			"token": tokenString,
			"user":  user,
		},))
}

func GoogleLogin(ctx *gin.Context) {
	code := ctx.Query("code")

	googleUser, err := auth.GetGoogleUserInfo(code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Fail to get user info"))
		return
	}

	var user models.User
	err = user.GetOne("google_id", googleUser.ID)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			ctx.JSON(http.StatusInternalServerError, 
				models.ErrorResponse("Fail to check if user is existing"))
			return
		}

		err = user.Create(googleUser.Name, googleUser.Email, "")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, 
				models.ErrorResponse("Fail to create user"))
		}

		user.GoogleID = googleUser.ID
		user.Picture = googleUser.Picture

		err = user.Update(user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, 
				models.ErrorResponse("Fail to create user"))
			return
		}
	}

	claims := auth.Claims{
		ID: user.ID,
	}

	tokenString, err := auth.GenerateTokenString(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to generate token string"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Login by google successfully!", gin.H{
			"token": tokenString,
			"user":  user,
		},))
}

func GetUser(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	user := models.User{}
	err = user.GetOne("_id", claims.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to get user"))
		return
	}

	// Hide password
	user.Password = "*"

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Get user successfully", gin.H{"user": user}))
}

func UpdateUser(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	var newUser models.User

	if err := ctx.BindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Fail to read new user"))
		return
	}

	if newUser.ID != claims.ID {
		ctx.JSON(http.StatusForbidden, 
			models.ErrorResponse("Can not update this user"))
		return
	}

	if err := newUser.Update(newUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Can not update the user"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Update user successfully", nil))
}

func UpdatePassword(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	var user models.User

	if err := user.GetOne("_id", claims.ID); err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("User not found"))
		return
	}

	oldPassword := ctx.Query("old-password")
	newPassword := ctx.Query("new-password")

	err = auth.ValidatePassword(user.Password, oldPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Wrong Password"))
		return
	}

	newHashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to generate hashed password"))
		return
	}

	err = user.UpdatePassword(string(newHashedPassword))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to update password"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Update Password Successfully", nil))
}

func ResetPassword(ctx *gin.Context) {
	email := ctx.Query("email")
	
	var user models.User 
	if err := user.ResetPassword(email); err != nil {
		ctx.JSON(http.StatusInternalServerError, 
			models.ErrorResponse("Fail to reset password"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Reset password successfully", gin.H{}))
}

func DeleteUser(ctx *gin.Context) {
	token := auth.GetTokenString(ctx)
	claims, err := auth.ParseToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, 
			models.ErrorResponse("Invalid token"))
		return
	}

	var user models.User
	user.ID = claims.ID

	if err := user.Delete(); err != nil {
		ctx.JSON(http.StatusForbidden, 
			models.ErrorResponse("Can not delete the user"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Delete user successfully", nil))
}