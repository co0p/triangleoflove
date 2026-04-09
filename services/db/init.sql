CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    first_name TEXT NOT NULL,
    invite_code TEXT NULL
);

CREATE TABLE couples (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id_a UUID NOT NULL REFERENCES accounts(id),
    account_id_b UUID NOT NULL REFERENCES accounts(id),
    formed_on TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (account_id_a, account_id_b)
);

CREATE TABLE checkins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id),
    date DATE NOT NULL,
    felt_close SMALLINT NULL CHECK (felt_close IS NULL OR felt_close BETWEEN -5 AND 5),
    positive_energy SMALLINT NULL CHECK (positive_energy IS NULL OR positive_energy BETWEEN -5 AND 5),
    supported SMALLINT NULL CHECK (supported IS NULL OR supported BETWEEN -5 AND 5),
    communication_healthy SMALLINT NULL CHECK (communication_healthy IS NULL OR communication_healthy BETWEEN -5 AND 5),
    stress_level SMALLINT NULL CHECK (stress_level IS NULL OR stress_level BETWEEN -5 AND 5),
    note TEXT NOT NULL DEFAULT '',
    saved_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (account_id, date)
);

CREATE INDEX checkins_account_date_idx ON checkins (account_id, date);

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

INSERT INTO checkins (account_id, date, felt_close, positive_energy, supported, communication_healthy, stress_level, note)
SELECT id, DATE '2026-03-25', 2, 1, 3, 2, -1, 'Good start to the week, felt connected.'       FROM accounts WHERE email = 'river@triangleoflove.app' UNION ALL
SELECT id, DATE '2026-03-26', 1, 2, 2, 1,  0, 'A bit distracted but overall okay.'              FROM accounts WHERE email = 'river@triangleoflove.app' UNION ALL
SELECT id, DATE '2026-03-27', 3, 3, 4, 3, -2, 'Really nice evening together.'                   FROM accounts WHERE email = 'river@triangleoflove.app' UNION ALL
SELECT id, DATE '2026-03-28', 0, 0, 1, 0,  2, 'Stressful day at work, felt a bit distant.'      FROM accounts WHERE email = 'river@triangleoflove.app' UNION ALL
SELECT id, DATE '2026-03-29', 2, 2, 2, 2,  0, 'Recovered well, good talk in the evening.'       FROM accounts WHERE email = 'river@triangleoflove.app' UNION ALL
SELECT id, DATE '2026-03-30', 4, 3, 4, 4, -3, 'Best day this week, really close and present.'   FROM accounts WHERE email = 'river@triangleoflove.app' UNION ALL
SELECT id, DATE '2026-03-31', 3, 2, 3, 3, -2, 'Wrapping up the month on a high note.'           FROM accounts WHERE email = 'river@triangleoflove.app';
