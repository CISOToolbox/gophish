
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `scanner_settings` (
    `id`             INTEGER PRIMARY KEY,
    `window_seconds` INTEGER,
    `min_chrome`     INTEGER,
    `min_firefox`    INTEGER,
    `min_safari`     INTEGER,
    `user_agents`    MEDIUMTEXT
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `scanner_settings`;
