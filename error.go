package appstore

import "fmt"

type Error struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %d message: %s", e.ErrorCode, e.ErrorMessage)
}
