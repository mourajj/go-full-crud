CREATE TABLE IF NOT EXISTS members (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL,
    role VARCHAR(255),
    duration INT,
    tags VARCHAR(255)
);
