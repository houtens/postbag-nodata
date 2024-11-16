-- up migration for pw_token and pw_token_expiry
alter table users add column pw_token varchar;
alter table users add column pw_token_expiry timestamp;
