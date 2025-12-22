-- Seed 100 test users with varied interests
-- Password for all users: Password123! (bcrypt hash)

DO $$
DECLARE
    user_id UUID;
    i INTEGER;
    first_names TEXT[] := ARRAY['Emma', 'Liam', 'Olivia', 'Noah', 'Ava', 'William', 'Sophia', 'James', 'Isabella', 'Oliver', 'Mia', 'Benjamin', 'Charlotte', 'Elijah', 'Amelia', 'Lucas', 'Harper', 'Mason', 'Evelyn', 'Logan', 'Aria', 'Alexander', 'Luna', 'Ethan', 'Chloe', 'Jacob', 'Penelope', 'Michael', 'Layla', 'Daniel', 'Riley', 'Henry', 'Zoey', 'Jackson', 'Nora', 'Sebastian', 'Lily', 'Aiden', 'Eleanor', 'Matthew', 'Hannah', 'Samuel', 'Lillian', 'David', 'Addison', 'Joseph', 'Aubrey', 'Carter', 'Ellie', 'Owen'];
    last_names TEXT[] := ARRAY['Smith', 'Johnson', 'Williams', 'Brown', 'Jones', 'Garcia', 'Miller', 'Davis', 'Rodriguez', 'Martinez', 'Hernandez', 'Lopez', 'Gonzalez', 'Wilson', 'Anderson', 'Thomas', 'Taylor', 'Moore', 'Jackson', 'Martin', 'Lee', 'Perez', 'Thompson', 'White', 'Harris', 'Sanchez', 'Clark', 'Ramirez', 'Lewis', 'Robinson', 'Walker', 'Young', 'Allen', 'King', 'Wright', 'Scott', 'Torres', 'Nguyen', 'Hill', 'Flores', 'Green', 'Adams', 'Nelson', 'Baker', 'Hall', 'Rivera', 'Campbell', 'Mitchell', 'Carter', 'Roberts'];
    places TEXT[] := ARRAY['New York', 'Los Angeles', 'Chicago', 'Houston', 'Phoenix', 'San Francisco', 'Seattle', 'Denver', 'Boston', 'Austin', 'Miami', 'Atlanta', 'Portland', 'San Diego', 'Dallas', 'London', 'Toronto', 'Sydney', 'Mumbai', 'Berlin'];
    bios TEXT[] := ARRAY['Software Developer', 'Product Manager', 'UX Designer', 'Data Scientist', 'Full Stack Developer', 'Mobile Developer', 'DevOps Engineer', 'Tech Lead', 'Startup Founder', 'Student', 'Freelancer', 'Content Creator', 'Digital Marketer', 'Entrepreneur', 'Open Source Contributor'];
    password_hash TEXT := '$2a$10$N9qo8uLOickgx2ZMRZoMye.IIuVRdVg.MvKhKl9mGZ3AZaLrCDR9e';
    interest_ids UUID[];
    random_interests UUID[];
    num_interests INTEGER;
BEGIN
    -- Get all interest IDs
    SELECT ARRAY_AGG(id) INTO interest_ids FROM interests;
    
    FOR i IN 1..100 LOOP
        user_id := gen_random_uuid();
        
        -- Insert user
        INSERT INTO users (
            id, email, email_verified, phone, phone_verified, username, 
            first_name, last_name, dob, password_hash, avatar_url, 
            bio, place, is_active, onboarding_completed, created_at, updated_at
        ) VALUES (
            user_id,
            'user' || i || '@mockhu.com',
            (i % 3 = 0), -- Every 3rd user has verified email
            CASE WHEN i % 4 = 0 THEN '+1555' || LPAD(i::text, 7, '0') ELSE NULL END,
            (i % 5 = 0), -- Every 5th user has verified phone
            'user_' || i,
            first_names[(i % 50) + 1],
            last_names[(i % 50) + 1],
            ('1990-01-01'::date + (i * 30)::integer),
            password_hash,
            'https://i.pravatar.cc/150?img=' || ((i % 70) + 1),
            bios[(i % 15) + 1] || ' | Loves to connect',
            places[(i % 20) + 1],
            true,
            true,
            NOW() - ((100 - i) || ' days')::interval,
            NOW()
        );
        
        -- Assign 2-5 random interests to each user
        num_interests := 2 + (i % 4);
        
        INSERT INTO user_interests (user_id, interest_id, created_at)
        SELECT user_id, id, NOW()
        FROM interests
        ORDER BY RANDOM()
        LIMIT num_interests
        ON CONFLICT DO NOTHING;
        
    END LOOP;
END $$;
