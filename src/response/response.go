package response

type Errors struct {
	Errors []string `json:"errors"`
}

type Msg struct {
	Message string `json:"message"`
}