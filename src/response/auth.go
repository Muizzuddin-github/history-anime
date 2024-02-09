package response

type Login struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
