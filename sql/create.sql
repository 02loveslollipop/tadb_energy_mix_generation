CREATE SCHEMA core;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE core.generator;
DROP TABLE core.type;
DROP TABLE core.production;

CREATE TABLE core.type(
    id  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(20) UNIQUE NOT NULL,
    description varchar(80) NOT NULL,
    isRenuevable bool NOT NULL
);

CREATE TABLE core.generator(
    id  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type UUID NOT NULL,
    capacity FLOAT NOT NULL,
    CONSTRAINT fk_type
        FOREIGN KEY (type)
        REFERENCES core.type(id)
        ON DELETE CASCADE
);

CREATE TABLE core.production(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    generator_id UUID NOT NULL,
    date DATE NOT NULL,
    production_mw DECIMAL NOT NULL,
    CONSTRAINT fk_generator
        FOREIGN KEY (generator_id)
        REFERENCES core.generator(id)
        ON DELETE CASCADE,
    CONSTRAINT uk_generator_date
        UNIQUE(generator_id,date)
);
