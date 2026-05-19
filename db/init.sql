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

CREATE INDEX IF NOT EXISTS idx_outbox_pending ON outbox(id) WHERE sent_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_outbox_sent ON outbox(sent_at) WHERE sent_at IS NOT NULL;

CREATE TABLE IF NOT EXISTS treasury (
    id INT PRIMARY KEY DEFAULT 1,
    balance_copper BIGINT NOT NULL DEFAULT 10000000,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT single_row CHECK (id = 1)
);

CREATE TABLE IF NOT EXISTS contracts (
    id UUID PRIMARY KEY,
    target_item VARCHAR(50) NOT NULL,
    target_material VARCHAR(50) NOT NULL,
    min_quality INT NOT NULL DEFAULT 0,
    required_qty INT NOT NULL,
    reward_copper BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

COMMIT;
