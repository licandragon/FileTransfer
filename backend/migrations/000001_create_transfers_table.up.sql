CREATE TABLE IF NOT EXISTS transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    download_token TEXT NOT NULL UNIQUE,
    sender_email TEXT NOT NULL,
    message_email TEXT NOT NULL,
    recipients JSONB DEFAULT '[]',
    user_id UUID DEFAULT NULL,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transfers_download_token ON transfers(download_token);