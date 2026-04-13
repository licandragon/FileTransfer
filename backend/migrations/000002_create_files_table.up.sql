CREATE TABLE IF NOT EXISTS files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transfer_id UUID NOT NULL,
    filename TEXT NOT NULL,
    original_name TEXT NOT NULL,
    size BIGINT NOT NULL,
    mime_type TEXT,
    storage_path TEXT NOT NULL,
    bucket TEXT NOT NULL DEFAULT 'transfers',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT fk_transfer
        FOREIGN KEY (transfer_id)
        REFERENCES transfers(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_files_transfer_id ON files(transfer_id);