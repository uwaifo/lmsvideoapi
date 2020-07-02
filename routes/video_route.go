package routes

import (
	"fmt"
	fileupload "github.com/uwaifo/lmsvideoapi/interfaces/upload"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uwaifo/lmsvideoapi/application"
	"github.com/uwaifo/lmsvideoapi/domain/entity"
	"github.com/uwaifo/lmsvideoapi/infrastructure/auth"
)

type Video struct {
	videoApp    application.VideoAppInterface
	userApp    application.UserAppInterface
	fileUpload fileupload.UploadFileInterface
	tk         auth.TokenInterface
	rd         auth.AuthInterface
}

//Video constructor
func NewVideo(fApp application.VideoAppInterface, uApp application.UserAppInterface, fd fileupload.UploadFileInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Video {
	return &Video{
		videoApp:    fApp,
		userApp:    uApp,
		fileUpload: fd,
		rd:         rd,
		tk:         tk,
	}
}

func (fo *Video) SaveVideo(c *gin.Context) {
	//check is the user is authenticated first
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//lookup the metadata in redis:
	userId, err := fo.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//We we are using a frontend(vuejs), our errors need to have keys for easy checking, so we use a map to hold our errors
	var saveVideoError = make(map[string]string)

	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "Invalid json",
		})
		return
	}
	//We initialize a new video for the purpose of validating: in case the payload is empty or an invalid data type is used
	emptyVideo := entity.Video{}
	emptyVideo.Title = title
	emptyVideo.Description = description
	saveVideoError = emptyVideo.Validate("")
	if len(saveVideoError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, saveVideoError)
		return
	}
	file, err := c.FormFile("video_image")
	if err != nil {
		saveVideoError["invalid_file"] = "a valid file is required"
		c.JSON(http.StatusUnprocessableEntity, saveVideoError)
		return
	}
	//check if the user exist
	_, err = fo.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}
	uploadedFile, err := fo.fileUpload.UploadFile(file)
	if err != nil {
		saveVideoError["upload_err"] = err.Error() //this error can be any we defined in the UploadFile method
		c.JSON(http.StatusUnprocessableEntity, saveVideoError)
		return
	}
	var video = entity.Video{}
	video.UserID = userId
	video.Title = title
	video.Description = description
	video.VideoImage = uploadedFile
	savedVideo, saveErr := fo.videoApp.SaveVideo(&video)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}
	c.JSON(http.StatusCreated, savedVideo)
}

func (fo *Video) UpdateVideo(c *gin.Context) {
	//Check if the user is authenticated first
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	//lookup the metadata in redis:
	userId, err := fo.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//We we are using a frontend(vuejs), our errors need to have keys for easy checking, so we use a map to hold our errors
	var updateVideoError = make(map[string]string)

	videoId, err := strconv.ParseUint(c.Param("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	//Since it is a multipart form data we sent, we will do a manual check on each item
	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json")
	}
	//We initialize a new video for the purpose of validating: in case the payload is empty or an invalid data type is used
	emptyVideo := entity.Video{}
	emptyVideo.Title = title
	emptyVideo.Description = description
	updateVideoError = emptyVideo.Validate("update")
	if len(updateVideoError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, updateVideoError)
		return
	}
	user, err := fo.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	//check if the video exist:
	video, err := fo.videoApp.GetVideo(videoId)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	//if the user id doesnt match with the one we have, dont update. This is the case where an authenticated user tries to update someone else post using postman, curl, etc
	if user.ID != video.UserID {
		c.JSON(http.StatusUnauthorized, "you are not the owner of this video")
		return
	}
	//Since this is an update request,  a new image may or may not be given.
	// If not image is given, an error occurs. We know this that is why we ignored the error and instead check if the file is nil.
	// if not nil, we process the file by calling the "UploadFile" method.
	// if nil, we used the old one whose path is saved in the database
	file, _ := c.FormFile("video_image")
	if file != nil {
		video.VideoImage, err = fo.fileUpload.UploadFile(file)
		//since i am using Digital Ocean(DO) Spaces to save image, i am appending my DO url here. You can comment this line since you may be using Digital Ocean Spaces.
		video.VideoImage = os.Getenv("DO_SPACES_URL") + video.VideoImage
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"upload_err": err.Error(),
			})
			return
		}
	}
	//we dont need to update user's id
	video.Title = title
	video.Description = description
	video.UpdatedAt = time.Now()
	updatedVideo, dbUpdateErr := fo.videoApp.UpdateVideo(video)
	if dbUpdateErr != nil {
		c.JSON(http.StatusInternalServerError, dbUpdateErr)
		return
	}
	c.JSON(http.StatusOK, updatedVideo)
}

func (fo *Video) GetAllVideo(c *gin.Context) {
	allvideo, err := fo.videoApp.GetAllVideo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, allvideo)
}

func (fo *Video) GetVideoAndCreator(c *gin.Context) {
	videoId, err := strconv.ParseUint(c.Param("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	video, err := fo.videoApp.GetVideo(videoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	user, err := fo.userApp.GetUser(video.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	videoAndUser := map[string]interface{}{
		"video":    video,
		"creator": user.PublicUser(),
	}
	c.JSON(http.StatusOK, videoAndUser)
}

func (fo *Video) DeleteVideo(c *gin.Context) {
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	videoId, err := strconv.ParseUint(c.Param("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	_, err = fo.userApp.GetUser(metadata.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	err = fo.videoApp.DeleteVideo(videoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "video deleted")
}
