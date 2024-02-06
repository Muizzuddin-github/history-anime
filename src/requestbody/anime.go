package requestbody

type Anime struct {
	Name        string   `json:"name" form:"name" validate:"required"`
	Genre       []string `json:"genre" form:"genre" validate:"required"`
	Description string   `json:"description" form:"description" validate:"required"`
	Image       string   `json:"image" form:"image"`
	Status      string   `json:"status" form:"status" validate:"required"`
}
