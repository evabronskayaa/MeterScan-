package schema

import "github.com/jackc/pgx/v5/pgtype"

type Prediction struct {
	ID        uint        `gorm:"primarykey" json:"-"`
	UserID    uint        `json:"-"`
	ImageName pgtype.UUID `gorm:"type:uuid;default:uuid_generate_v4();unique" json:"image_name"`

	PredictionInfos []PredictionInfo `gorm:"foreignKey:PredictionID" json:"results"`
}

type PredictionInfo struct {
	ID           uint `gorm:"primarykey" json:"id"`
	PredictionID uint `json:"-"`

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
