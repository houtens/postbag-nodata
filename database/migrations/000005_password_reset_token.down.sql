-- down migrations for pw_token and pw_token_expiry
alter table users drop column pw_token;
alter table users drop column pw_token_expiry;
