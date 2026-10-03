
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- New read-only permission: view every user's campaigns and results.
INSERT INTO `permissions` (`slug`, `name`, `description`)
VALUES ('view_all_campaigns', 'View All Campaigns', 'View campaigns and results owned by any user (read-only)');

-- New read-only role that can oversee all campaigns without modifying anything.
INSERT INTO `roles` (`slug`, `name`, `description`)
VALUES ('auditor', 'Auditor', 'Read-only access to view all users'' campaigns and results');

-- The auditor role can view standard objects...
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id FROM roles AS r, `permissions` AS p
WHERE r.slug='auditor' AND p.slug='view_objects';

-- ...and can view every user's campaigns.
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id FROM roles AS r, `permissions` AS p
WHERE r.slug='auditor' AND p.slug='view_all_campaigns';

-- Admins also get cross-user campaign visibility.
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id FROM roles AS r, `permissions` AS p
WHERE r.slug='admin' AND p.slug='view_all_campaigns';

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DELETE FROM `role_permissions` WHERE `permission_id`=(SELECT `id` FROM `permissions` WHERE `slug`='view_all_campaigns');
DELETE rp FROM `role_permissions` rp INNER JOIN `roles` r ON rp.role_id=r.id WHERE r.slug='auditor';
DELETE FROM `permissions` WHERE `slug`='view_all_campaigns';
DELETE FROM `roles` WHERE `slug`='auditor';
