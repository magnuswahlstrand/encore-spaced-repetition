package notes

import (
	"fmt"
	"math"
	"time"
)

func calculateReviewStatus(answer, repetitionNumber int32, easinessFactor float64, interval int32) (int32, float64, int32) {
	if answer >= 3 {
		switch repetitionNumber {
		case 0:
			interval = 0
		case 1:
			interval = 6
		default:
			interval = int32(math.Round(float64(interval) * easinessFactor))
		}
		repetitionNumber++
	} else {
		// Incorrect answer
		repetitionNumber = 0
		interval = 1
		fmt.Println("incorrect")
	}
	// Update easiness factor
	easinessFactor += (0.1 - float64(5-answer)) * (0.08 + float64(5-answer)*0.02)
	easinessFactor = math.Max(1.3, easinessFactor)

	return repetitionNumber, easinessFactor, interval
}

func nextReviewTime(interval int32) time.Time {
	return time.Now().UTC().AddDate(0, 0, int(interval))
}
