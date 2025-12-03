CREATE OR REPLACE FUNCTION trg_item_insert()
RETURNS trigger AS $$
DECLARE
    uid UUID;
    login TEXT;
BEGIN
    uid := app_current_user();
    login := app_current_user_login();
    INSERT INTO history(item_id, action, changed_by, changed_by_login, old_data, new_data)
    VALUES (NEW.id, 'created', uid, login, NULL, to_jsonb(NEW));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION trg_item_update()
RETURNS trigger AS $$
DECLARE
    uid UUID;
    login TEXT;
BEGIN
    uid := app_current_user();
    login := app_current_user_login();
    INSERT INTO history(item_id, action, changed_by, changed_by_login, old_data, new_data)
    VALUES (NEW.id, 'updated', uid, login, to_jsonb(OLD), to_jsonb(NEW));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION trg_item_delete()
RETURNS trigger AS $$
DECLARE
    uid UUID;
    login TEXT;
BEGIN
    uid := app_current_user();
    login := app_current_user_login();
    INSERT INTO history(item_id, action, changed_by, changed_by_login, old_data, new_data)
    VALUES (OLD.id, 'deleted', uid, login, to_jsonb(OLD), NULL);
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER item_insert_history
AFTER INSERT ON items
FOR EACH ROW EXECUTE FUNCTION trg_item_insert();

CREATE TRIGGER item_update_history
AFTER UPDATE ON items
FOR EACH ROW EXECUTE FUNCTION trg_item_update();

CREATE TRIGGER item_delete_history
AFTER DELETE ON items
FOR EACH ROW EXECUTE FUNCTION trg_item_delete();