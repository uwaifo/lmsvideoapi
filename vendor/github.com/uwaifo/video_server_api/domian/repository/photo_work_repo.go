package repository

import (
	"github.com/uwaifo/video_server_api/domian/entity"
)

// PhotoWorkRepository  . . .
type PhotoWorkRepository interface {
	SavePhotoWork(*entity.PhotoWork) (*entity.PhotoWork, map[string]string)
	GetPhotoWork(uint64) (*entity.PhotoWork, error)
	GetAllPhotoWork() ([]entity.PhotoWork, error)
	UpdatePhotoWork(*entity.PhotoWork) (*entity.PhotoWork, map[string]string)
	//DeletePhotoWork(uint64) error
}
