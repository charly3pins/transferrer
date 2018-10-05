insert into "user"(name, email) values ('carles','carles@email.com');
insert into "user"(name, email) values ('john','john@email.com');

insert into account (number, balance, owner, currency) values ('XXXYYYZZZ',1234.456,'carles@email.com', 'EU');
insert into account (number, balance, owner, currency) values ('AAABBBCCC',778899,'john@email.com', 'EU');

insert into movement (origin, destination, amount) values ('XXXYYYZZZ','AAABBBCCC',100);
