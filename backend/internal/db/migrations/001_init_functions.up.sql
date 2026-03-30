CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION generate_conversation_display_id()
RETURNS TRIGGER AS $$
DECLARE
    seq_name TEXT;
    next_id  BIGINT;
BEGIN
    seq_name := 'conv_disp_seq_' || NEW."account_id";
    EXECUTE format('CREATE SEQUENCE IF NOT EXISTS %I START 1', seq_name);
    EXECUTE format('SELECT nextval(%L)', seq_name) INTO next_id;
    NEW."display_id" := next_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
