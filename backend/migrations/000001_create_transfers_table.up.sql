CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; //Cambio en extension para uuid
CREATE TABLE IF NOT EXISTS transfers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    download_token UUID NOT NULL UNIQUE,
    upload_token UUID NOT NULL UNIQUE,
    sender_email TEXT NOT NULL,
    subject_email TEXT,
    message_email TEXT,
    recipients TEXT[] DEFAULT '{}',
    user_id UUID DEFAULT NULL,
    status_transfer VARCHAR(20) NOT NULL DEFAULT 'pending',
    download_count INTEGER DEFAULT 0,
    total_file INTEGER NOT NULL DEFAULT 0,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
