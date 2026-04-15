CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS postgis;

-- =========================
-- USERS
-- =========================
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,

    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- =========================
-- USER PROFILES
-- =========================
CREATE TABLE IF NOT EXISTS user_profiles(
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,

    username VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    avatar_url TEXT,

    trust_score FLOAT DEFAULT 0.5,

    total_posts INT DEFAULT 0,
    total_confirmations INT DEFAULT 0,
    total_reports INT DEFAULT 0,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- =========================
-- ENUMS
-- =========================
CREATE TYPE identity_type AS ENUM(
    'PUBLIC',
    'PSEUDONYMOUS',
    'ANONYMOUS'
);

CREATE TYPE post_type AS ENUM(
    'INCIDENT',
    'BROADCAST'
);

CREATE TYPE feedback_type AS ENUM (
    'HELPFUL',
    'MISLEADING'
);

-- =========================
-- INCIDENTS
-- =========================
CREATE TABLE IF NOT EXISTS incidents(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     
    location GEOGRAPHY(POINT, 4326) NOT NULL,

    title TEXT,
    category VARCHAR(50),

    post_count INT DEFAULT 0,
    confirmation_count INT DEFAULT 0,
    trust_score FLOAT DEFAULT 0.5,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_incidents_location ON incidents USING GIST(location);
CREATE INDEX idx_incidents_created_at ON incidents(created_at DESC);

-- =========================
-- POSTS 
-- =========================
CREATE TABLE IF NOT EXISTS posts(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID REFERENCES users(id) ON DELETE CASCADE,


    incident_id UUID REFERENCES incidents(id) ON DELETE CASCADE,

    content TEXT NOT NULL,

    location GEOGRAPHY(POINT, 4326) NOT NULL,

    radius INT DEFAULT 2000,

    identity_type identity_type NOT NULL,


    post_type post_type NOT NULL,

    trust_score FLOAT,

    media_urls TEXT[],

    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ DEFAULT (NOW() + INTERVAL '7 days'),

    is_deleted BOOLEAN DEFAULT FALSE,
    is_flagged BOOLEAN DEFAULT FALSE,

    CONSTRAINT post_incident_consistency CHECK (
        (post_type = 'INCIDENT' AND incident_id IS NOT NULL)
        OR
        (post_type = 'BROADCAST' AND incident_id IS NULL)
    )
);

CREATE INDEX idx_posts_location ON posts USING GIST(location);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_incident_id ON posts(incident_id);

-- =========================
-- POST FEEDBACK
-- =========================
CREATE TABLE post_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,

    feedback_type feedback_type,

    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(post_id, user_id)
);

-- =========================
-- INCIDENT CONFIRMATIONS
-- =========================
CREATE TABLE incident_confirmations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    incident_id UUID REFERENCES incidents(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,

    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(incident_id, user_id)
);

-- =========================
-- POST REPORTS
-- =========================
CREATE TABLE post_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,

    reason TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);