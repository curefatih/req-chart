CREATE TABLE requests (
    id SERIAL PRIMARY KEY,
    created_at timestamptz,
    updated_at timestamptz,
    deleted_at timestamptz,
    type varchar,
    spent bigint,
    timestamp bigint
);
CREATE OR REPLACE FUNCTION public.notify_event() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN PERFORM pg_notify('new_data', row_to_json(NEW)::text);
RETURN NULL;
END;
$$;
CREATE TRIGGER newly_added_data
AFTER
INSERT ON requests FOR EACH ROW EXECUTE PROCEDURE notify_event();