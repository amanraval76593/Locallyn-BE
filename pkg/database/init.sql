CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,

    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);


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

CREATE TYPE identity_type AS ENUM(
    'PUBLIC',
    'PSEUDONYMOUS',
    'ANONYMOUS'
);
CREATE TABLE IF NOT EXISTS incidents(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     
    location GEOGRAPHY(POINT, 4326) NOT NULL,

    title TEXT,
    category VARCHAR(50),

    -- aggregation stats
    post_count INT DEFAULT 1,
    confirmation_count INT DEFAULT 0,
    -- computed trust of incident
    trust_score FLOAT DEFAULT 0.5,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    -- lifecycle
    expires_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS posts(
    id UUID PRIMARY KEY default gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    incident_id UUID REFERENCES incidents(id) ON DELETE CASCADE,
    content TEXT NOT NULL,

     -- GEO LOCATION (PostGIS)
    location GEOGRAPHY(POINT, 4326) NOT NULL,

    -- visibility radius (meters)
    radius INT DEFAULT 2000,

    identity_type identity_type NOT NULL,

    -- snapshot trust at posting time
    trust_score FLOAT,

    -- metadata
    media_urls TEXT[],

    -- lifecycle
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ,

    -- moderation flags
    is_deleted BOOLEAN DEFAULT FALSE,
    is_flagged BOOLEAN DEFAULT FALSE
);


CREATE INDEX idx_posts_location ON posts USING GIST(location);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
CREATE INDEX idx_posts_user_id ON posts(user_id);


CREATE TABLE post_incidents (
    post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
    incident_id UUID REFERENCES incidents(id) ON DELETE CASCADE,

    PRIMARY KEY (post_id, incident_id)
);

CREATE INDEX idx_incidents_location ON incidents USING GIST(location);
CREATE INDEX idx_incidents_created_at ON incidents(created_at DESC);

CREATE TYPE feedback_type AS ENUM (
    'HELPFUL',
    'MISLEADING'
);

CREATE TABLE post_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,

    feedback_type feedback_type,

    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(post_id, user_id)
);

CREATE TABLE incident_confirmations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    incident_id UUID REFERENCES incidents(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,

    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(incident_id, user_id)
);

CREATE TABLE post_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id),

    reason TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
