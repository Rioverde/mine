BEGIN;

CREATE TABLE IF NOT EXISTS items (
    id              UUID        PRIMARY KEY,
    name            TEXT        NOT NULL,
    quality         INTEGER     NOT NULL CHECK (quality BETWEEN 0 AND 100),
    ore_material    TEXT        NOT NULL CHECK (ore_material IN ('iron', 'copper', 'gold')),
    ore_capacity    INTEGER     NOT NULL,
    ingot_quality   INTEGER     NOT NULL CHECK (ingot_quality BETWEEN 0 AND 100),
    factory         TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_items_created       ON items(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_items_ore_material  ON items(ore_material);

CREATE TABLE IF NOT EXISTS outbox (
    id          BIGSERIAL   PRIMARY KEY,
    item_id     UUID        NOT NULL REFERENCES items(id),
    payload     JSONB       NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    sent_at     TIMESTAMPTZ
);

-- Partial index makes "fetch pending" trivial — only scans unsent rows
CREATE INDEX IF NOT EXISTS idx_outbox_pending ON outbox(id) WHERE sent_at IS NULL;

-- Partial index for cleanup: only sent rows, ordered by sent_at
CREATE INDEX IF NOT EXISTS idx_outbox_sent ON outbox(sent_at) WHERE sent_at IS NOT NULL;

COMMIT;
