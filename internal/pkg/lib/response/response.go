package response

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = 200
	StatusError = 400
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}
