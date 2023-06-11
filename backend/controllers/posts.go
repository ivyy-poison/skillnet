package controllers

import (
	"errors"
	"io"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
)

const postNotFoundMessage = "Post not found"

type APIEnv struct {
	DB          *gorm.DB
	GoogleCloud *storage.Client
}

func (a *APIEnv) GetPosts(context *gin.Context) {
	posts, err := database.GetPosts(a.DB)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, posts)
}

func (a *APIEnv) CreatePost(context *gin.Context) {
	var newPost models.Post
	if err := bindInput(context, &newPost); err != nil {
		return
	}
	newPost.UserID = uuid.MustParse(context.Params.ByName("userID"))
	post, err := database.CreatePost(a.DB, &newPost)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, post)
}

func bindInput(context *gin.Context, obj any) error {
	if bindErr := context.ShouldBindJSON(obj); bindErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return bindErr
	}
	return nil
}

func (a *APIEnv) GetPostByID(context *gin.Context) {
	postID := context.Param("id")
	post, err := database.GetPostByID(a.DB, postID)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, post)
}

func (a *APIEnv) UpdatePost(context *gin.Context) {
	postID := context.Param("id")
	var inputUpdate models.Post
	if bindErr := bindInput(context, &inputUpdate); bindErr != nil {
		return
	}
	post, err := database.UpdatePost(a.DB, &inputUpdate, postID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"error": postNotFoundMessage})
		return
	} else if errors.Is(err, database.ErrNotOwner) {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	} else if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, post)
}

func (a *APIEnv) UpdateUserPicture(context *gin.Context) {
	// userID := context.Param("userID")
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Open the file
	openedFile, err := file.Open()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer openedFile.Close()

	bucket := a.GoogleCloud.Bucket("skillnet-posts")

	// Create a new writer in the bucket
	ctx := context.Request.Context()
	writer := bucket.Object("test").NewWriter(ctx)

	// Copy the file to the bucket
	_, err = io.Copy(writer, openedFile)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Close the writer
	if err := writer.Close(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the URL of the new object
	// attrs, err := writer.Attrs()
	// if err != nil {
	// 	context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	attrs := writer.Attrs()
	url := attrs.MediaLink

	context.JSON(http.StatusOK, gin.H{"url": url})
}

func (a *APIEnv) DeletePost(context *gin.Context) {
	postID := context.Param("id")
	userID := context.Param("userID")
	err := database.DeletePost(a.DB, postID, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		context.JSON(http.StatusNotFound, gin.H{"error": postNotFoundMessage})
		return
	} else if errors.Is(err, database.ErrNotOwner) {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	} else if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, nil)
}
