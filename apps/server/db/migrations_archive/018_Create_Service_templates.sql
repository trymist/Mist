-- Service Templates for Pre-made Services (Database / Service Type)
-- Phase 2: Database Services from roadmap
CREATE TABLE IF NOT EXISTS service_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE, -- 'postgres', 'redis', 'mysql', 'mongodb', etc.
    display_name TEXT NOT NULL, -- 'PostgreSQL 16', 'Redis 7', etc.
    category TEXT NOT NULL CHECK(category IN ('database', 'cache', 'queue', 'storage', 'other')) DEFAULT 'database',
    description TEXT,
    icon_url TEXT,
    
    -- Docker Configuration
    docker_image TEXT NOT NULL, -- e.g., 'postgres:16-alpine'
    docker_image_version TEXT, -- For tracking specific versions
    
    -- Default Configuration
    default_port INTEGER NOT NULL,
    default_env_vars TEXT, -- JSON object of default environment variables
    required_env_vars TEXT, -- JSON array of required env var keys
    
    -- Volume Configuration
    default_volume_path TEXT, -- Container path for data persistence
    volume_required BOOLEAN DEFAULT 1,
    
    -- Resource Recommendations
    recommended_cpu REAL, -- Recommended CPU limit
    recommended_memory INTEGER, -- Recommended memory in MB
    min_memory INTEGER, -- Minimum memory required
    
    -- Health Check
    healthcheck_command TEXT,
    healthcheck_interval INTEGER DEFAULT 30,
    
    -- UI Configuration
    admin_ui_image TEXT, -- Optional admin UI (e.g., 'phpmyadmin', 'pgadmin')
    admin_ui_port INTEGER,
    setup_instructions TEXT, -- Markdown text with setup guide
    
    -- Metadata
    is_active BOOLEAN DEFAULT 1,
    is_featured BOOLEAN DEFAULT 0,
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Insert default templates
INSERT INTO service_templates (name, display_name, category, description, docker_image, default_port, default_env_vars, default_volume_path, recommended_memory, min_memory)
VALUES 
    ('postgres', 'PostgreSQL 16', 'database', 'PostgreSQL is a powerful, open source object-relational database system', 'postgres:16-alpine', 5432, '{"POSTGRES_PASSWORD":"GENERATE","POSTGRES_DB":"myapp","POSTGRES_USER":"postgres"}', '/var/lib/postgresql/data', 512, 256),
    ('redis', 'Redis 7', 'cache', 'Redis is an in-memory data structure store, used as a database, cache, and message broker', 'redis:7-alpine', 6379, '{}', '/data', 256, 128),
    ('mysql', 'MySQL 8', 'database', 'MySQL is the world''s most popular open source database', 'mysql:8', 3306, '{"MYSQL_ROOT_PASSWORD":"GENERATE","MYSQL_DATABASE":"myapp"}', '/var/lib/mysql', 512, 256),
    ('mariadb', 'MariaDB 11', 'database', 'MariaDB is a community-developed fork of MySQL', 'mariadb:11', 3306, '{"MARIADB_ROOT_PASSWORD":"GENERATE","MARIADB_DATABASE":"myapp"}', '/var/lib/mysql', 512, 256),
    ('mongodb', 'MongoDB 7', 'database', 'MongoDB is a source-available cross-platform document-oriented database', 'mongo:7', 27017, '{"MONGO_INITDB_ROOT_USERNAME":"admin","MONGO_INITDB_ROOT_PASSWORD":"GENERATE"}', '/data/db', 512, 256),
    ('rabbitmq', 'RabbitMQ 3', 'queue', 'RabbitMQ is a reliable and mature messaging and streaming broker', 'rabbitmq:3-management', 5672, '{"RABBITMQ_DEFAULT_USER":"admin","RABBITMQ_DEFAULT_PASS":"GENERATE"}', '/var/lib/rabbitmq', 512, 256),
    ('minio', 'MinIO', 'storage', 'MinIO is a high-performance, S3 compatible object store', 'minio/minio', 9000, '{"MINIO_ROOT_USER":"admin","MINIO_ROOT_PASSWORD":"GENERATE"}', '/data', 512, 256);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_service_templates_category ON service_templates(category);
CREATE INDEX IF NOT EXISTS idx_service_templates_is_active ON service_templates(is_active);
CREATE INDEX IF NOT EXISTS idx_service_templates_sort_order ON service_templates(sort_order);
