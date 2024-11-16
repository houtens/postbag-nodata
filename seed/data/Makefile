default: resetdb importdb fix_sql export_tables iconv
	# finished

resetdb: drop_absp create_absp

drop_absp:
	-mysql -e 'drop database absp;'

create_absp:
	-mysql -e 'create database absp;'

importdb:
	sed -i '' 's/absp_database_org/absp/g' sql/absp-database-org.sql
	cat sql/absp-database-org.sql | mysql absp

fix_sql: rename_tables remove_nonprintables

rename_tables:
	mysql absp < sql/1-rename-tables-drop-absp-prefix.sql

remove_nonprintables:
	mysql absp < sql/2-remove-windows-newlines.sql

export_tables:
	mkdir -p export
	rm -f export/*.csv
	mysql absp < sql/3-export-all-tables.sql


iconv: iconv_tournaments_csv iconv_invoices_csv

iconv_tournaments_csv:
	-iconv -c -t utf-8 export/tournaments.csv > export/tmp.csv
	mv export/tmp.csv export/tournaments.csv

iconv_invoices_csv:
	-iconv -c -t utf-8 export/invoices.csv > export/tmp.csv
	mv export/tmp.csv export/invoices.csv

.PHONY: iconv_tournaments_csv
