package jsend

type Status int8

const (
	StatusError   = "error"
	StatusFail    = "fail"
	StatusSuccess = "success"
)

type Payload struct {
	Status  string `json:"status"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func Ok(payload any) Payload {
	return Payload{Status: StatusSuccess, Data: payload}
}

func Error(payload string) Payload {
	return Payload{Status: StatusError, Message: payload}
}

func Fail(payload any) Payload {
	return Payload{Status: StatusFail, Data: payload}
}
