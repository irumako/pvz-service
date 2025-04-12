CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'reception_status') THEN
CREATE TYPE reception_status AS ENUM ('in_progress', 'close');
END IF;
END
$$;

DO $$
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'product_type') THEN
CREATE TYPE product_type AS ENUM ('electronics', 'clothes', 'shoes');
END IF;
END
$$;

DO $$
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
CREATE TYPE user_role AS ENUM ('employee', 'moderator');
END IF;
END
$$;

CREATE TABLE IF NOT EXISTS pvz (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    registration_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    city VARCHAR(255) NOT NULL
    );

CREATE TABLE IF NOT EXISTS reception (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dateTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    pvz_id UUID NOT NULL REFERENCES pvz(id),
    status reception_status NOT NULL
    );

CREATE TABLE IF NOT EXISTS product (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    dateTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reception_id UUID NOT NULL REFERENCES reception(id),
    type product_type NOT NULL
    );

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    role user_role NOT NULL,
    registration_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
