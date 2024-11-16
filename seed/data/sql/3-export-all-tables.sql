-- clubs
SELECT
"country","county","club_name","type","contact","contact_name","email","phone","website","updated","id"
UNION ALL
SELECT
* FROM
clubs
INTO OUTFILE "/Users/shouten/src/postbag/seed/data/export/clubs.csv"
FIELDS
TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\r\n';

-- countries
SELECT
"id","name","filename","tier"
UNION ALL
SELECT
* FROM
countries
INTO OUTFILE "/Users/shouten/src/postbag/seed/data/export/countries.csv"
FIELDS
TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\r\n';

-- invoices
-- UPDATE invoices set comment = replace() 
SELECT
"id","date","tournament","players","games","supp_mems","multiday","penalty","amount","type","comment","paid","locked"
UNION ALL
SELECT
* FROM
invoices
INTO OUTFILE "/Users/shouten/src/postbag/seed/data/export/invoices.csv"
FIELDS
TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\r\n';

-- log_entries
SELECT
"timestamp","user","status","action","data","ip","id"
UNION ALL
SELECT
* FROM
log_entries
INTO OUTFILE "/Users/shouten/src/postbag/seed/data/export/log_entries.csv"
FIELDS
TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\r\n';

-- member_years
SELECT
"id","year","member"
UNION ALL
SELECT
* FROM
member_years
INTO OUTFILE "/Users/shouten/src/postbag/seed/data/export/member_years.csv"
FIELDS
TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\r\n';

-- members
SELECT
"id","absp_number","access_level","first_name","last_name","status","address1","address2","address3","address4",
"postcode","phone","email","club","club_id","country","year","paid_when","paid_how","last_payment","mobile",
"notes","life","ob_post","ob_pdf","dob","last_played","rating","total_games","updated","password_hash","fb_id"
UNION ALL
SELECT
* FROM
members
INTO OUTFILE "/Users/shouten/src/postbag/seed/data/export/members.csv"
FIELDS
TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\r\n';

-- nsc_qualifiers
SELECT
"id","year","player","date","tournament","qualified","registered"
UNION ALL
SELECT
* FROM
nsc_qualifiers
INTO OUTFILE "/Users/shouten/src/postbag/seed/data/export/nsc_qualifiers.csv"
FIELDS
TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\r\n';

-- payments
SELECT
"id","date","member","amount","method","subscription_year"
UNION ALL
SELECT
* FROM
payments
INTO OUTFILE "/Users/shouten/src/postbag/seed/data/export/payments.csv"
FIELDS
TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\r\n';

-- ratings
SELECT
"id","player","tournament","division","team","games","start_rating","points","date","locked","verified"
UNION ALL
SELECT
* FROM
ratings
INTO OUTFILE "/Users/shouten/src/postbag/seed/data/export/ratings.csv"
FIELDS
TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\r\n';

-- results
SELECT
"id","player1","player2","score1","score2","spread","type","tournament","round","locked"
UNION ALL
SELECT
* FROM
results
INTO OUTFILE "/Users/shouten/src/postbag/seed/data/export/results.csv"
FIELDS
TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\r\n';

-- tournaments
SELECT
"id","short_name","title","date","entries","division","rounds","results","team_event","organiser","td","co","creator","updated","comment","invoice","locked","datafile"
UNION ALL
SELECT
* FROM
tournaments
INTO OUTFILE "/Users/shouten/src/postbag/seed/data/export/tournaments.csv"
FIELDS
TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\r\n';

