/*
Contains controllers for Post API.
*/
package controllers

import (
	"errors"
	"io"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/models"
	"gorm.io/gorm"
)

const (
	postNotFoundMessage = "Post not found"
)

type APIEnv struct {
	DB          *gorm.DB
	GoogleCloud *storage.Client
}

func (a *APIEnv) InitialisePostHandler() {
	a.PostDBHandler = &database.PostDB{
		DB: a.DB,
	}
}

func (a *APIEnv) CreatePost(ctx *gin.Context) {
	var newPost models.Post

	// If unable to bind JSON in request to the Post object, return status
	// code 400 Bad Request
	if err := helpers.BindInput(ctx, &newPost); err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Add userID into the corresponding field in the newPost object so that
	// the client does not have to pass in any userID, and overwrites any userID
	// that a malicious client might have passed in.
	userID := helpers.GetUserIdFromContext(ctx)
	newPost.UserID = uuid.MustParse(userID)

	post, err := a.PostDBHandler.CreatePost(&newPost)

	// If post cannot be created, return status code 500 Internal Service Error
	if err != nil {
		helpers.OutputError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.OutputData(ctx, post)
}

func (a *APIEnv) DeletePost(ctx *gin.Context) {
	postID := helpers.GetPostIdFromContext(ctx)
	userID := helpers.GetUserIdFromContext(ctx)

	err := a.PostDBHandler.DeletePost(postID, userID)
	// If post cannot be found in the database return status code 404 Status Not Found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		helpers.OutputError(ctx, http.StatusNotFound, postNotFoundMessage)
		return
	}
	// If user is not the owner of the post, return status code 401 Unauthorized
	if errors.Is(err, database.ErrNotOwner) {
		helpers.OutputError(ctx, http.StatusUnauthorized, database.ErrNotOwner.Error())
		return
	}
	// If post cannot be deleted for any other reason, return status code 403 Bad Request
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	helpers.OutputMessage(ctx, "Post successfully deleted")
}

func (a *APIEnv) GetPosts(ctx *gin.Context) {
	posts, err := a.PostDBHandler.GetPosts()
	// If unable to retrieve posts, return status code 404 Not Found
	if err != nil {
		helpers.OutputError(ctx, http.StatusNotFound, postNotFoundMessage)
		return
	}
	helpers.OutputData(ctx, posts)
}

func (a *APIEnv) GetPostByID(ctx *gin.Context) {
	postID := helpers.GetPostIdFromContext(ctx)
	post, err := a.PostDBHandler.GetPostByID(postID)
	// If unable to retrieve post, return status code 404 Not Found
	if err != nil {
		helpers.OutputError(ctx, http.StatusNotFound, postNotFoundMessage)
		return
	}
	helpers.OutputData(ctx, post)
}

func (a *APIEnv) UpdatePost(ctx *gin.Context) {
	postID := helpers.GetPostIdFromContext(ctx)
	var inputUpdate models.Post

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
	// If unable to bind JSON in request to the Post object, return status
	// code 400 Bad Request
	if err := helpers.BindInput(ctx, &inputUpdate); err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, err.Error())
	}

	post, err := a.PostDBHandler.UpdatePost(&inputUpdate, postID)

	// If post cannot be found in the database, return status code 404 Status Not Found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		helpers.OutputError(ctx, http.StatusNotFound, postNotFoundMessage)
		return
	}
	// If user is not the owner of the post, return status code 401 Unauthorised
	if errors.Is(err, database.ErrNotOwner) {
		helpers.OutputError(ctx, http.StatusUnauthorized, database.ErrNotOwner.Error())
		return
	}
	// If post cannot be updated for any other reason, return status code 403 Bad Request
	if err != nil {
		helpers.OutputError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	helpers.OutputData(ctx, post)
}
