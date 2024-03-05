-- +goose Up
ALTER TABLE users ADD COLUMN apikey VARCHAR(64) UNIQUE;
UPDATE users SET apikey = (encode(sha256(random()::text::bytea), 'hex'));
ALTER TABLE users ALTER COLUMN apikey SET NOT NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN apikey;
```
