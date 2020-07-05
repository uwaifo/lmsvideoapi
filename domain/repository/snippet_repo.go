package repository

import "github.com/uwaifo/lmsvideoapi/domain/entity"

type SnippetRepository interface {
	SaveSnippet(*entity.Snippet) (*entity.Snippet, map[string]string)
}
