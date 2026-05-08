CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    first_name TEXT NOT NULL,
    invite_code TEXT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE couples (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id_a UUID NOT NULL REFERENCES accounts(id),
    account_id_b UUID NOT NULL REFERENCES accounts(id),
    formed_on TIMESTAMPTZ NOT NULL DEFAULT now(),
    ended_on TIMESTAMPTZ NULL,
    UNIQUE (account_id_a, account_id_b)
);

CREATE TABLE checkins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id),
    date DATE NOT NULL,
    felt_understood   SMALLINT NOT NULL DEFAULT 0 CHECK (felt_understood   BETWEEN 0 AND 5),
    meaningful_sharing SMALLINT NOT NULL DEFAULT 0 CHECK (meaningful_sharing BETWEEN 0 AND 5),
    could_count_on_them SMALLINT NOT NULL DEFAULT 0 CHECK (could_count_on_them BETWEEN 0 AND 5),
    effort_for_us     SMALLINT NOT NULL DEFAULT 0 CHECK (effort_for_us     BETWEEN 0 AND 5),
    desire            SMALLINT NOT NULL DEFAULT 0 CHECK (desire            BETWEEN 0 AND 5),
    spark             SMALLINT NOT NULL DEFAULT 0 CHECK (spark             BETWEEN 0 AND 5),
    mood              SMALLINT NOT NULL DEFAULT 0 CHECK (mood              BETWEEN 0 AND 5),
    note TEXT NOT NULL DEFAULT '',
    saved_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (account_id, date)
);

CREATE INDEX checkins_account_date_idx ON checkins (account_id, date);

INSERT INTO accounts (email, hashed_password, first_name, role)
VALUES (
    'admin@triangleoflove.app',
    '$2a$10$bfOlmuMfOzSg6mYVbPPYX.Hj7ktgS7jKVtzydtsjJq2i.gXJIhcue',
    'Admin',
    'admin'
);

INSERT INTO accounts (email, hashed_password, first_name)
VALUES (
    'river@triangleoflove.app',
    '$2a$10$pLlNqI6u3aN1f4qqRCI/huMGLpCgSUd43jon6WxrdYDw2878DeAEi',
    'River'
);

INSERT INTO accounts (email, hashed_password, first_name)
VALUES (
    'jordan@triangleoflove.app',
    '$2y$10$cxtnNv8hGdXzDGr11z9BGe30GCFoWI0cWURq1pgyedcWXW82pOyRi',
    'Jordan'
);

INSERT INTO accounts (email, hashed_password, first_name)
VALUES (
    'alex@triangleoflove.app',
    '$2a$10$K3sohwO.ZVaoahUghShek.T3PZ6ROKL0YbyRdjB4NMvJfLlNKomy.',
    'Alex'
);

INSERT INTO accounts (email, hashed_password, first_name)
VALUES (
    'sam@triangleoflove.app',
    '$2a$10$nsLfQTAtYuUvr6fyyeulDOGdiRIAgGr9VxVl1g2k3/2OpYhGBOInW',
    'Sam'
);

-- Alex and Sam are pre-formed as a couple so GivenPaired_* tests have
-- isolated, deterministic paired state independent of test ordering.
INSERT INTO couples (account_id_a, account_id_b, formed_on)
SELECT a.id, s.id, now()
FROM accounts a, accounts s
WHERE a.email = 'alex@triangleoflove.app'
  AND s.email = 'sam@triangleoflove.app';

INSERT INTO checkins (account_id, date, felt_understood, meaningful_sharing, could_count_on_them, effort_for_us, desire, spark, mood, note)
SELECT id, DATE '2026-03-25', 4, 3, 5, 4, 3, 4, 4, 'Good start to the week, felt connected.'     FROM accounts WHERE email = 'river@triangleoflove.app' UNION ALL
SELECT id, DATE '2026-03-26', 3, 2, 4, 3, 2, 3, 3, 'A bit distracted but overall okay.'           FROM accounts WHERE email = 'river@triangleoflove.app' UNION ALL
SELECT id, DATE '2026-03-27', 5, 4, 5, 5, 4, 5, 5, 'Really nice evening together.'                FROM accounts WHERE email = 'river@triangleoflove.app' UNION ALL
SELECT id, DATE '2026-03-28', 2, 1, 2, 2, 1, 2, 2, 'Stressful day at work, felt a bit distant.'   FROM accounts WHERE email = 'river@triangleoflove.app' UNION ALL
SELECT id, DATE '2026-03-29', 4, 3, 4, 4, 3, 3, 4, 'Recovered well, good talk in the evening.'    FROM accounts WHERE email = 'river@triangleoflove.app' UNION ALL
SELECT id, DATE '2026-03-30', 5, 5, 5, 5, 5, 5, 5, 'Best day this week, really close and present.' FROM accounts WHERE email = 'river@triangleoflove.app' UNION ALL
SELECT id, DATE '2026-03-31', 4, 4, 4, 4, 4, 4, 4, 'Wrapping up the month on a high note.'        FROM accounts WHERE email = 'river@triangleoflove.app';
