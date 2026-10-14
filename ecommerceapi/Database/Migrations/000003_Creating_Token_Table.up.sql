Create Table Refreshtokens(
    token_id uuid default gen_random_uuid(),
    user_id uuid Unique,
    token_hash varchar ,
    family_id uuid Unique,
    is_used boolean,
    is_revoked boolean,
    expires_at timestamptz,
    created_at timestamptz default now(),
    updated_at timestamptz default now(),
    Constraint Pk_refreshtokens Primary key(token_id),
    Constraint Fk_refreshtokens Foreign key(user_id) 
    REFERENCES Users(user_id) on delete cascade

);