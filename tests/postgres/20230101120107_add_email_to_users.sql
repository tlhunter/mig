--BEGIN MIGRATION UP--
ALTER TABLE users ADD COLUMN email varchar(64) UNIQUE;
UPDATE "users" SET email = 'tlhunter@example.com' WHERE username = 'tlhunter';
--END MIGRATION UP--
--BEGIN MIGRATION DOWN--
ALTER TABLE users DROP COLUMN email;
--END MIGRATION DOWN--