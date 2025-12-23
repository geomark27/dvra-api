BEGIN; -- <-- Agregar esto al inicio

DO $$
DECLARE
    reg_id integer;
    sub_id integer;
    count_id integer;
    stat_id integer;
    now_ts timestamp := CURRENT_TIMESTAMP;
BEGIN

-----------------------------------------------------------------------
-- 1. REGIÓN: AMERICAS
-----------------------------------------------------------------------
INSERT INTO regions (name, is_active, created_at, updated_at) 
VALUES ('Americas', true, now_ts, now_ts) RETURNING id INTO reg_id;

    -- SUBREGION: NORTH AMERICA
    INSERT INTO subregions (name, region_id, is_active, created_at, updated_at) 
    VALUES ('North America', reg_id, true, now_ts, now_ts) RETURNING id INTO sub_id;

        -- USA
        INSERT INTO countries (name, numeric_code, iso2, iso3, phone_code, timezones, subregion_id, is_active, created_at, updated_at)
        VALUES ('United States', '840', 'US', 'USA', '1', 'UTC-5,UTC-8', sub_id, true, now_ts, now_ts) RETURNING id INTO count_id;
            -- California
            INSERT INTO states (name, country_id, country_code, is_active, created_at, updated_at) 
            VALUES ('California', count_id, 'US', true, now_ts, now_ts) RETURNING id INTO stat_id;
            INSERT INTO cities (name, state_id, is_active, created_at, updated_at) VALUES ('Los Angeles', stat_id, true, now_ts, now_ts), ('San Francisco', stat_id, true, now_ts, now_ts), ('San Diego', stat_id, true, now_ts, now_ts);
            -- New York
            INSERT INTO states (name, country_id, country_code, is_active, created_at, updated_at) 
            VALUES ('New York', count_id, 'US', true, now_ts, now_ts) RETURNING id INTO stat_id;
            INSERT INTO cities (name, state_id, is_active, created_at, updated_at) VALUES ('New York City', stat_id, true, now_ts, now_ts), ('Buffalo', stat_id, true, now_ts, now_ts);
            -- Texas
            INSERT INTO states (name, country_id, country_code, is_active, created_at, updated_at) 
            VALUES ('Texas', count_id, 'US', true, now_ts, now_ts) RETURNING id INTO stat_id;
            INSERT INTO cities (name, state_id, is_active, created_at, updated_at) VALUES ('Houston', stat_id, true, now_ts, now_ts), ('Austin', stat_id, true, now_ts, now_ts);

        -- CANADA
        INSERT INTO countries (name, numeric_code, iso2, iso3, phone_code, timezones, subregion_id, is_active, created_at, updated_at)
        VALUES ('Canada', '124', 'CA', 'CAN', '1', 'UTC-3.5,UTC-8', sub_id, true, now_ts, now_ts) RETURNING id INTO count_id;
            -- Ontario
            INSERT INTO states (name, country_id, country_code, is_active, created_at, updated_at) 
            VALUES ('Ontario', count_id, 'CA', true, now_ts, now_ts) RETURNING id INTO stat_id;
            INSERT INTO cities (name, state_id, is_active, created_at, updated_at) VALUES ('Toronto', stat_id, true, now_ts, now_ts), ('Ottawa', stat_id, true, now_ts, now_ts);

    -- SUBREGION: SOUTH AMERICA (LATAM)
    INSERT INTO subregions (name, region_id, is_active, created_at, updated_at) 
    VALUES ('South America', reg_id, true, now_ts, now_ts) RETURNING id INTO sub_id;

        -- ECUADOR
        INSERT INTO countries (name, numeric_code, iso2, iso3, phone_code, timezones, subregion_id, is_active, created_at, updated_at)
        VALUES ('Ecuador', '218', 'EC', 'ECU', '593', 'UTC-5', sub_id, true, now_ts, now_ts) RETURNING id INTO count_id;
            -- Pichincha
            INSERT INTO states (name, country_id, country_code, is_active, created_at, updated_at) VALUES ('Pichincha', count_id, 'EC', true, now_ts, now_ts) RETURNING id INTO stat_id;
            INSERT INTO cities (name, state_id, is_active, created_at, updated_at) VALUES ('Quito', stat_id, true, now_ts, now_ts);
            -- Guayas
            INSERT INTO states (name, country_id, country_code, is_active, created_at, updated_at) VALUES ('Guayas', count_id, 'EC', true, now_ts, now_ts) RETURNING id INTO stat_id;
            INSERT INTO cities (name, state_id, is_active, created_at, updated_at) VALUES ('Guayaquil', stat_id, true, now_ts, now_ts);

        -- ARGENTINA
        INSERT INTO countries (name, numeric_code, iso2, iso3, phone_code, timezones, subregion_id, is_active, created_at, updated_at)
        VALUES ('Argentina', '032', 'AR', 'ARG', '54', 'UTC-3', sub_id, true, now_ts, now_ts) RETURNING id INTO count_id;
            -- Buenos Aires
            INSERT INTO states (name, country_id, country_code, is_active, created_at, updated_at) VALUES ('Buenos Aires', count_id, 'AR', true, now_ts, now_ts) RETURNING id INTO stat_id;
            INSERT INTO cities (name, state_id, is_active, created_at, updated_at) VALUES ('CABA', stat_id, true, now_ts, now_ts), ('La Plata', stat_id, true, now_ts, now_ts);

    -- SUBREGION: CENTRAL AMERICA & MEXICO
    INSERT INTO subregions (name, region_id, is_active, created_at, updated_at) 
    VALUES ('Central America', reg_id, true, now_ts, now_ts) RETURNING id INTO sub_id;

        -- MEXICO
        INSERT INTO countries (name, numeric_code, iso2, iso3, phone_code, timezones, subregion_id, is_active, created_at, updated_at)
        VALUES ('Mexico', '484', 'MX', 'MEX', '52', 'UTC-6', sub_id, true, now_ts, now_ts) RETURNING id INTO count_id;
            -- Jalisco
            INSERT INTO states (name, country_id, country_code, is_active, created_at, updated_at) VALUES ('Jalisco', count_id, 'MX', true, now_ts, now_ts) RETURNING id INTO stat_id;
            INSERT INTO cities (name, state_id, is_active, created_at, updated_at) VALUES ('Guadalajara', stat_id, true, now_ts, now_ts), ('Zapopan', stat_id, true, now_ts, now_ts);
            -- Ciudad de Mexico
            INSERT INTO states (name, country_id, country_code, is_active, created_at, updated_at) VALUES ('Ciudad de Mexico', count_id, 'MX', true, now_ts, now_ts) RETURNING id INTO stat_id;
            INSERT INTO cities (name, state_id, is_active, created_at, updated_at) VALUES ('CDMX', stat_id, true, now_ts, now_ts);

