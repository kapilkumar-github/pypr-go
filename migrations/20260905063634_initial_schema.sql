-- +goose Up
SELECT 'up SQL query';

-- ============================================================
-- ORGANIZATIONS
-- ============================================================

CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(255) NOT NULL,

    type VARCHAR(20) NOT NULL
        CHECK (type IN ('PERSONAL', 'BUSINESS')),

    timezone VARCHAR(50) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- ============================================================
-- USERS
-- ============================================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    email VARCHAR(320) NOT NULL UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,

    timezone VARCHAR(50) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX idx_users_email
    ON users(email);


-- ============================================================
-- ORGANIZATION MEMBERS
-- ============================================================

CREATE TABLE organization_members (
    organization_id UUID NOT NULL
        REFERENCES organizations(id)
        ON DELETE CASCADE,

    user_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    role VARCHAR(50) NOT NULL DEFAULT 'MEMBER',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (organization_id, user_id)
);


CREATE INDEX idx_organization_members_user_id
    ON organization_members(user_id);


-- ============================================================
-- ORGANIZATION INVITATIONS
-- ============================================================

CREATE TABLE organization_invites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    organization_id UUID NOT NULL
        REFERENCES organizations(id)
        ON DELETE CASCADE,

    email VARCHAR(320) NOT NULL,

    role VARCHAR(50) NOT NULL DEFAULT 'MEMBER',

    token_encrypted TEXT NOT NULL,

    token_hash TEXT NOT NULL UNIQUE,

    expires_at TIMESTAMPTZ NOT NULL
        DEFAULT (NOW() + INTERVAL '7 days'),

    accepted_at TIMESTAMPTZ,

    created_by UUID NOT NULL
        REFERENCES users(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX idx_organization_invites_email
    ON organization_invites(email);

CREATE INDEX idx_organization_invites_organization_id
    ON organization_invites(organization_id);

CREATE UNIQUE INDEX uq_org_invitation_active_email
    ON organization_invites (organization_id, email)
    WHERE accepted_at IS NULL;


-- ============================================================
-- USER CREDENTIALS
-- ============================================================

CREATE TABLE user_credentials (
    user_id UUID PRIMARY KEY
        REFERENCES users(id)
        ON DELETE CASCADE,

    password_hash TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- ============================================================
-- CONTACTS
--
-- Contacts are organization-wide.
-- Every member can see contacts in the organization.
-- Each contact has one owner.
--
-- Example:
--
-- Organization A
--     John Doe
--     owner_user_id = Alice
--
-- Bob can SEE John,
-- but Bob cannot start a sequence for John.
-- ============================================================

CREATE TABLE contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    organization_id UUID NOT NULL
        REFERENCES organizations(id)
        ON DELETE CASCADE,

    owner UUID NOT NULL
        REFERENCES users(id),

    email VARCHAR(320) NOT NULL,

    first_name VARCHAR(100),
    last_name VARCHAR(100),

    phone VARCHAR(30),

    company VARCHAR(255),
    job_title VARCHAR(255),

    timezone VARCHAR(50),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_contacts_organization_email
        UNIQUE (organization_id, email)
);


CREATE INDEX idx_contacts_organization_id
    ON contacts(organization_id);

CREATE INDEX idx_contacts_owner
    ON contacts(owner);


CREATE TABLE variables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    key TEXT NOT NULL UNIQUE,

    label TEXT NOT NULL,

    description TEXT,

    source TEXT NOT NULL
        CHECK (source IN (
            'CONTACT',
            'SENDER'
        )),

    data_type TEXT NOT NULL DEFAULT 'TEXT'
        CHECK (data_type IN (
            'TEXT',
            'EMAIL',
            'URL'
        )),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- ============================================================
-- SEQUENCES
-- ============================================================

CREATE TABLE sequences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    organization_id UUID NOT NULL
        REFERENCES organizations(id)
        ON DELETE CASCADE,

    owner UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    name TEXT NOT NULL,

    description TEXT,

    timezone_mode VARCHAR(30) NOT NULL DEFAULT 'OWNER',

    timezone VARCHAR(50),

    status TEXT NOT NULL DEFAULT 'DRAFT'
        CHECK (status IN (
            'DRAFT',
            'ACTIVE',
            'PAUSED',
            'ARCHIVED'
        )),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_sequence_timezone_mode
        CHECK (
            timezone_mode IN (
                'OWNER',
                'CONTACT',
                'CUSTOM'
            )
        ),

    CONSTRAINT chk_sequence_custom_timezone
        CHECK (
            (
                timezone_mode = 'CUSTOM'
                AND timezone IS NOT NULL
            )
            OR
            (
                timezone_mode <> 'CUSTOM'
            )
        )
);


CREATE INDEX idx_sequences_owner
    ON sequences(owner);


