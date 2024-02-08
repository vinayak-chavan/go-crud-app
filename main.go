package main

// CompileDaemon -command="./go-crud-app"

import (
	"context"
	"go-crud-app/controllers"
	"go-crud-app/initializers"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()

	r.POST("/post", controllers.PostsCreate)
	r.GET("/post", controllers.GetPosts)
	r.GET("/post/:id", controllers.GetPostById)
	r.PUT("/post/:id", controllers.UpdatePost)
	r.DELETE("/post/:id", controllers.DeletePost)

	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*")
	r.MaxMultipartMemory = 8 << 20

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.POST("/server", func(c *gin.Context) {
		file, err := c.FormFile("image")

		if err != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"error": "Failed to upload image",
			})
			return
		}

		f, openErr := file.Open()
		if openErr != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"error": "Failed to upload image",
			})
			return
		}

		result, uploadErr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String("my-bucket"),
			Key:    aws.String("my-object-key"),
			Body:   f,
			ACL: "public-read",
		})

		if uploadErr != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"error": "Failed to upload image",
			})
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"image": result.Location,
		})
	})

	r.POST("/local", func(c *gin.Context) {
		file, err := c.FormFile("image")

		if err != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"error": "Failed to upload image",
			})
			return
		}

		err = c.SaveUploadedFile(file, "assets/uploads/"+file.Filename)

		if err != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"error": "Failed to upload image",
			})
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"image": "assets/uploads/" + file.Filename,
		})
	})

	r.Run()
}
