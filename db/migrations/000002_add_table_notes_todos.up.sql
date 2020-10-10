CREATE TABLE IF NOT EXISTS notes
(
	note_id SERIAL PRIMARY KEY,
	note_name VARCHAR NOT NULL,
	user_id INTEGER NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS todos
(
    todo_id SERIAL PRIMARY KEY,
    todo_name VARCHAR NOT NULL,
    todo_done BOOLEAN NOT NULL,
    note_id INTEGER NOT NULL,
    FOREIGN KEY (note_id) REFERENCES notes(note_id)
);