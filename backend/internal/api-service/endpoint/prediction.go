package endpoint

import (
	"backend/internal/crud"
	"backend/internal/errors"
	"backend/internal/media"
	"backend/internal/proto"
	"backend/internal/schema"
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

const (
	ErrInvalidForm          errors.SimpleError = "Отсутствуют файлы"
	ErrNotFoundFiles        errors.SimpleError = "Файлы не найдены"
	ErrCreateStream         errors.SimpleError = "Произошла ошибка при создании запроса для распознавания"
	ErrWithReadFile         errors.SimpleError = "Произошла ошибка при чтении файла"
	ErrCannotSendFileToGRPC errors.SimpleError = "Произошла ошибка при отправке файла на сервер распознавания"
)

// PredictHandler godoc
//
//	@Summary		Предсказать цифры
//	@Tags			prediction
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			files	formData	file			true	"Файлы"
//	@Success		200		{object}	[]schema.Prediction
//	@Failure		400		{object}	string
//	@Failure		500		{object}	string
//	@Router			/predict [post]
//	@Security JWT
func (s *Service) PredictHandler(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidForm)
		return
	}
	files := form.File["files"]

	if len(files) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrNotFoundFiles)
		return
	}

	stream, err := s.ImageProcessingService.ProcessImage(context.Background())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrCreateStream)
		return
	}

	done := make(chan bool)

	images := make([][]byte, len(files))

	for i, file := range files {
		buf := bytes.NewBuffer(nil)
		open, err := file.Open()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrWithReadFile)
			return
		}
		if _, err := io.Copy(buf, open); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrWithReadFile)
			return
		}
		images[i] = buf.Bytes()
		req := &proto.ImageRequest{
			Image: images[i],
		}
		if err := stream.Send(req); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrCannotSendFileToGRPC)
			return
		}
	}
	if err := stream.CloseSend(); err != nil {
		log.Println(err)
	}

	user := c.MustGet("user").(*schema.User)

	var predictions = make([]*schema.Prediction, len(images))

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(done)
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}

			prediction, err := crud.AddPrediction(s.DB, user, resp.RecognitionResult, resp.Metric)
			if err != nil {
				log.Fatalf("can not add prediction %v", err)
				return
			}
			predictions[resp.Index] = prediction
			fileName := fmt.Sprintf("%v", prediction.ID)
			if _, err := media.SaveData("prediction", fileName, images[resp.Index]); err != nil {
				log.Fatalf("can not save %v", err)
			}
			if _, err := media.SaveData("prediction/contour", fileName, resp.ImageWithContour); err != nil {
				log.Fatalf("can not save %v", err)
			}
		}
	}()

	<-done
	c.AbortWithStatusJSON(http.StatusOK, predictions)
}
