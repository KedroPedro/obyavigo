-- Migration: Add password_reset_tokens table
-- Date: 2025-11-03
-- Description: Creates table for storing password reset tokens

CREATE TABLE IF NOT EXISTS site.password_reset_tokens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    token text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    CONSTRAINT password_reset_tokens_pkey PRIMARY KEY (id),
    CONSTRAINT password_reset_tokens_token_unique UNIQUE (token),
    CONSTRAINT password_reset_tokens_user_fk FOREIGN KEY (user_id) 
        REFERENCES site.users(id) ON DELETE CASCADE
);

-- Index for faster token lookups
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_token 
    ON site.password_reset_tokens(token);

-- Index for faster expiration checks
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_expires_at 
    ON site.password_reset_tokens(expires_at);

-- Clean up expired tokens periodically (optional, can be done via cron job)
-- This is just a helper query to run manually or via scheduled job:
-- DELETE FROM site.password_reset_tokens WHERE expires_at < NOW();
