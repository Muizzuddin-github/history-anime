package requestbody

type Genre struct {
	Name string `json:"name" validate:"required"`
}
