package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uwaifo/lmsvideoapi/application"
	fileupload "github.com/uwaifo/lmsvideoapi/interfaces/upload"
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
	snippetStart := c.PostForm("start")
	snippetEnd := c.PostForm("end")

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	files := form.File["files"]
	fmt.Println(files)
	fmt.Println(form)

	singleFile := filepath.Base(files[0].Filename)

	for _, file := range files {
		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, "temp_upload/"+filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
	}
	// send snippet to be edited
	editedSnippet, _ := fileupload.SingleSnippetEdit(singleFile, snippetStart, snippetEnd)

	c.String(http.StatusCreated, fmt.Sprintf("Uploaded successfully %d files with fields start=%s and end=%s and outputs %s.", len(files),
		snippetStart, snippetEnd, editedSnippet))

}
