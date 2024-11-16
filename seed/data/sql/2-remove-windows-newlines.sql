update members set address1 = replace(address1, '"', '');
update members set club = replace(club, '"', '');
update members set last_name = replace(last_name, '\r', '');
update members set last_name = replace(last_name, '\n', '');
update clubs set phone = replace(phone, char(0xA0), '');
update tournaments set title = replace(title, '\r\n', '');
update tournaments set title = replace(title, '\n', '');
update tournaments set title = replace(title, '"', '');
update tournaments set `comment` = replace(`comment`, '\r\n', '');
update tournaments set `comment` = replace(`comment`, '\n', '');
update invoices set `comment` = replace(`comment`, '"', '');
update invoices set `comment` = replace(`comment`, '\r\n', '');
update invoices set `comment` = replace(`comment`, '\n', '');

-- delete rating for non-existent tournament
delete from ratings where tournament = 16974;
