-- Seed default admin-defined titles
INSERT INTO titles (name, description, defined_by, used_by_count) VALUES
    ('Student', 'Currently pursuing education', 'admin', 0),
    ('Teacher', 'Educator or instructor', 'admin', 0),
    ('CAT Aspirant', 'Preparing for CAT entrance exam', 'admin', 0),
    ('UPSC Aspirant', 'Preparing for UPSC Civil Services exam', 'admin', 0),
    ('GATE Aspirant', 'Preparing for GATE entrance exam', 'admin', 0),
    ('JEE Aspirant', 'Preparing for JEE entrance exam', 'admin', 0),
    ('NEET Aspirant', 'Preparing for NEET entrance exam', 'admin', 0),
    ('SSC Aspirant', 'Preparing for SSC exams', 'admin', 0),
    ('Banking Aspirant', 'Preparing for banking exams (IBPS, SBI, RBI)', 'admin', 0),
    ('Working Professional', 'Currently employed professional', 'admin', 0),
    ('Fresher', 'Recent graduate looking for opportunities', 'admin', 0),
    ('Mentor', 'Guides and helps other learners', 'admin', 0),
    ('Content Creator', 'Creates educational content', 'admin', 0),
    ('Researcher', 'Engaged in academic research', 'admin', 0),
    ('Freelancer', 'Self-employed professional', 'admin', 0),
    ('Entrepreneur', 'Building or running a business', 'admin', 0)
ON CONFLICT (name) DO NOTHING;
