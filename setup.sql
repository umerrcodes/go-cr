-- Create tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert some sample data
INSERT INTO tasks (title, description, completed) VALUES 
    ('Learn Go', 'Study Go programming language', FALSE),
    ('Build API', 'Create a REST API with Go', TRUE),
    ('Add Database', 'Integrate PostgreSQL with Go API', FALSE);

-- Display the created table
SELECT * FROM tasks; 