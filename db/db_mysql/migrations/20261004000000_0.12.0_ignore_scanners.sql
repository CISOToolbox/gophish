
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `campaigns` ADD COLUMN `ignore_scanners` BOOLEAN NOT NULL DEFAULT 0;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `campaigns` DROP COLUMN `ignore_scanners`;
