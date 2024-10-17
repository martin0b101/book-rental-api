CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY, 
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (first_name, last_name) VALUES
('John', 'Doe'),
('Jane', 'Smith'),
('Michael', 'Johnson');
