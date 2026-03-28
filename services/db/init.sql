CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    first_name TEXT NOT NULL
);

INSERT INTO accounts (email, hashed_password, first_name)
VALUES (
    'river@triangleoflove.app',
    '$2a$10$pLlNqI6u3aN1f4qqRCI/huMGLpCgSUd43jon6WxrdYDw2878DeAEi',
    'River'
);
