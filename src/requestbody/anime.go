package requestbody

type Anime struct {
	Name        string   `json:"name" validate:"required"`
	Genre       []string `json:"genre" validate:"required,min=1"`
	Description string   `json:"description" validate:"required"`
	Image       string   `json:"image" validate:"required"`
	Status      string   `json:"status" validate:"required"`
}
