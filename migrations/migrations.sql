CREATE TABLE IF NOT EXISTS users (
	id serial NOT NULL CONSTRAINT users_pk PRIMARY KEY,
	name VARCHAR (100) NOT NULL,
	rank VARCHAR (50)
);

CREATE UNIQUE INDEX IF NOT EXISTS users_id_uindex on users (id);

CREATE TABLE IF NOT EXISTS cars (
	id serial NOT NULL CONSTRAINT cars_pk PRIMARY KEY,
	user_id integer CONSTRAINT cars_users_id_fk REFERENCES users,
	colour VARCHAR (50),
	brand VARCHAR (50),
	license VARCHAR (50)
);

CREATE UNIQUE INDEX IF NOT EXISTS cars_id_uindex on cars (id);

INSERT INTO users (name, rank) VALUES 
	('Lana', 'CEO'),
	('Anna', 'QA Engineer'),
	('Jimmy', 'Senior Engineer'),
	('Rob', 'Sales'),
	('Vanessa', 'Teamlead');

INSERT INTO cars (user_id, colour, brand, license) VALUES 
	(1, 'blue', 'BMW', 'L666YY'),
	(1, 'white', 'Mercedes', 'L888YY'),
	(3, 'white', 'Nissan', 'M327HG'),
	(4, 'red', 'Ford', 'Q184XD'),
	(5, 'yellow', 'Lada', 'P843FM'),
	(5, 'red', 'Niva', 'C325QB');
