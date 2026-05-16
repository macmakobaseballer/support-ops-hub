-- +goose Up
CREATE TABLE IF NOT EXISTS comments (
    id         BIGINT   NOT NULL AUTO_INCREMENT,
    ticket_id  BIGINT   NOT NULL,
    author_id  BIGINT   NOT NULL,
    body       TEXT     NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_comments_ticket_id (ticket_id),
    CONSTRAINT fk_comments_ticket_id FOREIGN KEY (ticket_id) REFERENCES tickets (id) ON DELETE CASCADE,
    CONSTRAINT fk_comments_author_id FOREIGN KEY (author_id) REFERENCES users   (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS comments;
