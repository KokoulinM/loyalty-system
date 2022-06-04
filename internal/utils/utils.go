package utils

import "time"

func CalculateAdditionTime(count int) time.Duration {
	return time.Second * time.Duration(30*(count+1))
}
