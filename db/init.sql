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

CREATE TABLE IF NOT EXISTS exporter_checkpoints (
    checkpoint_name  TEXT        PRIMARY KEY,
    last_timestamp   TIMESTAMPTZ NOT NULL,
    last_id          UUID        NOT NULL,
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO exporter_checkpoints (checkpoint_name, last_timestamp, last_id)
VALUES ('auction_merchant',  '1970-01-01 00:00:00+00', '00000000-0000-0000-0000-000000000000')
ON CONFLICT (checkpoint_name) DO NOTHING;

COMMIT;
