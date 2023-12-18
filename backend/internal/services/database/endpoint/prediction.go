package endpoint

import (
	"backend/internal/errors"
	"backend/internal/proto"
	"backend/internal/services/database/schema"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s GRPCServer) GetPredictions(ctx context.Context, request *proto.GetPredictionsRequest) (*proto.GetPredictionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPredictions not implemented")
}

func (s GRPCServer) UpdatePrediction(ctx context.Context, request *proto.UpdatePredictionRequest) (*proto.Empty, error) {
	if err := s.DB.Model(&schema.PredictionInfo{}).
		Joins("JOIN predictions ON predictions.id = prediction_infos.prediction_id").
		Where("prediction_infos.id = ? AND predictions.user_id = ?", request.Id, request.UserId).
		Update("valid_meter_readings", request.ValidMeterReadings).Error; err != nil {
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
