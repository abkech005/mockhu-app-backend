-- Add media column to postfeeds (array of media URLs from Cloudflare R2)
ALTER TABLE postfeeds ADD COLUMN media JSONB DEFAULT '[]';

-- Index for media queries
CREATE INDEX idx_postfeeds_media ON postfeeds USING GIN(media);

COMMENT ON COLUMN postfeeds.media IS 'Array of media objects: [{url, type, thumbnail_url, width, height}]';
