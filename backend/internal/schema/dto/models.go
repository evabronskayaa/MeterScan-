package dto

type RecognitionResult struct {
	ID uint `json:"id"`
}

type Paging struct {
	Sort   string `json:"sort"`
	Page   string `json:"page"`
	Size   string `json:"size"`
	Order  string `json:"order"`
	Filter string `json:"filter"`
}
