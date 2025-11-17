-- Listup sequences
SELECT 
    s.schemaname,
    s.sequencename,
    d.refobjid::regclass AS owned_by_table,
    a.attname AS owned_by_column
FROM 
    pg_sequences s
LEFT JOIN 
    pg_depend d ON d.objid = (s.schemaname || '.' || s.sequencename)::regclass::oid
LEFT JOIN 
    pg_attribute a ON a.attrelid = d.refobjid AND a.attnum = d.refobjsubid
WHERE 
    d.deptype = 'a'  -- 'a' = auto dependency (ownership)
ORDER BY 
    s.schemaname, s.sequencename;

-- remove ownership
DO $$
DECLARE
    rec RECORD;
BEGIN
    FOR rec IN 
        SELECT 
            n.nspname AS schema_name,
            c.relname AS sequence_name
        FROM 
            pg_class c
        JOIN 
            pg_namespace n ON n.oid = c.relnamespace
        JOIN 
            pg_depend d ON d.objid = c.oid AND d.deptype = 'a'
        WHERE 
            c.relkind = 'S'
    LOOP
        EXECUTE format('ALTER SEQUENCE %I.%I OWNED BY NONE', 
                      rec.schema_name, 
                      rec.sequence_name);
        RAISE NOTICE 'Removed ownership from %.%', rec.schema_name, rec.sequence_name;
    END LOOP;
END $$;
