package middleware

import (
	"encoding/json"
	"wb-l0/internal/service"
)

// ProcessModel check model and return OrderUid
func ProcessModel(body []byte) (string, error) {
	var model service.Model
	err := json.Unmarshal(body, &model)
	if err != nil {
		return "", &service.MyError{
			Message: err.Error(),
			Code:    401,
		}
	}

	if model.OrderUid == "" || model.TrackNumber == "" || model.Payment.Transaction == "" {
		return "", &service.MyError{
			Message: "Uncorrected model",
			Code:    401,
		}
	}
	return model.OrderUid, nil
}
