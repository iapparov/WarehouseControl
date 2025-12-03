CREATE TABLE history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id     UUID NOT NULL,
    action      VARCHAR(20) NOT NULL,
    changed_by  UUID NOT NULL,
    changed_by_login VARCHAR(100) NOT NULL,
    changed_at  TIMESTAMP NOT NULL DEFAULT now(),
    old_data    JSONB,
    new_data    JSONB
);