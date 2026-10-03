package common

type GenericResponse struct {
	Error   error  `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
