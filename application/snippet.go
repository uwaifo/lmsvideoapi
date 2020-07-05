package application

import (
	"github.com/uwaifo/lmsvideoapi/domain/entity"
	"github.com/uwaifo/lmsvideoapi/domain/repository"
)

type snippetApp struct {
	sn repository.SnippetRepository
}

var _ SnippetAppInterface = &snippetApp{}

type SnippetAppInterface interface {
	SaveSnippet(*entity.Snippet) (*entity.Snippet, map[string]string)
}

func (sn *snippetApp) SaveSnippet(snippet *entity.Snippet) (*entity.Snippet, map[string]string) {
	return sn.sn.SaveSnippet(snippet)
}
