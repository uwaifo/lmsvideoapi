package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uwaifo/lmsvideoapi/application"
	"net/http"
	"path/filepath"
)

type Snippet struct {
	snippetApp application.SnippetAppInterface
	userApp    application.UserAppInterface
}

// SnippetConstructor
func NewSnippet(snippApp application.SnippetAppInterface) *Snippet {
	return &Snippet{
		snippetApp: snippApp,
	}

}

func SaveSnippet(c *gin.Context) {
	snippet_file := c.PostForm("snippet_file")

	c.String(http.StatusOK, snippet_file)
	name := c.PostForm("name")
	email := c.PostForm("email")

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	files := form.File["files"]

	for _, file := range files {
		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
	}

	c.String(http.StatusOK, fmt.Sprintf("Uploaded successfully %d files with fields name=%s and email=%s.", len(files), name, email))

}
