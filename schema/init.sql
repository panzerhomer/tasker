CREATE TABLE tasks (
	id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE,
    description VARCHAR(600)
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
	email VARCHAR(255) UNIQUE,
	password VARCHAR(255)
);

CREATE TABLE users_tasks (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    task_id INT REFERENCES tasks(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMP,
    CONSTRAINT unique_user_task UNIQUE (user_id, task_id)
);

CREATE TABLE users_tasks_history (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    task_id INT REFERENCES tasks(id) ON DELETE CASCADE,
    action VARCHAR(30) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE OR REPLACE FUNCTION check_and_delete()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM users_tasks WHERE expired_at < NOW();
    RETURN NEW;
END;

$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER delete_expired_tasks_trigger
    AFTER INSERT ON users_tasks
    EXECUTE FUNCTION check_and_delete();