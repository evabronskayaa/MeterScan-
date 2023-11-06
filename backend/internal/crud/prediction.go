package crud

import (
	"backend/internal/errors"
	"backend/internal/proto"
	"backend/internal/schema"
	"gorm.io/gorm"
)

const (
	ErrSavePrediction errors.SimpleError = "Произошла ошибка при сохранении показаний"
)

func AddPrediction(db *gorm.DB, user *schema.User, results []*proto.RecognitionResult) (*schema.Prediction, error) {
	prediction := &schema.Prediction{
		UserID: user.ID,
	}

	for _, result := range results {
		predictionInfo := schema.PredictionInfo{
			Scope: schema.Scope{
				X1: result.GetScope().GetX1(),
				Y1: result.GetScope().GetY1(),
				X2: result.GetScope().GetX2(),
				Y2: result.GetScope().GetY2(),
			},
			MeterReadings:      result.GetRecognition(),
			ValidMeterReadings: result.GetRecognition(),
			Metric:             result.GetMetric(),
		}
		prediction.PredictionInfos = append(prediction.PredictionInfos, predictionInfo)
	}

	if db.Create(&prediction).Error != nil {
		return nil, ErrSavePrediction
	}
	return prediction, nil
}

func UpdateMeterReadings(db *gorm.DB, id, userId uint, meterReadings string) error {
	return db.Model(&schema.PredictionInfo{}).
		Joins("JOIN predictions ON predictions.id = prediction_infos.prediction_id").
		Where("prediction_infos.id = ? AND predictions.user_id = ?", id, userId).
		Update("valid_meter_readings", meterReadings).Error
}
