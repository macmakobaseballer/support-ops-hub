-- +goose Up
CREATE TABLE IF NOT EXISTS tickets (
    id          BIGINT       NOT NULL AUTO_INCREMENT,
    title       VARCHAR(255) NOT NULL,
    description TEXT         NULL,
    type        ENUM('question','bug','config','data')          NOT NULL,
    priority    ENUM('high','medium','low')                     NOT NULL,
    status      ENUM('new','in_progress','waiting','done')      NOT NULL DEFAULT 'new',
    customer_id BIGINT       NOT NULL,
    system_id   BIGINT       NOT NULL,
    assignee_id BIGINT       NULL,
    created_by  BIGINT       NOT NULL,
    received_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_tickets_status      (status),
    INDEX idx_tickets_priority    (priority),
    INDEX idx_tickets_type        (type),
    INDEX idx_tickets_customer_id (customer_id),
    INDEX idx_tickets_system_id   (system_id),
    INDEX idx_tickets_assignee_id (assignee_id),
    INDEX idx_tickets_received_at (received_at DESC),
    CONSTRAINT fk_tickets_customer_id FOREIGN KEY (customer_id) REFERENCES customers (id),
    CONSTRAINT fk_tickets_system_id   FOREIGN KEY (system_id)   REFERENCES systems   (id),
    CONSTRAINT fk_tickets_assignee_id FOREIGN KEY (assignee_id) REFERENCES users     (id) ON DELETE SET NULL,
    CONSTRAINT fk_tickets_created_by  FOREIGN KEY (created_by)  REFERENCES users     (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS tickets;