-- ============================================================
-- SEQUENCE STEPS
-- ============================================================
CREATE TABLE sequence_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    sequence_id UUID NOT NULL
        REFERENCES sequences(id) ON DELETE CASCADE,

    name TEXT NOT NULL,

    step_order INTEGER NOT NULL,

    type TEXT NOT NULL
        CHECK (type IN ('EMAIL', 'WAIT')),

    -- WAIT step
    wait_seconds BIGINT,

    -- EMAIL scheduling
    schedule_type TEXT
        CHECK (schedule_type IN (
            'IMMEDIATE',
            'WEEKDAY_TIME',
            'EXACT_DATE'
        )),

    scheduled_day SMALLINT,
    scheduled_time TIME,
    scheduled_at TIMESTAMPTZ,

    -- Email content
    subject TEXT,
    body TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (sequence_id, step_order),

    CHECK (
        -- WAIT step
        (
            type = 'WAIT'
            AND wait_seconds IS NOT NULL
            AND wait_seconds > 0
            AND schedule_type IS NULL
            AND scheduled_day IS NULL
            AND scheduled_time IS NULL
            AND scheduled_at IS NULL
            AND subject IS NULL
            AND body IS NULL
        )

        OR

        -- EMAIL step
        (
            type = 'EMAIL'
            AND wait_seconds IS NULL
            AND subject IS NOT NULL
            AND body IS NOT NULL
            AND schedule_type IS NOT NULL

            AND (
                -- Send immediately
                (
                    schedule_type = 'IMMEDIATE'
                    AND scheduled_day IS NULL
                    AND scheduled_time IS NULL
                    AND scheduled_at IS NULL
                )

                OR

                -- Send on a specific weekday/time
                (
                    schedule_type = 'WEEKDAY_TIME'
                    AND scheduled_day BETWEEN 0 AND 6
                    AND scheduled_time IS NOT NULL
                    AND scheduled_at IS NULL
                )

                OR

                -- Send at an exact date/time
                (
                    schedule_type = 'EXACT_DATE'
                    AND scheduled_day IS NULL
                    AND scheduled_time IS NULL
                    AND scheduled_at IS NOT NULL
                )
            )
        )
    )
);

CREATE INDEX idx_sequence_steps_sequence_id
    ON sequence_steps(sequence_id);

CREATE TABLE sequence_step_variables (
    sequence_step_id UUID NOT NULL
        REFERENCES sequence_steps(id) ON DELETE CASCADE,

    variable_id UUID NOT NULL
        REFERENCES variables(id) ON DELETE RESTRICT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (sequence_step_id, variable_id)
);

CREATE INDEX idx_sequence_step_variables_variable_id
    ON sequence_step_variables(variable_id);

-- ============================================================
-- SEQUENCE ENROLMENTS
-- ============================================================

CREATE TABLE sequence_enrolments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    sequence_id UUID NOT NULL
        REFERENCES sequences(id)
        ON DELETE CASCADE,

    contact_id UUID NOT NULL
        REFERENCES contacts(id)
        ON DELETE CASCADE,

    status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',

    current_step_order INTEGER,

    enrolled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    completed_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_sequence_enrolment
        UNIQUE (sequence_id, contact_id),

    CONSTRAINT chk_sequence_enrolment_status
        CHECK (
            status IN (
                'ACTIVE',
                'COMPLETED',
                'PAUSED',
                'CANCELLED'
            )
        )
);


CREATE INDEX idx_sequence_enrolment_sequence
    ON sequence_enrolments(sequence_id);

CREATE INDEX idx_sequence_enrolment_contact
    ON sequence_enrolments(contact_id);


-- ============================================================
-- SEQUENCE STEP EXECUTIONS
-- ============================================================

CREATE TABLE sequence_step_executions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    enrolment_id UUID NOT NULL
        REFERENCES sequence_enrolments(id)
        ON DELETE CASCADE,

    sequence_step_id UUID NOT NULL
        REFERENCES sequence_steps(id)
        ON DELETE CASCADE,

    step_order INTEGER NOT NULL,

    scheduled_at TIMESTAMPTZ NOT NULL,

    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',

    attempt_count INTEGER NOT NULL DEFAULT 0,

    locked_at TIMESTAMPTZ,

    started_at TIMESTAMPTZ,

    completed_at TIMESTAMPTZ,

    last_error TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_step_execution_status
        CHECK (
            status IN (
                'PENDING',
                'PROCESSING',
                'COMPLETED',
                'FAILED',
                'SKIPPED',
                'CANCELLED'
            )
        )
);


-- Used by the polling worker:
--
-- WHERE status = 'PENDING'
-- AND scheduled_at <= NOW()
--
CREATE INDEX idx_step_execution_pending_schedule
    ON sequence_step_executions (scheduled_at)
    WHERE status = 'PENDING';


CREATE INDEX idx_step_execution_enrolment
    ON sequence_step_executions(enrolment_id);

-- +goose Down
SELECT 'down SQL query';

DROP TABLE IF EXISTS sequence_step_variables;
DROP TABLE IF EXISTS sequence_step_executions;
DROP TABLE IF EXISTS sequence_enrolments;
DROP TABLE IF EXISTS sequence_steps;
DROP TABLE IF EXISTS sequences;

DROP TABLE IF EXISTS variables;

DROP TABLE IF EXISTS contacts;

DROP TABLE IF EXISTS user_credentials;
DROP TABLE IF EXISTS organization_invites;
DROP TABLE IF EXISTS organization_members;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS organizations;