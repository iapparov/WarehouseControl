DROP TRIGGER IF EXISTS item_insert_history ON items;
DROP TRIGGER IF EXISTS item_update_history ON items;
DROP TRIGGER IF EXISTS item_delete_history ON items;

DROP FUNCTION IF EXISTS trg_item_insert();
DROP FUNCTION IF EXISTS trg_item_update();
DROP FUNCTION IF EXISTS trg_item_delete();