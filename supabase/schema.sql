-- FocusFlow Database Schema

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Table: categories
-- User-defined categories for grouping activities
CREATE TABLE IF NOT EXISTS categories (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id      UUID REFERENCES auth.users(id) NOT NULL,
  name         TEXT NOT NULL,
  color        TEXT NOT NULL DEFAULT '#2E75B6',
  icon         TEXT NOT NULL DEFAULT 'circle',
  created_at   TIMESTAMPTZ DEFAULT now()
);

-- Table: activity_logs
-- Every event (window focus, URL visit, app open) is logged here
CREATE TABLE IF NOT EXISTS activity_logs (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id      UUID REFERENCES auth.users(id) NOT NULL,
  device_id    TEXT NOT NULL,
  app_name     TEXT NOT NULL,
  title        TEXT, -- Window title
  url          TEXT, -- URL (if browser)
  category_id  UUID REFERENCES categories(id),
  started_at   TIMESTAMPTZ NOT NULL,
  ended_at     TIMESTAMPTZ,
  duration_sec INTEGER GENERATED ALWAYS AS (
    EXTRACT(EPOCH FROM (ended_at - started_at))::INTEGER
  ) STORED,
  created_at   TIMESTAMPTZ DEFAULT now()
);

-- Table: app_category_mappings
-- Map apps/domains to categories automatically
CREATE TABLE IF NOT EXISTS app_category_mappings (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id      UUID REFERENCES auth.users(id) NOT NULL,
  matcher      TEXT NOT NULL, -- domain or app name
  match_type   TEXT NOT NULL CHECK (match_type IN ('domain', 'app', 'title')),
  category_id  UUID REFERENCES categories(id) NOT NULL,
  device_type  TEXT NOT NULL DEFAULT 'all', -- 'all', 'linux', 'android'
  created_at   TIMESTAMPTZ DEFAULT now()
);

-- Row Level Security (RLS) Policies
ALTER TABLE categories ENABLE ROW LEVEL SECURITY;
ALTER TABLE activity_logs ENABLE ROW LEVEL SECURITY;
ALTER TABLE app_category_mappings ENABLE ROW LEVEL SECURITY;

-- Categories: Users can only see/edit their own categories
CREATE POLICY "Users can manage their own categories" ON categories
  FOR ALL TO authenticated USING (auth.uid() = user_id);

-- Activity Logs: Users can only see/insert their own logs
CREATE POLICY "Users can manage their own activity logs" ON activity_logs
  FOR ALL TO authenticated USING (auth.uid() = user_id);

-- Mappings: Users can only see/edit their own mappings
CREATE POLICY "Users can manage their own mappings" ON app_category_mappings
  FOR ALL TO authenticated USING (auth.uid() = user_id);
