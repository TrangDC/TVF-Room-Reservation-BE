CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    work_email VARCHAR(255) NOT NULL,
    oid UUID NOT NULL UNIQUE

    -- Validate work email
    CONSTRAINT check_work_email_domain CHECK (work_email ~* '^[A-Za-z0-9._%+-]+@techvify\.com\.vn$')
);