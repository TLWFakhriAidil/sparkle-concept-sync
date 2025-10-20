package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// Connect establishes a connection to PostgreSQL with connection pooling
func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Configure connection pooling for 5000+ concurrent users
	db.SetMaxOpenConns(500)                // Maximum number of open connections
	db.SetMaxIdleConns(100)                // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum amount of time a connection may be reused

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	log.Println("✅ Database connected successfully with connection pooling")
	return db, nil
}

// RunMigrations executes database migrations
func RunMigrations(db *sql.DB) error {
	log.Println("🔄 Running database migrations...")

	migrations := []string{
		createUsersTable,
		createUserSessionsTable,
		createDeviceSettingTable,
		createChatbotFlowsTable,
		createAIWhatsAppTable,
		createConversationLogTable,
		createOrdersTable,
		createWasapBotTable,
		createIndexes,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration %d failed: %v", i+1, err)
		}
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id CHAR(36) NOT NULL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login TIMESTAMP WITH TIME ZONE DEFAULT NULL
);`

const createUserSessionsTable = `
CREATE TABLE IF NOT EXISTS user_sessions (
    id CHAR(36) NOT NULL PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);`

const createDeviceSettingTable = `
CREATE TABLE IF NOT EXISTS device_setting (
    id VARCHAR(255) PRIMARY KEY,
    device_id VARCHAR(255),
    api_key_option VARCHAR(100) DEFAULT 'openai/gpt-4.1' CHECK (api_key_option IN (
        'openai/gpt-5-chat', 
        'openai/gpt-5-mini', 
        'openai/chatgpt-4o-latest', 
        'openai/gpt-4.1', 
        'google/gemini-2.5-pro', 
        'google/gemini-pro-1.5'
    )),
    webhook_id VARCHAR(500),
    provider VARCHAR(20) DEFAULT 'wablas' CHECK (provider IN ('whacenter', 'wablas', 'waha')),
    phone_number VARCHAR(20),
    api_key TEXT,
    id_device VARCHAR(255),
    user_id CHAR(36),
    instance TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);`

const createChatbotFlowsTable = `
CREATE TABLE IF NOT EXISTS chatbot_flows (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    niche TEXT,
    id_device VARCHAR(255),
    nodes JSONB,
    edges JSONB,
    user_id CHAR(36),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);`

const createAIWhatsAppTable = `
CREATE TABLE IF NOT EXISTS ai_whatsapp (
    id_prospect SERIAL PRIMARY KEY,
    flow_reference VARCHAR(255) DEFAULT NULL,
    execution_id VARCHAR(255) DEFAULT NULL,
    date_order TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    id_device VARCHAR(255) DEFAULT NULL,
    niche VARCHAR(255) DEFAULT NULL,
    prospect_name VARCHAR(255) DEFAULT NULL,
    prospect_num VARCHAR(255) DEFAULT NULL,
    stage VARCHAR(255) DEFAULT NULL,
    conv_last TEXT,
    conv_current TEXT,
    execution_status VARCHAR(20) DEFAULT NULL CHECK (execution_status IN ('active','completed','failed')),
    flow_id VARCHAR(255) DEFAULT NULL,
    current_node_id VARCHAR(255) DEFAULT NULL,
    last_node_id VARCHAR(255) DEFAULT NULL,
    waiting_for_reply BOOLEAN DEFAULT false,
    human INTEGER DEFAULT 0,
    user_id CHAR(36),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);`

const createConversationLogTable = `
CREATE TABLE IF NOT EXISTS conversation_log (
    id VARCHAR(255) PRIMARY KEY,
    prospect_num VARCHAR(20) NOT NULL,
    sender VARCHAR(10) NOT NULL CHECK (sender IN ('user', 'bot', 'staff')),
    message TEXT NOT NULL,
    message_type VARCHAR(10) DEFAULT 'text' CHECK (message_type IN ('text', 'image', 'document', 'audio', 'video')),
    stage VARCHAR(255),
    ai_response JSONB,
    device_id VARCHAR(255),
    user_id CHAR(36),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);`

const createOrdersTable = `
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    amount DECIMAL(10,2) NOT NULL,
    collection_id VARCHAR(255),
    status VARCHAR(20) DEFAULT 'Pending' CHECK (status IN ('Pending', 'Processing', 'Success', 'Failed')),
    bill_id VARCHAR(255),
    url TEXT,
    product VARCHAR(255) NOT NULL,
    method VARCHAR(50) DEFAULT 'billplz',
    user_id CHAR(36),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);`

const createWasapBotTable = `
CREATE TABLE IF NOT EXISTS wasapBot (
    id_prospect SERIAL PRIMARY KEY,
    flow_reference VARCHAR(255) DEFAULT NULL,
    execution_id VARCHAR(255) DEFAULT NULL,
    execution_status VARCHAR(20) DEFAULT NULL CHECK (execution_status IN ('active','completed','failed')),
    flow_id VARCHAR(255) DEFAULT NULL,
    current_node_id VARCHAR(255) DEFAULT NULL,
    last_node_id VARCHAR(255) DEFAULT NULL,
    waiting_for_reply BOOLEAN DEFAULT false,
    prospect_num VARCHAR(100) DEFAULT NULL,
    niche VARCHAR(300) DEFAULT NULL,
    instance VARCHAR(255) DEFAULT NULL,
    nama VARCHAR(100) DEFAULT NULL,
    stage VARCHAR(200) DEFAULT NULL,
    conv_last TEXT,
    status VARCHAR(200) DEFAULT 'Prospek',
    user_input TEXT,
    user_id CHAR(36),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);`

const createIndexes = `
-- Performance indexes for 5000+ concurrent users
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_device_setting_user_id ON device_setting(user_id);
CREATE INDEX IF NOT EXISTS idx_device_setting_id_device ON device_setting(id_device);
CREATE INDEX IF NOT EXISTS idx_chatbot_flows_user_id ON chatbot_flows(user_id);
CREATE INDEX IF NOT EXISTS idx_chatbot_flows_id_device ON chatbot_flows(id_device);
CREATE INDEX IF NOT EXISTS idx_ai_whatsapp_user_id ON ai_whatsapp(user_id);
CREATE INDEX IF NOT EXISTS idx_ai_whatsapp_prospect_num ON ai_whatsapp(prospect_num);
CREATE INDEX IF NOT EXISTS idx_ai_whatsapp_execution_status ON ai_whatsapp(execution_status);
CREATE INDEX IF NOT EXISTS idx_conversation_log_prospect_num ON conversation_log(prospect_num);
CREATE INDEX IF NOT EXISTS idx_conversation_log_user_id ON conversation_log(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_wasapbot_prospect_num ON wasapBot(prospect_num);
CREATE INDEX IF NOT EXISTS idx_wasapbot_user_id ON wasapBot(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_token ON user_sessions(token);
`
