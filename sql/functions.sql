-- =====================================================
-- QUERY FUNCTIONS FOR GET ENDPOINTS
-- =====================================================

-- Get all types
CREATE OR REPLACE FUNCTION core.get_all_types()
RETURNS TABLE(
    id uuid,
    name varchar(20),
    description varchar(80),
    isrenuevable boolean
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, t.name, t.description, t.isrenuevable
    FROM core.type t
    ORDER BY t.name;
END;
$$;

-- Get type by ID
CREATE OR REPLACE FUNCTION core.get_type_by_id(p_type_id uuid)
RETURNS TABLE(
    id uuid,
    name varchar(20),
    description varchar(80),
    isrenuevable boolean
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, t.name, t.description, t.isrenuevable
    FROM core.type t
    WHERE t.id = p_type_id;
END;
$$;

-- Get generators by type
CREATE OR REPLACE FUNCTION core.get_generators_by_type(p_type_id uuid)
RETURNS TABLE(
    id uuid,
    type_id uuid,
    type_name varchar(20),
    capacity float
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT g.id, g.type, t.name, g.capacity
    FROM core.generator g
    JOIN core.type t ON g.type = t.id
    WHERE g.type = p_type_id
    ORDER BY g.capacity DESC;
END;
$$;

-- Get all generators
CREATE OR REPLACE FUNCTION core.get_all_generators()
RETURNS TABLE(
    id uuid,
    type_id uuid,
    type_name varchar(20),
    type_description varchar(80),
    isrenuevable boolean,
    capacity float
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT g.id, g.type, t.name, t.description, t.isrenuevable, g.capacity
    FROM core.generator g
    JOIN core.type t ON g.type = t.id
    ORDER BY t.name, g.capacity DESC;
END;
$$;

-- Get generator by ID
CREATE OR REPLACE FUNCTION core.get_generator_by_id(p_generator_id uuid)
RETURNS TABLE(
    id uuid,
    type_id uuid,
    type_name varchar(20),
    type_description varchar(80),
    isrenuevable boolean,
    capacity float
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT g.id, g.type, t.name, t.description, t.isrenuevable, g.capacity
    FROM core.generator g
    JOIN core.type t ON g.type = t.id
    WHERE g.id = p_generator_id;
END;
$$;

-- Get productions by generator
CREATE OR REPLACE FUNCTION core.get_productions_by_generator(p_generator_id uuid)
RETURNS TABLE(
    id uuid,
    generator_id uuid,
    date date,
    production_mw decimal
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, p.generator_id, p.date, p.production_mw
    FROM core.production p
    WHERE p.generator_id = p_generator_id
    ORDER BY p.date DESC;
END;
$$;

-- Get productions by generator and date range
CREATE OR REPLACE FUNCTION core.get_productions_by_generator_date_range(
    p_generator_id uuid,
    p_start_date date,
    p_end_date date DEFAULT NULL
)
RETURNS TABLE(
    id uuid,
    generator_id uuid,
    date date,
    production_mw decimal
)
LANGUAGE plpgsql
AS $$
BEGIN
    IF p_end_date IS NULL THEN
        -- If no end date provided, get all records from start_date onwards
        RETURN QUERY
        SELECT p.id, p.generator_id, p.date, p.production_mw
        FROM core.production p
        WHERE p.generator_id = p_generator_id AND p.date >= p_start_date
        ORDER BY p.date DESC;
    ELSE
        -- Get records within the specified date range
        RETURN QUERY
        SELECT p.id, p.generator_id, p.date, p.production_mw
        FROM core.production p
        WHERE p.generator_id = p_generator_id 
        AND p.date >= p_start_date 
        AND p.date <= p_end_date
        ORDER BY p.date DESC;
    END IF;
END;
$$;

-- Get all productions
CREATE OR REPLACE FUNCTION core.get_all_productions()
RETURNS TABLE(
    id uuid,
    generator_id uuid,
    generator_capacity float,
    type_name varchar(20),
    isrenuevable boolean,
    date date,
    production_mw decimal
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, p.generator_id, g.capacity, t.name, t.isrenuevable, p.date, p.production_mw
    FROM core.production p
    JOIN core.generator g ON p.generator_id = g.id
    JOIN core.type t ON g.type = t.id
    ORDER BY p.date DESC, t.name;
END;
$$;

-- Get production by ID
CREATE OR REPLACE FUNCTION core.get_production_by_id(p_production_id uuid)
RETURNS TABLE(
    id uuid,
    generator_id uuid,
    generator_capacity float,
    type_name varchar(20),
    isrenuevable boolean,
    date date,
    production_mw decimal
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, p.generator_id, g.capacity, t.name, t.isrenuevable, p.date, p.production_mw
    FROM core.production p
    JOIN core.generator g ON p.generator_id = g.id
    JOIN core.type t ON g.type = t.id
    WHERE p.id = p_production_id;
END;
$$;

-- =====================================================
-- ANALYTICAL FUNCTIONS
-- =====================================================

-- Get total production by date range
CREATE OR REPLACE FUNCTION core.get_total_production_by_date_range(
    p_start_date date,
    p_end_date date DEFAULT NULL
)
RETURNS TABLE(
    date date,
    total_production decimal,
    renewable_production decimal,
    non_renewable_production decimal
)
LANGUAGE plpgsql
AS $$
BEGIN
    IF p_end_date IS NULL THEN
        p_end_date := CURRENT_DATE;
    END IF;

    RETURN QUERY
    SELECT 
        p.date,
        SUM(p.production_mw) as total_production,
        SUM(CASE WHEN t.isrenuevable = true THEN p.production_mw ELSE 0 END) as renewable_production,
        SUM(CASE WHEN t.isrenuevable = false THEN p.production_mw ELSE 0 END) as non_renewable_production
    FROM core.production p
    JOIN core.generator g ON p.generator_id = g.id
    JOIN core.type t ON g.type = t.id
    WHERE p.date >= p_start_date AND p.date <= p_end_date
    GROUP BY p.date
    ORDER BY p.date DESC;
END;
$$;

-- Get generator efficiency (production vs capacity)
CREATE OR REPLACE FUNCTION core.get_generator_efficiency(
    p_start_date date,
    p_end_date date DEFAULT NULL
)
RETURNS TABLE(
    generator_id uuid,
    type_name varchar(20),
    capacity float,
    total_production decimal,
    avg_daily_production decimal,
    efficiency_percentage decimal
)
LANGUAGE plpgsql
AS $$
BEGIN
    IF p_end_date IS NULL THEN
        p_end_date := CURRENT_DATE;
    END IF;

    RETURN QUERY
    SELECT 
        g.id as generator_id,
        t.name as type_name,
        g.capacity,
        SUM(p.production_mw) as total_production,
        AVG(p.production_mw) as avg_daily_production,
        ROUND((AVG(p.production_mw) / g.capacity * 100), 2) as efficiency_percentage
    FROM core.generator g
    JOIN core.type t ON g.type = t.id
    LEFT JOIN core.production p ON g.id = p.generator_id 
        AND p.date >= p_start_date AND p.date <= p_end_date
    GROUP BY g.id, t.name, g.capacity
    ORDER BY efficiency_percentage DESC NULLS LAST;
END;
$$;

-- Get renewable vs non-renewable summary
CREATE OR REPLACE FUNCTION core.get_renewable_summary(
    p_start_date date,
    p_end_date date DEFAULT NULL
)
RETURNS TABLE(
    energy_type varchar(20),
    total_capacity float,
    generator_count bigint,
    total_production decimal,
    avg_production decimal,
    percentage_of_total decimal
)
LANGUAGE plpgsql
AS $$
DECLARE
    total_all_production decimal;
BEGIN
    IF p_end_date IS NULL THEN
        p_end_date := CURRENT_DATE;
    END IF;

    -- Get total production for percentage calculation
    SELECT SUM(prod.production_mw) INTO total_all_production
    FROM core.production prod
    WHERE prod.date >= p_start_date AND prod.date <= p_end_date;

    IF total_all_production IS NULL THEN
        total_all_production := 0;
    END IF;

    RETURN QUERY
    SELECT 
        CASE WHEN t.isrenuevable = true THEN 'Renewable' ELSE 'Non-Renewable' END as energy_type,
        SUM(g.capacity) as total_capacity,
        COUNT(DISTINCT g.id) as generator_count,
        COALESCE(SUM(p.production_mw), 0) as total_production,
        COALESCE(AVG(p.production_mw), 0) as avg_production,
        CASE 
            WHEN total_all_production > 0 THEN ROUND((COALESCE(SUM(p.production_mw), 0) / total_all_production * 100), 2)
            ELSE 0 
        END as percentage_of_total
    FROM core.type t
    JOIN core.generator g ON t.id = g.type
    LEFT JOIN core.production p ON g.id = p.generator_id 
        AND p.date >= p_start_date AND p.date <= p_end_date
    GROUP BY t.isrenuevable
    ORDER BY energy_type;
END;
$$;
