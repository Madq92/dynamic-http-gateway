package common

type Payload struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

const SUCCESS_CODE = 0

func (p *Payload) IsOk() bool {
	return p.Code == SUCCESS_CODE
}
