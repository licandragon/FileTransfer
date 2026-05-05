CREATE TABLE IF NOT EXISTS files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transfer_id UUID NOT NULL,
    file_name TEXT NOT NULL,
    original_name TEXT NOT NULL,
    size_file BIGINT NOT NULL,
    mime_type TEXT,
    storage_path TEXT NOT NULL,
    file_index INTEGER NOT NULL DEFAULT 0,
    bucket TEXT NOT NULL DEFAULT 'transfers',
    status_file VARCHAR(20) NOT NULL DEFAULT 'uploaded',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT fk_transfer
        FOREIGN KEY (transfer_id)
        REFERENCES transfers(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_files_transfer_id ON files(transfer_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_files_transfer_file_index ON files(transfer_id, file_index)