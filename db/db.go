package db

import "github.com/jmoiron/sqlx"

//Config for making functions. Stores connection to db.
type Config struct {
	DB *sqlx.DB
}

//New is a factory for Config
func New(connection *sqlx.DB) *Config {
	return &Config{
		DB: connection,
	}
}

//Init creates database and tables if needed.
func (c *Config) Init() error {
	const addExtension = `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	if _, err := c.DB.Exec(addExtension); err != nil {
		return err
	}
	return insertCartTable(c.DB)
}

func insertCartTable(db *sqlx.DB) error {
	const createTable = `CREATE TABLE IF NOT EXISTS public.cart
	(
		uuid uuid NOT NULL,
		item_uuid uuid NOT NULL,
		cart_uuid uuid NOT NULL,
		CONSTRAINT cart_pkey PRIMARY KEY (uuid)
	) WITH (OIDS=FALSE);`
	if _, err := db.Exec(createTable); err != nil {
		return err
	}
	return insertClientTable(db)
}

func insertClientTable(db *sqlx.DB) error {
	const createTable = `CREATE TABLE IF NOT EXISTS public.client
	(
		uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
		first_name text COLLATE pg_catalog."default" NOT NULL,
		last_name text COLLATE pg_catalog."default" NOT NULL,
		email text COLLATE pg_catalog."default" NOT NULL,
		phone text COLLATE pg_catalog."default" NOT NULL,
		address text COLLATE pg_catalog."default" NOT NULL,
		CONSTRAINT "Cart_pkey" PRIMARY KEY (uuid),
		CONSTRAINT email_check CHECK (email ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'::text),
		CONSTRAINT "Check_text_empty" CHECK (length(first_name) > 0 AND length(last_name) > 0 AND length(email) > 0) NOT VALID
	) WITH (OIDS=FALSE);`
	if _, err := db.Exec(createTable); err != nil {
		return err
	}
	return insertItemTable(db)
}

func insertItemTable(db *sqlx.DB) error {
	const createTable = `CREATE TABLE IF NOT EXISTS public.item
	(
		uuid uuid NOT NULL,
		item_name character varying COLLATE pg_catalog."default" NOT NULL,
		display_name character varying COLLATE pg_catalog."default",
		price numeric,
		currency integer,
		available integer,
		description text COLLATE pg_catalog."default",
		image_path character varying COLLATE pg_catalog."default",
		item_type_id integer,
		CONSTRAINT item_pkey PRIMARY KEY (uuid),
		CONSTRAINT check_item_name CHECK (length(item_name::text) > 0)
	) WITH (OIDS=FALSE);`
	if _, err := db.Exec(createTable); err != nil {
		return err
	}
	return insertOderTable(db)
}

func insertOderTable(db *sqlx.DB) error {
	const createTable = `CREATE TABLE IF NOT EXISTS public.orders
	(
		uuid uuid NOT NULL,
		client_uuid uuid NOT NULL,
		cart_uuid uuid NOT NULL,
		date timestamp(6) with time zone DEFAULT now(),
		is_payed boolean,
		status integer,
		status_date timestamp(6) with time zone DEFAULT now(),
		CONSTRAINT order_pkey PRIMARY KEY (uuid)
	) WITH (OIDS=FALSE);`
	if _, err := db.Exec(createTable); err != nil {
		return err
	}
	return insertItemTypeTable(db)
}

func insertItemTypeTable(db *sqlx.DB) error {
	const createTable = `CREATE TABLE IF NOT EXISTS public.item_type
	(
		id serial,
		display character varying(255) COLLATE pg_catalog."default" NOT NULL,
		CONSTRAINT item_type_pkey PRIMARY KEY (id),
		CONSTRAINT check_display_not_empty CHECK (length(display::text) > 0) NOT VALID
	) WITH (OIDS=FALSE);`
	if _, err := db.Exec(createTable); err != nil {
		return err
	}
	return nil
}
