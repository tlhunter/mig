--BEGIN MIGRATION UP--
ALTER TABLE users ADD COLUMN email varchar(64);
CREATE UNIQUE INDEX users_email ON users(email);
UPDATE "users" SET email = 'tlhunter@example.com' WHERE username = 'tlhunter';
--END MIGRATION UP--
--BEGIN MIGRATION DOWN--
DROP INDEX users_email;
ALTER TABLE users DROP COLUMN email;
--END MIGRATION DOWN--