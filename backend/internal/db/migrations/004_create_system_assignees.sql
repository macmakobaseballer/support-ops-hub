-- +goose Up
CREATE TABLE IF NOT EXISTS system_assignees (
    system_id BIGINT NOT NULL,
    user_id   BIGINT NOT NULL,
    PRIMARY KEY (system_id, user_id),
    INDEX idx_system_assignees_user_id (user_id),
    CONSTRAINT fk_system_assignees_system_id FOREIGN KEY (system_id) REFERENCES systems (id) ON DELETE CASCADE,
    CONSTRAINT fk_system_assignees_user_id   FOREIGN KEY (user_id)   REFERENCES users   (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS system_assignees;
