package drmtime

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

// CustomTime hold time information.
type CustomTime struct {
	Time time.Duration `validate:"required"`
}

// UnmarshalJSON for time in string to time.Duration.
func (ct *CustomTime) UnmarshalJSON(data []byte) (err error) {
	var tmp string

	if err = json.Unmarshal(data, &tmp); err != nil {
		return errors.Wrap(err, "Custom Time Unmarshal")
	}

	dur, err := time.ParseDuration(tmp)
	if err != nil {
		return errors.Wrap(err, "Custom Time Unmarshal")
	}

	ct.Time = dur

	return errors.Wrap(err, "Custom Time Unmarshal")
}
