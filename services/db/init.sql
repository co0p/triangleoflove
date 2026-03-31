CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    first_name TEXT NOT NULL
);

CREATE TABLE checkins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id),
    date DATE NOT NULL,
    felt_close SMALLINT NOT NULL CHECK (felt_close = -1 OR felt_close BETWEEN 1 AND 10),
    positive_energy SMALLINT NOT NULL CHECK (positive_energy = -1 OR positive_energy BETWEEN 1 AND 10),
    supported SMALLINT NOT NULL CHECK (supported = -1 OR supported BETWEEN 1 AND 10),
    communication_healthy SMALLINT NOT NULL CHECK (communication_healthy = -1 OR communication_healthy BETWEEN 1 AND 10),
    stress_level SMALLINT NOT NULL CHECK (stress_level = -1 OR stress_level BETWEEN 1 AND 10),
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

INSERT INTO checkins (account_id, date, felt_close, positive_energy, supported, communication_healthy, stress_level, note)
SELECT id, '2026-03-25', 7, 6, 8, 7, 4, 'Good start to the week, felt connected.'       FROM accounts WHERE email = 'river@triangleoflove.app' UNION ALL
SELECT id, '2026-03-26', 6, 7, 7, 6, 5, 'A bit distracted but overall okay.'              FROM accounts WHERE email = 'river@triangleoflove.app' UNION ALL
SELECT id, '2026-03-27', 8, 8, 9, 8, 3, 'Really nice evening together.'                   FROM accounts WHERE email = 'river@triangleoflove.app' UNION ALL
SELECT id, '2026-03-28', 5, 5, 6, 5, 7, 'Stressful day at work, felt a bit distant.'      FROM accounts WHERE email = 'river@triangleoflove.app' UNION ALL
SELECT id, '2026-03-29', 7, 7, 7, 7, 5, 'Recovered well, good talk in the evening.'       FROM accounts WHERE email = 'river@triangleoflove.app' UNION ALL
SELECT id, '2026-03-30', 9, 8, 9, 9, 2, 'Best day this week, really close and present.'   FROM accounts WHERE email = 'river@triangleoflove.app' UNION ALL
SELECT id, '2026-03-31', 8, 7, 8, 8, 3, 'Wrapping up the month on a high note.'           FROM accounts WHERE email = 'river@triangleoflove.app';
