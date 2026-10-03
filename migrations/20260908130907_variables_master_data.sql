-- +goose Up

INSERT INTO variables
    (key, label, description, source, data_type)
VALUES
    ('contact_first_name', 'First Name', 'Contact first name', 'CONTACT', 'TEXT'),
    ('contact_last_name', 'Last Name', 'Contact last name', 'CONTACT', 'TEXT'),
    ('contact_email', 'Email', 'Contact email address', 'CONTACT', 'EMAIL'),
    ('contact_company_name', 'Company Name', 'Contact company name', 'CONTACT', 'TEXT'),
    ('contact_job_title', 'Job Title', 'Contact job title', 'CONTACT', 'TEXT'),

    ('sender_first_name', 'Sender First Name', 'Sender first name', 'SENDER', 'TEXT'),
    ('sender_last_name', 'Sender Last Name', 'Sender last name', 'SENDER', 'TEXT'),
    ('sender_email', 'Sender Email', 'Sender email address', 'SENDER', 'EMAIL')

ON CONFLICT (key) DO NOTHING;


-- No-op. Master data is intentionally retained.