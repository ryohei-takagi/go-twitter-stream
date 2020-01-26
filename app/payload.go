package app

type (
	Payload struct {
		ID   int      `json:"id"`
		Body interface{} `json:"body"`
	}

	MessagePayload struct {
		Message string `json:"message"`
	}
)
