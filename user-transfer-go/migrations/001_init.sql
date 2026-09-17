-- Skema database untuk User Transfer API
-- Menggunakan SERIAL (auto-increment bawaan Postgres) untuk primary key,
-- supaya sederhana dan tidak perlu generate UUID manual.

CREATE TABLE IF NOT EXISTS users (
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(255)   NOT NULL,
    age     INTEGER        NOT NULL CHECK (age > 0),
    balance NUMERIC(18, 2) NOT NULL DEFAULT 0 CHECK (balance >= 0)
);

CREATE TABLE IF NOT EXISTS transfers (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_user_id  INTEGER        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    nominal         NUMERIC(18, 2) NOT NULL CHECK (nominal > 0),
    created_date    TIMESTAMPTZ    NOT NULL DEFAULT now()
);

-- Index untuk mempercepat query riwayat transfer per user
-- (dipakai di GET /users/{user_id}/transfers, mencari sebagai pengirim
-- maupun penerima).
CREATE INDEX IF NOT EXISTS idx_transfers_user_id ON transfers (user_id);
CREATE INDEX IF NOT EXISTS idx_transfers_target_user_id ON transfers (target_user_id);
