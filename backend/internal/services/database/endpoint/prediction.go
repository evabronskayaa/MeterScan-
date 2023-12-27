package endpoint

import (
	"backend/internal/errors"
	"backend/internal/proto"
	"backend/internal/services/database/schema"
	"context"
)

func (s GRPCServer) GetPredictions(_ context.Context, request *proto.GetPredictionsRequest) (*proto.GetPredictionsResponse, error) {
	var predictions []schema.Prediction

	result := s.DB.Preload("PredictionInfos").Where("user_id = ?", request.Id).Find(&predictions)
	if result.Error != nil {
		return nil, result.Error
	}

	var protoPredictions []*proto.Prediction
	for _, prediction := range predictions {
		protoPrediction := prediction.Proto()
		protoPredictions = append(protoPredictions, protoPrediction)
	}

	response := &proto.GetPredictionsResponse{
		Predictions: protoPredictions,
	}

	return response, nil
}

func (s GRPCServer) UpdatePrediction(ctx context.Context, request *proto.UpdatePredictionRequest) (*proto.Empty, error) {
	if err := s.DB.Exec(`
    UPDATE prediction_infos 
    SET valid_meter_readings = ? 
    FROM predictions 
    WHERE prediction_infos.id = ? 
      AND prediction_infos.prediction_id = predictions.id 
      AND predictions.user_id = ?`,
		request.ValidMeterReadings, request.Id, request.UserId).Error; err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

func (s GRPCServer) AddPrediction(ctx context.Context, request *proto.AddPredictionRequest) (*proto.Prediction, error) {
	userId := request.UserId

	prediction := &schema.Prediction{
		UserID: userId,
	}

	for _, result := range request.Results {
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

	if s.DB.Create(&prediction).Error != nil {
		return nil, errors.ErrSavePrediction
	}
	return prediction.Proto(), nil
}

func (s GRPCServer) UpdateFullPrediction(ctx context.Context, request *proto.UpdateFullPredictionRequest) (*proto.Empty, error) {
	prediction := &schema.Prediction{
		ID: request.Id,
	}

	for _, result := range request.Results {
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

	if err := s.DB.Updates(prediction).Error; err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

func (s GRPCServer) GetPrediction(_ context.Context, request *proto.GetPredictionsRequest) (*proto.Prediction, error) {
	var prediction *schema.Prediction
	if s.DB.Preload("PredictionInfos").Where("id = ?", request.Id).First(&prediction).Error != nil {
		return nil, errors.ErrNotFoundPrediction
	}
	return prediction.Proto(), nil
}

func (s GRPCServer) RemovePredictionInfo(_ context.Context, request *proto.RemovePredictionInfoRequest) (*proto.Empty, error) {
	var predictionInfo *schema.PredictionInfo
	if s.DB.Where("id = ?", request.Id).First(&predictionInfo).Error != nil {
		return nil, errors.ErrNotFoundPrediction
	}

	var prediction *schema.Prediction
	if s.DB.Preload("PredictionInfos").Where("id = ?", predictionInfo.PredictionID).First(&prediction).Error != nil || prediction.UserID != request.UserId {
		return nil, errors.ErrNotFoundPrediction
	}

	if s.DB.Delete(predictionInfo).Error != nil {
		return nil, errors.ErrSavePrediction
	}

	if len(prediction.PredictionInfos) <= 1 {
		if s.DB.Delete(prediction).Error != nil {
			return nil, errors.ErrSavePrediction
		}
	}

	return &proto.Empty{}, nil
}
