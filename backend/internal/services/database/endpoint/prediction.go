package endpoint

import (
	"backend/internal/proto"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s GRPCServer) GetPredictions(ctx context.Context, request *proto.GetPredictionsRequest) (*proto.GetPredictionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPredictions not implemented")
}
func (s GRPCServer) UpdatePrediction(ctx context.Context, request *proto.UpdatePredictionRequest) (*proto.Empty, error) {

	/*
		db.Model(&schema.PredictionInfo{}).
				Joins("JOIN predictions ON predictions.id = prediction_infos.prediction_id").
				Where("prediction_infos.id = ? AND predictions.user_id = ?", id, userId).
				Update("valid_meter_readings", meterReadings).Error
	*/

	return nil, status.Errorf(codes.Unimplemented, "method UpdatePrediction not implemented")
}
func (s GRPCServer) AddPrediction(ctx context.Context, request *proto.AddPredictionRequest) (*proto.Prediction, error) {

	/*
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
	*/

	return nil, status.Errorf(codes.Unimplemented, "method AddPrediction not implemented")
}
