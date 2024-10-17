CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    quantity INT NOT NULL
);

INSERT INTO books (title, quantity) VALUES
('The Catcher in the Rye', 5),
('To Kill a Mockingbird', 3),
('1984', 7),
('The Great Gatsby', 2),
('Moby Dick', 4),
('Pride and Prejudice', 6),
('War and Peace', 1);
