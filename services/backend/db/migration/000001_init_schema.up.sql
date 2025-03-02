CREATE TABLE users (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    username varchar NOT NULL UNIQUE,
    source varchar NOT NULL
);

CREATE TABLE projects (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    name varchar NOT NULL,
    description varchar NOT NULL,
    user_id uuid NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE tasks (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    project_id uuid NOT NULL,
    name varchar NOT NULL,
    description varchar NOT NULL,
    CONSTRAINT fk_project_id FOREIGN KEY (project_id) REFERENCES projects (id)
);

CREATE TABLE tasks_users (
    task_id uuid NOT NULL,
    user_id uuid NOT NULL,
    started_at timestamptz NOT NULL,
    finished_at timestamptz,
    PRIMARY KEY(task_id, user_id),
    CONSTRAINT fk_task_id FOREIGN KEY (task_id) REFERENCES tasks (id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);
