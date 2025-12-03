CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE OR REPLACE FUNCTION app_current_user()
RETURNS uuid AS $$
BEGIN
  RETURN current_setting('app.current_user', true)::uuid;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION app_current_user_login()
RETURNS text AS $$
BEGIN
  RETURN current_setting('app.current_user_login', true);
END;
$$ LANGUAGE plpgsql;