CREATE SCHEMA IF NOT EXISTS simpleext;

CREATE OR REPLACE FUNCTION simpleext.hello_world()
RETURNS text
LANGUAGE sql
AS $$
  SELECT 'hello from simpleext';
$$;

CREATE OR REPLACE FUNCTION simpleext.add_ints(a int, b int)
RETURNS int
LANGUAGE sql
AS $$
  SELECT a + b;
$$;
