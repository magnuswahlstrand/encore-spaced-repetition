package notes

import (
	"encore.app/notes/db"
	"time"
)

type Note db.Note

const (
	StartEasinessFactor   = 2.5
	StartRepetitionNumber = 0
	StartInterval         = 1
)

func New(front, back string) Note {
	return Note{
		EasinessFactor:   StartEasinessFactor,
		RepetitionNumber: StartRepetitionNumber,
		Interval:         StartInterval,
		NoteFront:        front,
		NoteBack:         back,
		NextReview:       time.Now().UTC(),
		IsLearning:       true,
	}
}

func (n Note) ToCreateParams() db.CreateParams {
	return db.CreateParams{
		NoteFront:        n.NoteFront,
		NoteBack:         n.NoteBack,
		EasinessFactor:   n.EasinessFactor,
		RepetitionNumber: n.RepetitionNumber,
		Interval:         n.Interval,
		NextReview:       n.NextReview,
	}
}

func (n Note) UpdateParams() db.UpdateReviewStatusParams {
	return db.UpdateReviewStatusParams{
		ID:               n.ID,
		EasinessFactor:   n.EasinessFactor,
		RepetitionNumber: n.RepetitionNumber,
		Interval:         n.Interval,
		NextReview:       n.NextReview,
	}
}

func (n Note) toResponse() *ResponseNote {
	return &ResponseNote{
		ID:         n.ID,
		Front:      n.NoteFront,
		Back:       n.NoteBack,
		NextReview: n.NextReview,
	}
}
