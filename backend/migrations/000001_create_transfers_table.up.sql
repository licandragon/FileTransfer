CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    download_token TEXT NOT NULL UNIQUE,
    download_count INTEGER DEFAULT 0,
    sender_email TEXT NOT NULL,
    subject_email TEXT,
    message_email TEXT,
    recipients TEXT[] DEFAULT '{}',
    user_id UUID DEFAULT NULL,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
