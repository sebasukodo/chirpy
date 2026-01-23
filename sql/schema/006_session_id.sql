-- +goose UP
ALTER TABLE refresh_tokens RENAME TO session_ids;
ALTER TABLE session_ids RENAME COLUMN token TO id;

-- +goose Down
ALTER TABLE session_ids RENAME COLUMN id TO token;
ALTER TABLE session_ids RENAME TO refresh_tokens;