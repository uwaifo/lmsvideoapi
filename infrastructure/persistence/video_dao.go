package persistence

import (
	"errors"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/uwaifo/lmsvideoapi/domain/entity"
	"github.com/uwaifo/lmsvideoapi/domain/repository"
   )

type VideoRepo struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepo {
	return &VideoRepo{db}
}

//VideoRepo implements the repository.VideoRepository interface
var _ repository.VideoRepository = &VideoRepo{}

func (r *VideoRepo) SaveVideo(video *entity.Video) (*entity.Video, map[string]string) {
	dbErr := map[string]string{}
	//The images are uploaded to digital ocean spaces. So we need to prepend the url. This might not be your use case, if you are not uploading image to Digital Ocean.
	video.VideoImage = os.Getenv("DO_SPACES_URL") + video.VideoImage

	err := r.db.Debug().Create(&video).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "video title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return video, nil
}

func (r *VideoRepo) GetVideo(id uint64) (*entity.Video, error) {
	var video entity.Video
	err := r.db.Debug().Where("id = ?", id).Take(&video).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("video not found")
	}
	return &video, nil
}

func (r *VideoRepo) GetAllVideo() ([]entity.Video, error) {
	var videos []entity.Video
	err := r.db.Debug().Limit(100).Order("created_at desc").Find(&videos).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return videos, nil
}

func (r *VideoRepo) UpdateVideo(video *entity.Video) (*entity.Video, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Save(&video).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return video, nil
}

func (r *VideoRepo) DeleteVideo(id uint64) error {
	var video entity.Video
	err := r.db.Debug().Where("id = ?", id).Delete(&video).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}
