-- Remove seeded users (users with email pattern user*@mockhu.com)
DELETE FROM users WHERE email LIKE 'user%@mockhu.com';
