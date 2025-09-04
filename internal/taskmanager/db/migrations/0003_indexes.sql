CREATE INDEX IF NOT EXISTS tasks_user_id_idx ON tasks (user_id);
CREATE INDEX IF NOT EXISTS tasks_user_id_status_idx ON tasks (user_id, status);

CREATE INDEX IF NOT EXISTS tasks_user_id_active_idx
    ON tasks (user_id)
    WHERE status != 'done';
