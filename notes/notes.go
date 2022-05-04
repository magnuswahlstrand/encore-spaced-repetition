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

type ListResponse struct {
	Notes []Note `json:"notes"`
}

//encore:api public method=GET path=/note
func ListNotes(ctx context.Context) (*ListResponse, error) {
	dbNotes, err := notedb.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	notes := []Note{}
	for _, n := range dbNotes {
		notes = append(notes, fromDBnote(n))
	}

	return &ListResponse{Notes: notes}, nil
}

func fromDBnote(n db.Note) Note {
	return Note{
		ID:         n.ID,
		Front:      n.NoteFront,
		Back:       n.NoteBack,
		NextReview: n.NextReview,
	}
}

type NewNoteRequest struct {
	Front string `json:"front"`
	Back  string `json:"back"`
}

type Note struct {
	ID         uuid.UUID `json:"id"`
	Front      string    `json:"front"`
	Back       string    `json:"back"`
	NextReview time.Time `json:"next_review"`
}

//encore:api public method=POST path=/note
func NewNote(ctx context.Context, req *NewNoteRequest) (*Note, error) {
	n, err := notedb.Create(ctx, db.CreateParams{
		NoteFront:        req.Front,
		NoteBack:         req.Back,
		EasinessFactor:   1.3,
		RepetitionNumber: 1,
		Interval:         0,
		NextReview:       nextReviewTime(0),
	})
	if err != nil {
		return nil, err
	}

	note := fromDBnote(n)
	return &note, nil
}

type ReviewNoteRequest struct {
	Answer string `json:"answer"`
}

var answerMapping = map[string]int32{
	"again": 0,
	"hard":  3,
	"good":  4,
	"easy":  5,
}

//encore:api public method=POST path=/note/:id
func ReviewNote(ctx context.Context, id uuid.UUID, req *ReviewNoteRequest) (*Note, error) {
	n, err := notedb.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	q, ok := answerMapping[req.Answer]
	if !ok {
		return nil, fmt.Errorf("invalid answer: %q", req.Answer)
	}

	n.RepetitionNumber, n.EasinessFactor, n.Interval = calculateReviewStatus(q, n.RepetitionNumber, n.EasinessFactor, n.Interval)
	update := db.UpdateReviewStatusParams{
		ID:               n.ID,
		EasinessFactor:   n.EasinessFactor,
		RepetitionNumber: n.RepetitionNumber,
		Interval:         n.Interval,
		NextReview:       nextReviewTime(n.Interval),
	}
	if err := notedb.UpdateReviewStatus(ctx, update); err != nil {
		return nil, err
	}

	note := fromDBnote(n)
	return &note, nil
}
