CREATE OR REPLACE PROCEDURE core.insert_type(
    p_name text,
    p_description text,
    p_is_renuevable boolean
)
LANGUAGE plpgsql
as $$
    declare same_count integer;

    begin
        if p_name is null or
            p_description is null or
            length(p_name) = 0 or
            length(p_description) = 0 then
                raise exception 'Receive null or empty values';
        end if;

        select count(t.id) into same_count
        from core.type as t
        where lower(t.name) = lower(p_name);

        if same_count > 0 then
            raise exception 'A register with the same name already exist, please use another name or modify the existing register';
        end if;

        insert into core.type(name, description, isrenuevable)
        values(p_name, p_description, p_is_renuevable);

    end;
$$;

-- =====================================================
-- TYPE PROCEDURES
-- =====================================================

-- Update type procedure
CREATE OR REPLACE PROCEDURE core.update_type(
    p_type_id uuid,
    p_name text,
    p_description text,
    p_is_renuevable boolean
)
LANGUAGE plpgsql
as $$
    declare 
        same_count integer;
        type_exists integer;

    begin
        -- Check if type exists
        select count(t.id) into type_exists
        from core.type as t
        where t.id = p_type_id;

        if type_exists = 0 then
            raise exception 'Type with ID % does not exist', p_type_id;
        end if;

        -- Validate input
        if p_name is null or
            p_description is null or
            length(p_name) = 0 or
            length(p_description) = 0 then
                raise exception 'Receive null or empty values';
        end if;

        -- Check for duplicate name (excluding current record)
        select count(t.id) into same_count
        from core.type as t
        where lower(t.name) = lower(p_name) and t.id != p_type_id;

        if same_count > 0 then
            raise exception 'A register with the same name already exist, please use another name or modify the existing register';
        end if;

        update core.type 
        set name = p_name, description = p_description, isrenuevable = p_is_renuevable
        where id = p_type_id;

    end;
$$;

-- Delete type procedure
CREATE OR REPLACE PROCEDURE core.delete_type(
    p_type_id uuid
)
LANGUAGE plpgsql
as $$
    declare type_exists integer;

    begin
        -- Check if type exists
        select count(t.id) into type_exists
        from core.type as t
        where t.id = p_type_id;

        if type_exists = 0 then
            raise exception 'Type with ID % does not exist', p_type_id;
        end if;

        delete from core.type where id = p_type_id;

    end;
$$;

-- =====================================================
-- GENERATOR PROCEDURES
-- =====================================================

-- Insert generator procedure
CREATE OR REPLACE PROCEDURE core.insert_generator(
    p_generator_type uuid,
    p_generator_capacity float
)
LANGUAGE plpgsql
as $$
    declare type_exists integer;

    begin
        -- Validate input
        if p_generator_type is null or p_generator_capacity is null then
            raise exception 'Receive null values';
        end if;

        if p_generator_capacity <= 0 then
            raise exception 'Capacity must be greater than 0';
        end if;

        -- Check if type exists
        select count(t.id) into type_exists
        from core.type as t
        where t.id = p_generator_type;

        if type_exists = 0 then
            raise exception 'Type with ID % does not exist', p_generator_type;
        end if;

        insert into core.generator(type, capacity)
        values(p_generator_type, p_generator_capacity);

    end;
$$;

-- Update generator procedure
CREATE OR REPLACE PROCEDURE core.update_generator(
    p_generator_id uuid,
    p_generator_type uuid,
    p_generator_capacity float
)
LANGUAGE plpgsql
as $$
    declare 
        generator_exists integer;
        type_exists integer;

    begin
        -- Check if generator exists
        select count(g.id) into generator_exists
        from core.generator as g
        where g.id = p_generator_id;

        if generator_exists = 0 then
            raise exception 'Generator with ID % does not exist', p_generator_id;
        end if;

        -- Validate input
        if p_generator_type is null or p_generator_capacity is null then
            raise exception 'Receive null values';
        end if;

        if p_generator_capacity <= 0 then
            raise exception 'Capacity must be greater than 0';
        end if;

        -- Check if type exists
        select count(t.id) into type_exists
        from core.type as t
        where t.id = p_generator_type;

        if type_exists = 0 then
            raise exception 'Type with ID % does not exist', p_generator_type;
        end if;

        update core.generator 
        set type = p_generator_type, capacity = p_generator_capacity
        where id = p_generator_id;

    end;
