package repository

import "github.com/uwaifo/lmsvideoapi/domain/entity"

type VideoRepository interface {
	SaveVideo(*entity.Video) (*entity.Video, map[string]string)
	GetVideo(uint64) (*entity.Video, error)
	GetAllVideo() ([]entity.Video, error)
	UpdateVideo(*entity.Video) (*entity.Video, map[string]string)
	DeleteVideo(uint64) error
}
