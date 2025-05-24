ALTER TABLE users
ADD COLUMN password VARCHAR(255);

-- Update provider_id to be nullable since local users won't have it
ALTER TABLE users
ALTER COLUMN provider_id DROP NOT NULL; 