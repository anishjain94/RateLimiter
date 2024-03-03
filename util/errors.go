package util

import (
	"fmt"
)

func ToError(errorCode string, statusCode int, errorMsg string) {
	panic(ToErrorString(errorCode, statusCode, errorMsg))

}

func ToErrorString(errorCode string, statusCode int, errorMsg string) string {
	return fmt.Sprint(statusCode) + ":" + string(errorCode) + ":" + errorMsg
}

func ErrorIf(condition bool, errorCode string, statusCode int, errorMsg string) {
	if condition {

		ToError(errorCode, statusCode, errorMsg)
	}
}
