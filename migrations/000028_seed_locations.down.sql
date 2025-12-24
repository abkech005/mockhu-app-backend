-- Remove seeded locations
DELETE FROM locations WHERE used_by_count = 0;
