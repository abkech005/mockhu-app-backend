-- Seed test users with hashed password (password: Password123!)
-- Hash generated with bcrypt cost 10

INSERT INTO users (id, email, email_verified, phone, phone_verified, username, first_name, last_name, middle_name, dob, password_hash, avatar_url, bio, place, is_active, onboarding_completed, created_at, updated_at)
VALUES 
    ('11111111-1111-1111-1111-111111111111', 'john.doe@example.com', true, '+1234567890', false, 'johndoe_seed', 'John', 'Doe', NULL, '1990-05-15', '$2a$10$rQ7rV7ZKDQx.VqG4N2VQz.E3KzKzVQZ5BBzKzKzVQZ5BBzKzKzV', 'https://i.pravatar.cc/150?img=1', 'Software Developer from NYC', 'New York, USA', true, true, NOW(), NOW()),
    ('22222222-2222-2222-2222-222222222222', 'jane.smith@example.com', true, '+1987654321', true, 'janesmith_seed', 'Jane', 'Smith', 'Marie', '1992-08-20', '$2a$10$rQ7rV7ZKDQx.VqG4N2VQz.E3KzKzVQZ5BBzKzKzVQZ5BBzKzKzV', 'https://i.pravatar.cc/150?img=5', 'Product Manager | Tech Enthusiast', 'San Francisco, USA', true, true, NOW(), NOW()),
    ('33333333-3333-3333-3333-333333333333', 'bob.wilson@example.com', true, NULL, false, 'bobwilson_seed', 'Bob', 'Wilson', NULL, '1988-12-01', '$2a$10$rQ7rV7ZKDQx.VqG4N2VQz.E3KzKzVQZ5BBzKzKzVQZ5BBzKzKzV', 'https://i.pravatar.cc/150?img=8', 'Full Stack Developer', 'London, UK', true, true, NOW(), NOW()),
    ('44444444-4444-4444-4444-444444444444', 'alice.johnson@example.com', true, '+1555123456', true, 'alicejohnson_seed', 'Alice', 'Johnson', 'Rose', '1995-03-25', '$2a$10$rQ7rV7ZKDQx.VqG4N2VQz.E3KzKzVQZ5BBzKzKzVQZ5BBzKzKzV', 'https://i.pravatar.cc/150?img=9', 'UX Designer | Coffee Lover', 'Toronto, Canada', true, true, NOW(), NOW()),
    ('55555555-5555-5555-5555-555555555555', 'mike.brown@example.com', false, '+1666789012', false, 'mikebrown_seed', 'Mike', 'Brown', NULL, '1993-07-10', '$2a$10$rQ7rV7ZKDQx.VqG4N2VQz.E3KzKzVQZ5BBzKzKzVQZ5BBzKzKzV', 'https://i.pravatar.cc/150?img=11', 'Backend Engineer', 'Sydney, Australia', true, false, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Add some interests to seeded users
INSERT INTO user_interests (user_id, interest_id, created_at)
SELECT '11111111-1111-1111-1111-111111111111', id, NOW() FROM interests WHERE slug IN ('javascript', 'react', 'nodejs')
ON CONFLICT DO NOTHING;

INSERT INTO user_interests (user_id, interest_id, created_at)
SELECT '22222222-2222-2222-2222-222222222222', id, NOW() FROM interests WHERE slug IN ('python', 'machine-learning', 'data-science')
ON CONFLICT DO NOTHING;

INSERT INTO user_interests (user_id, interest_id, created_at)
SELECT '33333333-3333-3333-3333-333333333333', id, NOW() FROM interests WHERE slug IN ('golang', 'docker', 'kubernetes')
ON CONFLICT DO NOTHING;

INSERT INTO user_interests (user_id, interest_id, created_at)
SELECT '44444444-4444-4444-4444-444444444444', id, NOW() FROM interests WHERE slug IN ('ui-design', 'figma', 'css')
ON CONFLICT DO NOTHING;
