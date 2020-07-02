package interfaces

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uwaifo/video_server_api/application"
	"github.com/uwaifo/video_server_api/domian/entity"
	"github.com/uwaifo/video_server_api/infrastructure/auth"
)

// PhotoWork . .
type PhotoWork struct {
	photoWorkApp application.PhotoWorkAppInterface
	userApp      application.UserAppInterface
	//fileUpload   fileupload.UploadFileInterface
	tk auth.TokenInterface
	rd auth.AuthInterface
}

// NewPhotoWork . . .
func NewPhotoWork(pwApp application.PhotoWorkAppInterface, uApp application.UserAppInterface, rd auth.AuthInterface, tk auth.TokenInterface) *PhotoWork {
	return &PhotoWork{
		photoWorkApp: pwApp,
		//foodApp:    fApp,
		userApp: uApp,
		//fileUpload: fd,
		rd: rd,
		tk: tk,
	}
}

// SavePhotoWork . . .
func (fo *PhotoWork) SavePhotoWork(c *gin.Context) {
	//check is the user is authenticated first
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//lookup the metadata in redis:
	userID, err := fo.rd.FetchAuth(metadata.TokenUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	//We we are using a frontend(vuejs), our errors need to have keys for easy checking, so we use a map to hold our errors
	var savePhotoWorkError = make(map[string]string)

	title := c.PostForm("title")
	description := c.PostForm("description")
	if fmt.Sprintf("%T", title) != "string" || fmt.Sprintf("%T", description) != "string" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "Invalid json",
		})
		return
	}
	//We initialize a new food for the purpose of validating: in case the payload is empty or an invalid data type is used
	emptyFood := entity.Food{}
	emptyFood.Title = title
	emptyFood.Description = description
	savePhotoWorkError = emptyFood.Validate("")
	if len(savePhotoWorkError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, savePhotoWorkError)
		return
	}
	/*file, err := c.FormFile("food_image")
	if err != nil {
		savePhotoWorkError["invalid_file"] = "a valid file is required"
		c.JSON(http.StatusUnprocessableEntity, savePhotoWorkError)
		return
	}*/
	//check if the user exist
	_, err = fo.userApp.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}
	/*uploadedFile, err := fo.fileUpload.UploadFile(file)
	if err != nil {
		savePhotoWorkError["upload_err"] = err.Error() //this error can be any we defined in the UploadFile method
		c.JSON(http.StatusUnprocessableEntity, savePhotoWorkError)
		return
	}*/
	var photoWork = entity.PhotoWork{}
	photoWork.UserID = userID
	photoWork.Title = title
	photoWork.Description = description
	//food.FoodImage = uploadedFile
	savedPhotoWork, saveErr := fo.photoWorkApp.SavePhotoWork(&photoWork)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, saveErr)
		return
	}
	c.JSON(http.StatusCreated, savedPhotoWork)
}

// GetUserPhotoWork  . .
func (fo *PhotoWork) GetUserPhotoWork(c *gin.Context) {
	photoWorkID, err := strconv.ParseUint(c.Param("photo_work_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	photoWork, err := fo.photoWorkApp.GetPhotoWork(photoWorkID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	user, err := fo.userApp.GetUser(photoWork.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	foodAndUser := map[string]interface{}{
		"work":    photoWork,
		"creator": user.PublicUser(),
	}
	c.JSON(http.StatusOK, foodAndUser)
}

// GetAllPhotoWork . . .
func (fo *PhotoWork) GetAllPhotoWork(c *gin.Context) {
	allPhotoWork, err := fo.photoWorkApp.GetAllPhotoWork()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, allPhotoWork)
}

// UpdatePhotoWork . . .
func (fo *PhotoWork) UpdatePhotoWork(c *gin.Context) {
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
	var updatePhotoWorkError = make(map[string]string)

	photoWorkID, err := strconv.ParseUint(c.Param("photo_work_id"), 10, 64)
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
	//We initialize a new food for the purpose of validating: in case the payload is empty or an invalid data type is used
	emptyPhotoWork := entity.PhotoWork{}
	emptyPhotoWork.Title = title
	emptyPhotoWork.Description = description
	updatePhotoWorkError = emptyPhotoWork.Validate("update")
	if len(updatePhotoWorkError) > 0 {
		c.JSON(http.StatusUnprocessableEntity, updatePhotoWorkError)
		return
	}
	user, err := fo.userApp.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user not found, unauthorized")
		return
	}

	//check if the food exist:
	photoWork, err := fo.photoWorkApp.GetPhotoWork(photoWorkID)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	//if the user id doesnt match with the one we have, dont update. This is the case where an authenticated user tries to update someone else post using postman, curl, etc
	if user.ID != photoWork.UserID {
		c.JSON(http.StatusUnauthorized, "you are not the owner of this photo work")
		return
	}
	//Since this is an update request,  a new image may or may not be given.
	// If not image is given, an error occurs. We know this that is why we ignored the error and instead check if the file is nil.
	// if not nil, we process the file by calling the "UploadFile" method.
	// if nil, we used the old one whose path is saved in the database
	/*file, _ := c.FormFile("food_image")
	if file != nil {
		food.FoodImage, err = fo.fileUpload.UploadFile(file)
		//since i am using Digital Ocean(DO) Spaces to save image, i am appending my DO url here. You can comment this line since you may be using Digital Ocean Spaces.
		food.FoodImage = os.Getenv("DO_SPACES_URL") + food.FoodImage
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"upload_err": err.Error(),
			})
			return
		}
	}
	*/
	//we dont need to update user's id
	photoWork.Title = title
	photoWork.Description = description
	photoWork.UpdatedAt = time.Now()
	updatedPhotoWork, dbUpdateErr := fo.photoWorkApp.UpdatePhotoWork(photoWork)
	if dbUpdateErr != nil {
		c.JSON(http.StatusInternalServerError, dbUpdateErr)
		return
	}
	c.JSON(http.StatusOK, updatedPhotoWork)
}

/*


func (fo *Food) DeleteFood(c *gin.Context) {
	metadata, err := fo.tk.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	foodId, err := strconv.ParseUint(c.Param("food_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}
	_, err = fo.userApp.GetUser(metadata.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	err = fo.foodApp.DeleteFood(foodId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "food deleted")
}
*/
