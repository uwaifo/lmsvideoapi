package application

import (
	"github.com/uwaifo/video_server_api/domian/entity"
	"github.com/uwaifo/video_server_api/domian/repository"
)

type photoWorkApp struct {
	pwr repository.PhotoWorkRepository
}

var _ PhotoWorkAppInterface = &photoWorkApp{}

// PhotoWorkAppInterface . .
type PhotoWorkAppInterface interface {
	SavePhotoWork(*entity.PhotoWork) (*entity.PhotoWork, map[string]string)
	GetPhotoWork(uint64) (*entity.PhotoWork, error)
	GetAllPhotoWork() ([]entity.PhotoWork, error)
	UpdatePhotoWork(*entity.PhotoWork) (*entity.PhotoWork, map[string]string)
	//DeletePhotoWork(uint64) error

}

func (p *photoWorkApp) SavePhotoWork(photoWork *entity.PhotoWork) (*entity.PhotoWork, map[string]string) {
	return p.pwr.SavePhotoWork(photoWork)
}

func (p *photoWorkApp) GetPhotoWork(photoWorkID uint64) (*entity.PhotoWork, error) {
	return p.pwr.GetPhotoWork(photoWorkID)
}

func (p *photoWorkApp) GetAllPhotoWork() ([]entity.PhotoWork, error) {
	return p.pwr.GetAllPhotoWork()
}

func (p *photoWorkApp) UpdatePhotoWork(photoWork *entity.PhotoWork) (*entity.PhotoWork, map[string]string) {
	return p.pwr.UpdatePhotoWork(photoWork)
}

/*
func (f *foodApp) DeleteFood(foodId uint64) error {
	return f.fr.DeleteFood(foodId)
}
*/
