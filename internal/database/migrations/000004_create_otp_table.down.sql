-- Drop OTP table and verification column
DROP TABLE IF EXISTS otps;
ALTER TABLE users DROP COLUMN IF EXISTS is_verified;
