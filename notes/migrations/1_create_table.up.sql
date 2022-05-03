CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE notes
(
    id                UUID DEFAULT uuid_generate_v4() NOT NULL,

    easiness_factor   FLOAT                           NOT NULL,
    repetition_number INT                             NOT NULL,
    interval          INT                             NOT NULL,

    note_front        TEXT                            NOT NULL,
    note_back         TEXT                            NOT NULL,
    next_review       TIMESTAMP                       NOT NULL
);
