package controllers

import (
	"go-crud-app/initializers"
	"go-crud-app/models"

	"github.com/gin-gonic/gin"
)

func PostsCreate(c *gin.Context) {
	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)

	post := models.Post{Title: body.Title, Body: body.Body}
	result := initializers.DB.Create(&post)

	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(201, gin.H{
		"data": post,
	})
}

func GetPosts(c *gin.Context) {
	var posts []models.Post
	initializers.DB.Find(&posts)

	c.JSON(200, gin.H{
		"data": posts,
	})
}

func GetPostById(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	initializers.DB.First(&post, id)

	c.JSON(200, gin.H{
		"data": post,
	})
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)
	var post models.Post
	initializers.DB.First(&post, id)

	initializers.DB.Model(&post).Updates(models.Post{
		Title: body.Title, 
		Body: body.Body,
	})

	c.JSON(200, gin.H{
		"data": post,
	})
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	initializers.DB.First(&post, id)

	initializers.DB.Model(&post).Delete(&models.Post{}, id)
	c.JSON(200, gin.H{
		"data": post,
	})
}
