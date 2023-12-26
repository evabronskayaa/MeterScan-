package endpoint

import (
	"backend/internal/dto"
	"backend/internal/errors"
	"backend/internal/proto"
	"backend/internal/rabbit"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

const (
	ErrInvalidForm   errors.SimpleError = "Отсутствуют файлы"
	ErrNotFoundFiles errors.SimpleError = "Файлы не найдены"
)

// PredictHandler godoc
//
//	@Summary		Предсказать цифры
//	@Tags			predictions
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			files	formData	file			true	"Файлы"
//	@Success		200		{object}	[]proto.Prediction
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

	user := c.MustGet("user").(*proto.UserResponse)
	predictions := make([]*proto.Prediction, 0)
	recognitionErrs := make(chan error)

	for _, file := range files {
		open, err := file.Open()
		if err != nil {
			recognitionErrs <- fmt.Errorf("file open error: %v", err)
			continue
		}
		defer func(open multipart.File) {
			if err := open.Close(); err != nil {
				log.Println(err)
			}
		}(open)

		prediction, err := s.DatabaseService.AddPrediction(c, &proto.AddPredictionRequest{UserId: user.Id})
		predictions = append(predictions, prediction)
		fileName := prediction.ImageName
		if info, err := s.S3Client.PutObject(context.Background(), "recognized-images", fileName, open, file.Size, minio.PutObjectOptions{
			Expires: time.Now().Add(3 * 30 * 24 * time.Hour),
		}); err != nil {
			log.Printf("file save error: %v", err)
			recognitionErrs <- fmt.Errorf("file save error: %v", err)
			continue
		} else if err := rabbit.PublishToRabbitMQ(s.RabbitMQ, "predictions", rabbit.Prediction{
			Index: prediction.Id,
			Image: fmt.Sprintf("%v/%v/%v", s.S3Client.EndpointURL(), info.Bucket, fileName),
		}); err != nil {
			log.Printf("publish to rabbit error: %v", err)
			recognitionErrs <- fmt.Errorf("publish to rabbit error: %v", err)
			continue
		}
	}

	if len(recognitionErrs) > 0 {
		var errMessages []string
		for err := range recognitionErrs {
			errMessages = append(errMessages, err.Error())
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": errMessages})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	defer func() {
		close(recognitionErrs)
	}()

	needCancel := false
	for {
		select {
		case <-ticker.C:
			for i, prediction := range predictions {
				newPrediction, err := s.DatabaseService.GetPrediction(ctx, &proto.GetPredictionsRequest{Id: prediction.Id})
				if err != nil {
					recognitionErrs <- err
					continue
				}
				predictions[i] = newPrediction
				if len(newPrediction.Results) > 0 {
					needCancel = true
				}
			}
			if needCancel {
				if len(recognitionErrs) > 0 {
					var errMessages []string
					for err := range recognitionErrs {
						errMessages = append(errMessages, err.Error())
					}
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": errMessages})
					return
				}
				c.JSON(http.StatusOK, predictions)
				return
			}
		case <-ctx.Done():
			c.Status(http.StatusRequestTimeout)
			return
		}
	}
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

	user := c.MustGet("user").(*proto.UserResponse)

	if _, err := s.DatabaseService.UpdatePrediction(c, &proto.UpdatePredictionRequest{Id: form.ID, UserId: user.Id, ValidMeterReadings: form.MeterReadings}); err != nil {
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
	user := c.MustGet("user").(*proto.UserResponse)

	if response, err := s.DatabaseService.GetPredictions(c, &proto.GetPredictionsRequest{Id: user.Id}); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, response.GetPredictions())
	}
}
