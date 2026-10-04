
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- The per-page "Redirect to" URL is retired in favor of per-campaign
-- educational pages.
ALTER TABLE `pages` DROP COLUMN redirect_url;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `pages` ADD COLUMN redirect_url TEXT;
