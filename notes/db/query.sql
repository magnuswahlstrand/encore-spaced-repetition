-- name: Create :one
INSERT INTO notes (note_front, note_back, easiness_factor, repetition_number, interval, next_review)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: Get :one
SELECT *
FROM notes
WHERE id = $1;

-- name: UpdateReviewStatus :exec
UPDATE notes
SET easiness_factor   = $2,
    repetition_number = $3,
    interval          = $4,
    next_review       = $5
WHERE id = $1;

-- name: ListAll :many
SELECT *
FROM notes;
--

-- name: ListDue :many
SELECT *
FROM notes
WHERE next_review < NOW();
--
