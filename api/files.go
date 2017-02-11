package api

import (
	"bytes"
	"fmt"

	"io"
	"io/ioutil"

	"image"

	"github.com/disintegration/imaging"
	uuid "github.com/satori/go.uuid"

	"image/jpeg"

	"net/http"

	"os"

	"github.com/xDarkicex/FileServer/models"
	"github.com/xDarkicex/FileServer/server"
)

type apiFiles struct{}

var files = new(apiFiles)
var ImageParams = map[string]ImageResizeParams{
	"original": {0, 0, 0},
	"800x800":  {800, 800, 80},
	"400x400":  {400, 400, 80},
	"150x150":  {150, 150, 80},
	"50x50":    {50, 50, 80},
}

type FileResponse struct {
	ID   int               `json:"id"`
	URLs map[string]string `json:"urls"`
}

type ImageResizeParams struct {
	Width   int
	Height  int
	Quality int
}

func handleFiles(request *server.Router) {
	request.GET("/{name}", files.show)
	request.POST("", files.upload)
	request.GET("", files.index)
}
func getURLS(dir string) map[string]string {
	urls := map[string]string{}
	for k := range ImageParams {
		urls[k] = fmt.Sprintf("http://127.0.0.1:8080/image/%s/%s.jpg", dir, k)
	}
	return urls
}

func (apiFiles) response(file *models.File) *FileResponse {
	return &FileResponse{
		ID:   file.ID,
		URLs: getURLS(file.Name),
	}
}
func (apiFiles *apiFiles) multiResponce(files []*models.File) []*FileResponse {
	fileResponse := []*FileResponse{}
	for _, file := range files {
		response := apiFiles.response(file)
		if response != nil {
			fileResponse = append(fileResponse, response)
		}
	}
	return fileResponse
}

type HTTPFile interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

func resizeAndStore(dir string, file HTTPFile, name string, params ImageResizeParams) error {
	if name == "original" {
		file.Seek(0, 0)
		buf := &bytes.Buffer{}
		_, err := buf.ReadFrom(file)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(fmt.Sprintf("assets/image/%s/%s.jpg", dir, name), buf.Bytes(), 0777)
		if err != nil {
			return err
		}
	} else {
		file.Seek(0, 0)
		img, _, err := image.Decode(file)
		if err != nil {
			return err
		}
		resizedImage := imaging.Fit(img, params.Width, params.Height, imaging.Lanczos)
		var b []byte
		buf := bytes.NewBuffer(b)
		jpeg.Encode(buf, resizedImage, &jpeg.Options{Quality: params.Quality})
		err = ioutil.WriteFile(fmt.Sprintf("assets/image/%s/%s.jpg", dir, name), buf.Bytes(), 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func (apiFiles *apiFiles) upload(ctx *server.Context) {
	httpFile, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.RenderError(http.StatusBadRequest, err)
		return
	}
	dir := uuid.NewV1().String()
	err = os.Mkdir(fmt.Sprintf("assets/image/%s", dir), 0777)
	if err != nil {
		ctx.RenderError(http.StatusBadRequest, err)
		return
	}
	for k, p := range ImageParams {
		err = resizeAndStore(dir, httpFile, k, p)
		if err != nil {
			ctx.RenderError(http.StatusBadRequest, err)
			return
		}
	}
	file, err := models.Files.Create(dir)
	if err != nil {
		ctx.RenderError(http.StatusBadRequest, err)
		return
	}
	ctx.RenderJSON(http.StatusCreated, apiFiles.response(file))
}

func (apiFiles *apiFiles) show(ctx *server.Context) {
	file, err := models.Files.ByName(ctx.Param("name"))
	if err != nil {
		ctx.RenderError(http.StatusBadRequest, err)
	}
	ctx.RenderJSON(http.StatusOK, apiFiles.response(file))
}

func (apiFiles *apiFiles) index(ctx *server.Context) {
	files, err := models.Files.Index()
	if err != nil {
		ctx.RenderError(http.StatusBadRequest, err)
		return
	}
	ctx.RenderJSON(http.StatusOK, apiFiles.multiResponce(files))
}
