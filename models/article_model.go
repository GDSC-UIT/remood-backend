package models

import (
	"context"
	"log"
	"sort"
	"time"

	"remood/pkg/const/collections"
	"remood/pkg/database"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	BaseModel `json:",inline" bson:",inline"`

	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Url    string             `json:"url"`
	Title  string             `json:"title,omitempty"`
	Topic  string             `json:"topic,omitempty"`
	Author string             `json:"author,omitempty"`
	Image  string             `json:"image,omitempty"`
}

// Extract article's information from the given url
func newArticleFromURL(url string) (Article, error) {
	var article Article
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return article, err
	}

	article.Url = url

	article.CreatedAt = time.Now().Unix()
	article.UpdatedAt = time.Now().Unix()
	article.ID = primitive.NewObjectID()

	// Extract the title
	title := doc.Find("title").Text()
	article.Title = title

	// Extract the author
	author, _ := doc.Find("meta[property='article:author']").Attr("content")
	article.Author = author

	// Extract the image
	image, _ := doc.Find("meta[property='og:image']").Attr("content")
	article.Image = image

	// Extract the title
	section, _ := doc.Find("meta[property='article:section']").Attr("content")
	article.Topic = section

	return article, nil
}

func (a *Article) CreateMany(urls []string) ([]Article, error) {
	articles := make([]Article, 0)
	for _, url := range urls {
		article, err := newArticleFromURL(url)
		if err == nil {
			articles = append(articles, article)
		}
	}

	insert := make([]interface{}, 0)
	for _, article := range articles {
		insert = append(insert, article)
	}

	collection := database.GetMongoInstance().Db.Collection(string(collections.Article))
	collection.InsertMany(context.Background(), insert)

	// Return successful insert even when some articles can not be insert
	return articles, nil
}

func (a *Article) GetAll() ([]Article, error) {
	var articles []Article
	collection := database.GetMongoInstance().Db.Collection(string(collections.Article))
	cursor, err := collection.Find(context.Background(), gin.H{})

	if err != nil {
		return articles, err
	}

	err = cursor.All(context.Background(), &articles)
	log.Println(err)
	return articles, err
}

func (a *Article) GetRandom(number int) ([]Article, error) {
	var articles []Article
	collection := database.GetMongoInstance().Db.Collection(string(collections.Article))

	pipeline := []interface{}{
		bson.M{
			"$sample": bson.M{
				"size": number,
			},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return articles, err
	}

	err = cursor.All(context.Background(), &articles)
	return articles, err
}

func (a *Article) GetAllByTopic(topics []string) ([]Article, error) {
	var result []Article
	collection := database.GetMongoInstance().Db.Collection(string(collections.Article))

	for _, topic := range topics {
		var articles []Article
		cursor, err := collection.Find(context.Background(), gin.H{"topic": topic})

		if err != nil {
			return result, err
		}

		err = cursor.All(context.Background(), &articles)
		if err != nil {
			return result, err
		}
		result = append(result, articles...)
	}
	return result, nil
}

func (a *Article) GetRandomByTopics(topics []string, number int) ([]Article, error) {
	var articles []Article
	collection := database.GetMongoInstance().Db.Collection(string(collections.Article))

	pipeline := []interface{}{
		bson.M{
			"$match": bson.M{
				"topic": bson.M{
					"$in": topics,
				},
			},
		},
		bson.M{
			"$sample": bson.M{
				"size": number,
			},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return articles, err
	}

	err = cursor.All(context.Background(), &articles)
	return articles, err
}

func (a *Article) DeleteMany(ids []string) error {
	log.Println(ids)
	collection := database.GetMongoInstance().Db.Collection(string(collections.Article))
	for _, rawId := range ids {
		id, err := primitive.ObjectIDFromHex(rawId)
		log.Println(id)
		if err == nil {
			collection.DeleteOne(context.Background(), bson.M{"_id": id})
		}
	}

	// return successful delete even when can not delete some articles
	return nil
}

func (a *Article) GetAllTopics() ([]string, error) {
	var topics []string
	collection := database.GetMongoInstance().Db.Collection(string(collections.Article))

	cursor, err := collection.Find(context.Background(), gin.H{})

	if err != nil {
		log.Println(err)
		return topics, err
	}

	var articles []Article

	err = cursor.All(context.Background(), &articles)
	if err != nil {
		return topics, nil
	}

	// Find all the separate topics by sorting the slice and
	// find the topics of the articles not having the same topic with the previous article
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Topic < articles[j].Topic
	})

	// Bug: the first topic is ""

	for i := range articles {
		if articles[i].Topic == ""{
			log.Println(articles[i])
		}
		if i == 0 {
			topics = append(topics, articles[i].Topic)
		} else {
			if articles[i].Topic != articles[i-1].Topic {
				topics = append(topics, articles[i].Topic)
			}
		}
	}

	return topics, nil
}