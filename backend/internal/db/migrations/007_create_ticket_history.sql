-- +goose Up
CREATE TABLE IF NOT EXISTS ticket_history (
    id         BIGINT       NOT NULL AUTO_INCREMENT,
    ticket_id  BIGINT       NOT NULL,
    changed_by BIGINT       NOT NULL,
    field_name VARCHAR(100) NOT NULL,
    old_value  TEXT         NULL,
    new_value  TEXT         NULL,
    changed_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_ticket_history_ticket_id (ticket_id),
    CONSTRAINT fk_ticket_history_ticket_id  FOREIGN KEY (ticket_id)  REFERENCES tickets (id) ON DELETE CASCADE,
    CONSTRAINT fk_ticket_history_changed_by FOREIGN KEY (changed_by) REFERENCES users   (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS ticket_history;
