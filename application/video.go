package application

import (
	"github.com/uwaifo/lmsvideoapi/domain/entity"
	"github.com/uwaifo/lmsvideoapi/domain/repository"
)

type videoApp struct {
	fr repository.VideoRepository
}

var _ VideoAppInterface = &videoApp{}

type VideoAppInterface interface {
	SaveVideo(*entity.Video) (*entity.Video, map[string]string)
	GetAllVideo() ([]entity.Video, error)
	GetVideo(uint64) (*entity.Video, error)
	UpdateVideo(*entity.Video) (*entity.Video, map[string]string)
	DeleteVideo(uint64) error
}

func (f *videoApp) SaveVideo(video *entity.Video) (*entity.Video, map[string]string) {
	return f.fr.SaveVideo(video)
}

func (f *videoApp) GetAllVideo() ([]entity.Video, error) {
	return f.fr.GetAllVideo()
}

func (f *videoApp) GetVideo(videoId uint64) (*entity.Video, error) {
	return f.fr.GetVideo(videoId)
}

func (f *videoApp) UpdateVideo(video *entity.Video) (*entity.Video, map[string]string) {
	return f.fr.UpdateVideo(video)
}

func (f *videoApp) DeleteVideo(videoId uint64) error {
	return f.fr.DeleteVideo(videoId)
}
