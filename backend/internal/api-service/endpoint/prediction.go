package endpoint

import (
	"backend/internal/crud"
	"backend/internal/errors"
	"backend/internal/media"
	"backend/internal/proto"
	"backend/internal/schema"
	"backend/internal/schema/dto"
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"sync"
)

const (
	ErrInvalidForm   errors.SimpleError = "Отсутствуют файлы"
	ErrNotFoundFiles errors.SimpleError = "Файлы не найдены"
	ErrCreateStream  errors.SimpleError = "Произошла ошибка при создании запроса для распознавания"
)

// PredictHandler godoc
//
//	@Summary		Предсказать цифры
//	@Tags			predictions
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			files	formData	file			true	"Файлы"
//	@Success		200		{object}	[]schema.Prediction
//	@Failure		400		{object}	string
//	@Failure		500		{object}	string
//	@Router			/predictions [post]
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

	user := c.MustGet("user").(*schema.User)
	var wg sync.WaitGroup
	predictions := make([]*schema.Prediction, 0)
	errors := make(chan error)

	for i, file := range files {
		wg.Add(1)
		go func(i int, file *multipart.FileHeader) {
			defer wg.Done()
			open, err := file.Open()
			if err != nil {
				errors <- fmt.Errorf("file open error: %v", err)
				return
			}
			defer open.Close()

			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, open); err != nil {
				errors <- fmt.Errorf("file read error: %v", err)
				return
			}

			req := &proto.ImageRequest{Image: buf.Bytes()}
			if err := stream.Send(req); err != nil {
				errors <- fmt.Errorf("stream send error: %v", err)
				return
			}

			resp, err := stream.Recv()
			if err == io.EOF {
				return
			} else if err != nil {
				errors <- fmt.Errorf("stream receive error: %v", err)
				return
			}

			prediction, err := crud.AddPrediction(s.DB, user, resp.Results)
			if err != nil {
				errors <- fmt.Errorf("add prediction error: %v", err)
				return
			}
			value, _ := prediction.ImageName.Value()
			fileName := value.(string)
			predictions = append(predictions, prediction)
			if _, err := media.SaveData("prediction", fileName, buf.Bytes()); err != nil {
				errors <- fmt.Errorf("file save error: %v", err)
				return
			}
		}(i, file)
	}

	wg.Wait()

	if err := stream.CloseSend(); err != nil {
		log.Println("Failed to close stream:", err)
	}

	close(errors)

	if len(errors) > 0 {
		var errMessages []string
		for err := range errors {
			errMessages = append(errMessages, err.Error())
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": errMessages})
		return
	}

	c.JSON(http.StatusOK, predictions)
}

// UpdatePredictHandler godoc
//
//	@Summary		Подтвердить показания
//	@Tags			predictions
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			prediction	body		dto.UpdatePredictionForm	true	"UpdatePredictionForm"
//	@Success		200
//	@Failure		400		{object}	string
//	@Failure		500		{object}	string
//	@Router			/predictions [put]
//	@Security JWT
func (s *Service) UpdatePredictHandler(c *gin.Context) {
	var form dto.UpdatePredictionForm
	if err := c.ShouldBind(&form); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	user := c.MustGet("user").(*schema.User)

	if err := crud.UpdateMeterReadings(s.DB, form.ID, user.ID, form.MeterReadings); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	} else {
		c.Status(http.StatusOK)
	}
}

// GetPredictionsHandler godoc
//
//	@Summary		Достать показания пользователя
//	@Tags			predictions
//	@Accept			json
//	@Produce		json
//	@Param			paging	query		dto.Paging	false	"Paging params"
//	@Success		200		{object}	paginate.Page
//	@Failure		400		{object}	string
//	@Failure		500		{object}	string
//	@Router			/predictions [get]
//	@Security JWT
func (s *Service) GetPredictionsHandler(c *gin.Context) {
	user := c.MustGet("user").(*schema.User)

	model := s.DB.Model(&user.Predictions).Preload("PredictionInfos")
	page := s.Pagination.With(model).Request(c.Request).Response(&[]schema.Prediction{})

	c.JSON(http.StatusOK, page)
}
