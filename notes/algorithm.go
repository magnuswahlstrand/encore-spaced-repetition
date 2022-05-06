package notes

import (
	"log"
	"math"
	"time"
)

const (
	learningInterval1Minute   = 1 * time.Minute
	learningInterval10Minutes = 10 * time.Minute
)

var easiness = map[int32]float64{
	AnswerAgain: -0.8,
	AnswerHard:  -0.14,
	AnswerGood:  0,
	AnswerEasy:  0.1,
}

func Review(n Note, answer int32) Note {
	// Learning mode
	if n.IsLearning {
		handleLearningPhase(&n, answer)
		return n
	}

	// Review mode
	// https://docs.ankiweb.net/studying.html#review-cards
	return handleReviewPhase(n, answer)
}

func handleReviewPhase(n Note, answer int32) Note {
	n.EasinessFactor += easiness[answer]
	n.EasinessFactor = math.Max(1.3, n.EasinessFactor)

	switch answer {
	case AnswerAgain:
		// Reset
		n.IsLearning = true
		n.RepetitionNumber = 0
		n.Interval = 0
		n.NextReview = n.NextReview.Add(learningInterval1Minute)
		return n
	case AnswerHard:
		n.Interval = int32(float64(n.Interval) * 1.2)
	case AnswerGood:
		n.Interval = int32(float64(n.Interval) * n.EasinessFactor)
	case AnswerEasy:
		// https://docs.ankiweb.net/deck-options.html#easy-bonus
		n.Interval = int32(float64(n.Interval) * n.EasinessFactor * 1.3)
	}
	n.NextReview = n.NextReview.Add(time.Duration(n.Interval) * time.Minute)
	return n
}

func gotoReviewPhase(n *Note) {
	n.IsLearning = false
	n.RepetitionNumber = 0
	n.Interval = 60 * 24
	n.NextReview = n.NextReview.Add(time.Duration(n.Interval) * time.Minute)
}

func handleLearningPhase(n *Note, answer int32) {
	switch answer {
	case AnswerAgain:
		n.RepetitionNumber = 0
		n.NextReview = n.NextReview.Add(learningInterval1Minute)
	case AnswerHard:
		n.NextReview = n.NextReview.Add(learningInterval1Minute)
	case AnswerGood:
		n.RepetitionNumber++

		switch n.RepetitionNumber {
		case 1:
			n.NextReview = n.NextReview.Add(learningInterval1Minute)
		case 2:
			n.NextReview = n.NextReview.Add(learningInterval10Minutes)
		default:
			gotoReviewPhase(n)
		}
	case AnswerEasy:
		gotoReviewPhase(n)
	default:
		log.Fatalln("shouldn't happen", answer)
	}
}
