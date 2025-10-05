package httpserver

type RespErr struct {
	Error  string `json:"error,omitempty"`
	ErrorMessage  string `json:"error_message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

func Error(err string) RespErr {
	return RespErr{
		Error: err,
	}
}
