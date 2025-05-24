-- Make provider_id not nullable again
ALTER TABLE users
ALTER COLUMN provider_id SET NOT NULL;

-- Remove password column
ALTER TABLE users
DROP COLUMN password; 