$$;

-- Delete generator procedure
CREATE OR REPLACE PROCEDURE core.delete_generator(
    p_generator_id uuid
)
LANGUAGE plpgsql
as $$
    declare generator_exists integer;

    begin
        -- Check if generator exists
        select count(g.id) into generator_exists
        from core.generator as g
        where g.id = p_generator_id;

        if generator_exists = 0 then
            raise exception 'Generator with ID % does not exist', p_generator_id;
        end if;

        delete from core.generator where id = p_generator_id;

    end;
$$;

-- =====================================================
-- PRODUCTION PROCEDURES
-- =====================================================

-- Insert production procedure
CREATE OR REPLACE PROCEDURE core.insert_production(
    p_generator_id uuid,
    p_date date,
    p_production_mw decimal
)
LANGUAGE plpgsql
as $$
    declare 
        generator_exists integer;
        production_exists integer;

    begin
        -- Validate input
        if p_generator_id is null or p_date is null or p_production_mw is null then
            raise exception 'Receive null values';
        end if;

        if p_production_mw < 0 then
            raise exception 'Production must be greater than or equal to 0';
        end if;

        -- Check if generator exists
        select count(id) into generator_exists
        from core.generator as g
        where g.id = p_generator_id;

        if generator_exists = 0 then
            raise exception 'Generator with ID % does not exist', p_generator_id;
        end if;

        -- Check if production record already exists for this generator and date
        select count(id) into production_exists
        from core.production as p
        where p.generator_id = p_generator_id and p.date = p_date;

        if production_exists > 0 then
            raise exception 'Production record already exists for generator % on date %', p_generator_id, p_date;
        end if;

        insert into core.production(generator_id, date, production_mw)
        values(p_generator_id, p_date, p_production_mw);

    end;
$$;

-- Update production procedure
CREATE OR REPLACE PROCEDURE core.update_production(
    p_production_id uuid,
    p_generator_id uuid,
    p_date date,
    p_production_mw decimal
)
LANGUAGE plpgsql
as $$
    declare 
        production_exists integer;
        generator_exists integer;
        duplicate_exists integer;

    begin
        -- Check if production record exists
        select count(p.id) into production_exists
        from core.production as p
        where p.id = p_production_id;

        if production_exists = 0 then
            raise exception 'Production record with ID % does not exist', p_production_id;
        end if;

        -- Validate input
        if p_generator_id is null or p_date is null or p_production_mw is null then
            raise exception 'Receive null values';
        end if;

        if p_production_mw < 0 then
            raise exception 'Production must be greater than or equal to 0';
        end if;

        -- Check if generator exists
        select count(g.id) into generator_exists
        from core.generator as g
        where g.id = p_generator_id;

        if generator_exists = 0 then
            raise exception 'Generator with ID % does not exist', p_generator_id;
        end if;

        -- Check for duplicate generator-date combination (excluding current record)
        select count(p.id) into duplicate_exists
        from core.production as p
        where p.generator_id = p_generator_id and p.date = p_date and p.id != p_production_id;

        if duplicate_exists > 0 then
            raise exception 'Production record already exists for generator % on date %', p_generator_id, p_date;
        end if;

        update core.production 
        set generator_id = p_generator_id, date = p_date, production_mw = p_production_mw
        where id = p_production_id;

    end;
$$;

-- Delete production procedure
CREATE OR REPLACE PROCEDURE core.delete_production(
    p_production_id uuid
)
LANGUAGE plpgsql
as $$
    declare production_exists integer;

    begin
        -- Check if production record exists
        select count(p.id) into production_exists
        from core.production as p
        where p.id = p_production_id;

        if production_exists = 0 then
            raise exception 'Production record with ID % does not exist', p_production_id;
        end if;

        delete from core.production where id = p_production_id;

    end;
$$;
)
LANGUAGE plpgsql
as $$
    declare production_exists integer;

    begin
        -- Check if production record exists
        select count(id) into production_exists
        from core.production as p
        where p.id = production_id;

        if production_exists = 0 then
            raise exception 'Production record with ID % does not exist', production_id;
        end if;

        delete from core.production where id = production_id;

    end;
$$;