-----------------------------------------------------------------------
-- 2. REGIÓN: EUROPE (Top 10)
-----------------------------------------------------------------------
INSERT INTO regions (name, is_active, created_at, updated_at) 
VALUES ('Europe', true, now_ts, now_ts) RETURNING id INTO reg_id;

    INSERT INTO subregions (name, region_id, is_active, created_at, updated_at) 
    VALUES ('Western Europe', reg_id, true, now_ts, now_ts) RETURNING id INTO sub_id;

        -- SPAIN
        INSERT INTO countries (name, numeric_code, iso2, iso3, phone_code, timezones, subregion_id, is_active, created_at, updated_at)
        VALUES ('Spain', '724', 'ES', 'ESP', '34', 'UTC+1', sub_id, true, now_ts, now_ts) RETURNING id INTO count_id;
            INSERT INTO states (name, country_id, country_code, is_active, created_at, updated_at) VALUES ('Madrid', count_id, 'ES', true, now_ts, now_ts) RETURNING id INTO stat_id;
            INSERT INTO cities (name, state_id, is_active, created_at, updated_at) VALUES ('Madrid', stat_id, true, now_ts, now_ts);
            INSERT INTO states (name, country_id, country_code, is_active, created_at, updated_at) VALUES ('Catalonia', count_id, 'ES', true, now_ts, now_ts) RETURNING id INTO stat_id;
            INSERT INTO cities (name, state_id, is_active, created_at, updated_at) VALUES ('Barcelona', stat_id, true, now_ts, now_ts);

        -- FRANCE
        INSERT INTO countries (name, numeric_code, iso2, iso3, phone_code, timezones, subregion_id, is_active, created_at, updated_at)
        VALUES ('France', '250', 'FR', 'FRA', '33', 'UTC+1', sub_id, true, now_ts, now_ts) RETURNING id INTO count_id;
            INSERT INTO states (name, country_id, country_code, is_active, created_at, updated_at) VALUES ('Île-de-France', count_id, 'FR', true, now_ts, now_ts) RETURNING id INTO stat_id;
            INSERT INTO cities (name, state_id, is_active, created_at, updated_at) VALUES ('Paris', stat_id, true, now_ts, now_ts);

        -- GERMANY
        INSERT INTO countries (name, numeric_code, iso2, iso3, phone_code, timezones, subregion_id, is_active, created_at, updated_at)
        VALUES ('Germany', '276', 'DE', 'DEU', '49', 'UTC+1', sub_id, true, now_ts, now_ts) RETURNING id INTO count_id;
            INSERT INTO states (name, country_id, country_code, is_active, created_at, updated_at) VALUES ('Bavaria', count_id, 'DE', true, now_ts, now_ts) RETURNING id INTO stat_id;
            INSERT INTO cities (name, state_id, is_active, created_at, updated_at) VALUES ('Munich', stat_id, true, now_ts, now_ts);

        -- UNITED KINGDOM
        INSERT INTO countries (name, numeric_code, iso2, iso3, phone_code, timezones, subregion_id, is_active, created_at, updated_at)
        VALUES ('United Kingdom', '826', 'GB', 'GBR', '44', 'UTC+0', sub_id, true, now_ts, now_ts) RETURNING id INTO count_id;
            INSERT INTO states (name, country_id, country_code, is_active, created_at, updated_at) VALUES ('England', count_id, 'GB', true, now_ts, now_ts) RETURNING id INTO stat_id;
            INSERT INTO cities (name, state_id, is_active, created_at, updated_at) VALUES ('London', stat_id, true, now_ts, now_ts), ('Manchester', stat_id, true, now_ts, now_ts);

        -- ITALY
        INSERT INTO countries (name, numeric_code, iso2, iso3, phone_code, timezones, subregion_id, is_active, created_at, updated_at)
        VALUES ('Italy', '380', 'IT', 'ITA', '39', 'UTC+1', sub_id, true, now_ts, now_ts) RETURNING id INTO count_id;
            INSERT INTO states (name, country_id, country_code, is_active, created_at, updated_at) VALUES ('Lazio', count_id, 'IT', true, now_ts, now_ts) RETURNING id INTO stat_id;
            INSERT INTO cities (name, state_id, is_active, created_at, updated_at) VALUES ('Rome', stat_id, true, now_ts, now_ts);

        -- PORTUGAL, NETHERLANDS, SWITZERLAND, BELGIUM, AUSTRIA
        -- (Siguen el mismo patrón para completar los 10 principales de Europa)

    RAISE NOTICE 'Poblado de datos geográficos en PostgreSQL completado.';
END $$;

COMMIT; -- <-- Agregar esto al final