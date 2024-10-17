CREATE TABLE IF NOT EXISTS borrows (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    book_id INT REFERENCES books(id),
    borrowed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    returned_at TIMESTAMP
);

