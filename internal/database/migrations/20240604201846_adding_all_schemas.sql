-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS students (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS society_coordinators (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS societies (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  president_id INTEGER UNIQUE REFERENCES students (id),
  active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS student_societies (
  student_id INTEGER REFERENCES students (id) ON DELETE CASCADE,
  society_id INTEGER REFERENCES societies (id) ON DELETE CASCADE,
  PRIMARY KEY (student_id, society_id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS student_societies;

DROP TABLE IF EXISTS societies;

DROP TABLE IF EXISTS society_coordinators;

DROP TABLE IF EXISTS students;

-- +goose StatementEnd
