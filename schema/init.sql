CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50),
    email VARCHAR(100) UNIQUE NOT NULL,
    -- role VARCHAR(100)
    -- roles VARCHAR(100)[], -- select * where (unnest(users.roles) intersect unnest(projects.roles)) is not null
    password VARCHAR(255) NOT NULL
);

CREATE TABLE projects (
    project_id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT
    -- author_id INT NOT NULL,
    -- roles VARCHAR(100)[], 
);

CREATE TABLE user_projects (
    user_id INT,
    project_id INT,
    user_role SMALLINT NOT NULL,
    PRIMARY KEY (user_id, project_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (project_id) REFERENCES projects(project_id)
);

CREATE TABLE tasks (
    task_id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    status SMALLINT NOT NULL,
    deadline TIMESTAMPTZ,
    project_id INT,
    assigned_user_id INT,
    FOREIGN KEY (project_id) REFERENCES projects(project_id),
    FOREIGN KEY (assigned_user_id) REFERENCES users(user_id)
);

-- CREATE TABLE task_assignments (
--     task_assignment_id SERIAL PRIMARY KEY,
--     task_id INT,
--     user_id INT,
--     assignment_date TIMESTAMPTZ,
--     completion_date TIMESTAMPTZ,
--     FOREIGN KEY (task_id) REFERENCES tasks(task_id),
--     FOREIGN KEY (user_id) REFERENCES users(user_id)
-- );
