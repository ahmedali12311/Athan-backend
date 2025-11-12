-- Functions
CREATE OR REPLACE FUNCTION app_func_insert_user_stuff()
RETURNS TRIGGER
LANGUAGE plpgsql
AS
$$
BEGIN
    IF OLD.id IS NULL
    THEN
        INSERT INTO wallets (id) VALUES (NEW.id);
    END IF;

    RETURN NEW;
END;
$$;

-- TRIGGERS

-- trigger: wallets_change
CREATE TRIGGER app_trigger_user_insert
AFTER INSERT ON users FOR EACH ROW
EXECUTE PROCEDURE app_func_insert_user_stuff();
