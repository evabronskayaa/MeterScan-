package schema

import (
	"backend/internal/proto"
	"github.com/jackc/pgx/v5/pgtype"
)

type Prediction struct {
	ID        uint64      `gorm:"primarykey" json:"-"`
	UserID    uint64      `json:"-"`
	ImageName pgtype.UUID `gorm:"type:uuid;default:uuid_generate_v4();unique" json:"image_name"`

	PredictionInfos []PredictionInfo `gorm:"foreignKey:PredictionID" json:"results"`
}

type PredictionInfo struct {
	ID           uint64 `gorm:"primarykey" json:"id"`
	PredictionID uint64 `json:"-"`

	Scope              Scope   `gorm:"embedded;embeddedPrefix:scope_" json:"scope"`
	MeterReadings      string  `gorm:"type:varchar(20);not null" json:"meter_readings"`
	ValidMeterReadings string  `gorm:"type:varchar(20);not null" json:"valid_meter_readings"`
	Metric             float32 `json:"metric"`
}

type Scope struct {
	X1 uint32 `json:"x1"`
	Y1 uint32 `json:"y1"`
	X2 uint32 `json:"x2"`
	Y2 uint32 `json:"y2"`
}

func (i *PredictionInfo) Proto() *proto.RecognitionResult {
	return &proto.RecognitionResult{
		Id:          &i.ID,
		Metric:      i.Metric,
		Recognition: i.ValidMeterReadings,
		Scope: &proto.Scope{
			X1: i.Scope.X1,
			Y1: i.Scope.Y1,
			X2: i.Scope.X2,
			Y2: i.Scope.Y2,
		},
	}
}

func (p *Prediction) Proto() *proto.Prediction {
	value, _ := p.ImageName.Value()
	imageName := value.(string)
	var results []*proto.RecognitionResult
	for _, value := range p.PredictionInfos {
		results = append(results, value.Proto())
	}
	return &proto.Prediction{
		Id:        p.ID,
		ImageName: imageName,
		Results:   results,
	}
}
