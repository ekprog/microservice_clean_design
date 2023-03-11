package conv

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func ValueOrDefault[T any](val *T, defaultVal ...T) T {

	var realV T
	if len(defaultVal) > 0 {
		realV = defaultVal[0]
	}

	if val != nil {
		realV = *val
	}
	return realV
}

func NullableTime(timeObj *time.Time) *timestamppb.Timestamp {
	if timeObj != nil {
		return timestamppb.New(*timeObj)
	}
	return nil
}
