CREATE TABLE IF NOT EXISTS attachments (
    id           BIGINT       NOT NULL AUTO_INCREMENT,
    ticket_id    BIGINT       NOT NULL,
    uploaded_by  BIGINT       NOT NULL,
    file_name    VARCHAR(255) NOT NULL,
    file_key     VARCHAR(512) NOT NULL,
    file_size    INT          NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    created_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_attachments_ticket_id (ticket_id),
    CONSTRAINT fk_attachments_ticket_id   FOREIGN KEY (ticket_id)   REFERENCES tickets (id) ON DELETE CASCADE,
    CONSTRAINT fk_attachments_uploaded_by FOREIGN KEY (uploaded_by) REFERENCES users   (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
