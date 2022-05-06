package notes

import (
	"context"
	"encore.app/notes/db"
	"encore.dev/storage/sqldb"
	"encore.dev/types/uuid"
	"fmt"
	"time"
)

var notedb = db.New(sqldb.Named("notes").Stdlib())

type ResponseNote struct {
	ID         uuid.UUID `json:"id"`
	Front      string    `json:"front"`
	Back       string    `json:"back"`
	NextReview time.Time `json:"next_review"`
}

func toResponse(n db.Note) ResponseNote {
	return ResponseNote{
		ID:         n.ID,
		Front:      n.NoteFront,
		Back:       n.NoteBack,
		NextReview: n.NextReview,
	}
}

type ListResponse struct {
	Notes []ResponseNote `json:"notes"`
}

//encore:api public method=GET path=/note
func ListNotes(ctx context.Context) (*ListResponse, error) {
	dbNotes, err := notedb.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	notes := []ResponseNote{}
	for _, n := range dbNotes {
		notes = append(notes, toResponse(n))
	}

	return &ListResponse{Notes: notes}, nil
}

type NewNoteRequest struct {
	Front string `json:"front"`
	Back  string `json:"back"`
}

//encore:api public method=POST path=/note
func NewNote(ctx context.Context, req *NewNoteRequest) (*ResponseNote, error) {
	params := New(req.Front, req.Back).ToCreateParams()
	n, err := notedb.Create(ctx, params)
	if err != nil {
		return nil, err
	}

	return Note(n).toResponse(), nil
}

type ReviewNoteRequest struct {
	Answer string `json:"answer"`
}

//encore:api public method=POST path=/note/:id
func ReviewNote(ctx context.Context, id uuid.UUID, req *ReviewNoteRequest) (*ResponseNote, error) {
	var answer int32
	switch req.Answer {
	case "again":
		answer = AnswerAgain
	case "hard":
		answer = AnswerHard
	case "good":
		answer = AnswerGood
	case "easy":
		answer = AnswerEasy
	default:
		return nil, fmt.Errorf("invalid answer: %q", req.Answer)
	}

	n, err := notedb.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	//n.RepetitionNumber, n.EasinessFactor, n.Interval = calculateReviewStatus(q, n.RepetitionNumber, n.EasinessFactor, n.Interval)
	note := Review(Note(n), answer)

	//n.RepetitionNumber, n.EasinessFactor, n.Interval, n.NextReview = calculateReviewStatusV2(q, n.RepetitionNumber, n.EasinessFactor, n.Interval, n.NextReview)
	//n.NextReview = nextReviewTime(n.Interval)
	if err := notedb.UpdateReviewStatus(ctx, note.UpdateParams()); err != nil {
		return nil, err
	}

	return note.toResponse(), nil
}
