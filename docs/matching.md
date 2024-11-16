# Fuzzy string matching in postgresql

https://www.freecodecamp.org/news/fuzzy-string-matching-with-postgresql/
https://www.postgresql.org/docs/current/fuzzystrmatch.html
https://www.crunchydata.com/blog/fuzzy-name-matching-in-postgresql

# Trigram matching

```
create extension pg_trgm;
select first_name, last_name, alt_name from users where 'steve oxford' % concat_ws(' ', first_name, last_name, alt_name);
```

# Soundex

```
create extension fuzzystrmatch;
select first_name, last_name, alt_name from users where soundex(concat_ws(' ', first_name, last_name, alt_name)) = soundex('paul tompson');
```
