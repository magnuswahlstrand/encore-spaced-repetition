package main

import (
	"encore.app/notes"
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

func main() {
	tcs := []struct {
		title  string
		answer []int32
	}{
		{"Answer again", []int32{notes.AnswerAgain, notes.AnswerAgain, notes.AnswerAgain, notes.AnswerAgain, notes.AnswerAgain}},
		{"Answer hard", []int32{notes.AnswerHard, notes.AnswerHard, notes.AnswerHard, notes.AnswerHard, notes.AnswerHard}},
		{"Answer good", []int32{notes.AnswerGood, notes.AnswerGood, notes.AnswerGood, notes.AnswerGood, notes.AnswerGood}},
		{"Answer easy", []int32{notes.AnswerEasy, notes.AnswerEasy, notes.AnswerEasy, notes.AnswerEasy, notes.AnswerEasy}},
		{"Test lapse", []int32{notes.AnswerEasy, notes.AnswerAgain, notes.AnswerEasy, notes.AnswerAgain}},
		{"Test hard", []int32{notes.AnswerEasy, notes.AnswerHard, notes.AnswerHard, notes.AnswerHard}},
		{"Test good", []int32{notes.AnswerEasy, notes.AnswerGood, notes.AnswerGood, notes.AnswerGood}},
		{"Test rebuild ease", []int32{notes.AnswerEasy, notes.AnswerAgain, notes.AnswerEasy, notes.AnswerEasy, notes.AnswerEasy, notes.AnswerEasy}},
	}

	w := tabwriter.NewWriter(os.Stdout, 5, 5, 3, ' ', 0)
	for _, tc := range tcs {
		n := notes.New("", "")
		fmt.Fprintf(w, "\n%s\t%s\n", tc.title, n.NextReview.Format(time.Kitchen))
		for _, answer := range tc.answer {
			n = notes.Review(n, answer)
			fmt.Fprintf(w, "%s\t%0.2f\t%d\t%d\t%v\t%d\n", n.NextReview.Format("2006 Jan 02\t3:04PM"), n.EasinessFactor, n.Interval, n.RepetitionNumber, n.IsLearning, answer)
		}
	}
	w.Flush()
}
