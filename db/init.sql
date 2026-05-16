-- Schema for the mine storage.
-- Denormalized: each item row carries its ore + ingot provenance inline.

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

COMMIT;
