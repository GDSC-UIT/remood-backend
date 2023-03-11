package handlers 

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"log"
	"net/url"

	"remood/models"

	"github.com/gin-gonic/gin"
)

func CreateManyArticle(ctx *gin.Context) {
	var urls []string
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Fail read to articles info 1"))
		return	
	}

	err = json.Unmarshal(body, &urls)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, 
			models.ErrorResponse("Fail to read articles info 2"))
		return
	}


	var a models.Article
	articles, _ := a.CreateMany(urls)

	ctx.JSON(http.StatusAccepted, 
		models.SuccessResponse("Created articles successfully", gin.H{
			"articles": articles,
		}))
}

func GetAllArticles(ctx *gin.Context) {
	var a models.Article
	articles, err := a.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("Fail to get all articles"))
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("Get all articles succesfully", gin.H{
		"articles": articles,
	}))

}

func GetRandomArticles(ctx *gin.Context) {
	param := ctx.Query("number")
	number, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid parameter"))
		return
	}

	var a models.Article 
	articles, err := a.GetRandom(number)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("Fail to get random articles"))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("Get random articles succesfully", gin.H{
		"articles": articles,
	}))
}

func GetAllArticlesByTopic(ctx *gin.Context) {
	params := ctx.Query("topics")
	params, err := url.PathUnescape(params)
	log.Println(params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid parameter"))
		return
	}


	topics := strings.Split(params, ",")
	var a models.Article 
	articles, err := a.GetAllByTopic(topics)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("Fail to get all article by topics"))
		return 
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("Get all articles by topics successfully", gin.H{
		"articles": articles,
	}))
}

func GetRandomArticlesByTopic(ctx *gin.Context) {
	numberParam := ctx.Query("number")
	topicsParam := ctx.Query("topics")
	number, err := strconv.Atoi(numberParam)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid parameter"))
		return
	}

	topics := strings.Split(topicsParam, ",")

	var a models.Article 
	articles, err := a.GetRandomByTopics(topics, number)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("Fail to get random articles"))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("Get random articles from topics succesfully", gin.H{
		"articles": articles,
	}))
}

func DeleteArticles(ctx *gin.Context) {
	params := ctx.Query("ids")
	ids := strings.Split(params, ",")
	
	var a models.Article
	if err := a.DeleteMany(ids); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("Fail to delete articles"))
		return 
	}
	
	ctx.JSON(http.StatusOK, models.SuccessResponse("Delete articles successfully", gin.H{}))
}


func GetAllTopic(ctx *gin.Context) {
	var a models.Article 
	topics, err := a.GetAllTopics()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("Fail to get all topics"))
		return
	}

	ctx.JSON(http.StatusOK, 
		models.SuccessResponse("Get all topics successfully", gin.H{
			"topics": topics,
		}))
}