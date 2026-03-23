package data

import (
	"fmt"
	"strconv"
)

// Runtime represents a movie duration in minutes.
type Runtime int32

// MarshalJSON formats runtime values as JSON strings like "102 mins".
func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	// Quote the string so the returned bytes are valid JSON.
	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}
