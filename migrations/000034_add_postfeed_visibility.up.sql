-- Add visibility column to postfeeds
ALTER TABLE postfeeds ADD COLUMN visibility VARCHAR(20) NOT NULL DEFAULT 'public' 
    CHECK (visibility IN ('public', 'private', 'followers_only'));

-- Index for visibility filtering
CREATE INDEX idx_postfeeds_visibility ON postfeeds(visibility);

COMMENT ON COLUMN postfeeds.visibility IS 'Post visibility: public, private, followers_only';
