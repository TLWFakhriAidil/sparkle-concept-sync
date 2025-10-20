package models

import (
	"encoding/json"
	"time"
)

// User represents a user in the system
type User struct {
	ID           string     `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	FullName     string     `json:"full_name" db:"full_name"`
	PasswordHash string     `json:"-" db:"password_hash"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin    *time.Time `json:"last_login" db:"last_login"`
}

// UserSession represents a user session for JWT authentication
type UserSession struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// DeviceSetting represents WhatsApp device configuration
type DeviceSetting struct {
	ID           string    `json:"id" db:"id"`
	DeviceID     *string   `json:"device_id" db:"device_id"`
	APIKeyOption string    `json:"api_key_option" db:"api_key_option"`
	WebhookID    *string   `json:"webhook_id" db:"webhook_id"`
	Provider     string    `json:"provider" db:"provider"`
	PhoneNumber  *string   `json:"phone_number" db:"phone_number"`
	APIKey       *string   `json:"api_key" db:"api_key"`
	IDDevice     *string   `json:"id_device" db:"id_device"`
	UserID       *string   `json:"user_id" db:"user_id"`
	Instance     *string   `json:"instance" db:"instance"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// FlowNode represents a node in the chatbot flow
type FlowNode struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Position map[string]interface{} `json:"position"`
	Data     map[string]interface{} `json:"data"`
}

// FlowEdge represents an edge in the chatbot flow
type FlowEdge struct {
	ID           string  `json:"id"`
	Source       string  `json:"source"`
	Target       string  `json:"target"`
	SourceHandle *string `json:"sourceHandle,omitempty"`
	TargetHandle *string `json:"targetHandle,omitempty"`
	Label        *string `json:"label,omitempty"`
}

// ChatbotFlow represents a complete chatbot flow
type ChatbotFlow struct {
	ID          string     `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description *string    `json:"description" db:"description"`
	Niche       *string    `json:"niche" db:"niche"`
	IDDevice    *string    `json:"id_device" db:"id_device"`
	Nodes       []FlowNode `json:"nodes" db:"nodes"`
	Edges       []FlowEdge `json:"edges" db:"edges"`
	UserID      *string    `json:"user_id" db:"user_id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// AIWhatsApp represents an AI WhatsApp conversation
type AIWhatsApp struct {
	IDProspect      int        `json:"id_prospect" db:"id_prospect"`
	FlowReference   *string    `json:"flow_reference" db:"flow_reference"`
	ExecutionID     *string    `json:"execution_id" db:"execution_id"`
	DateOrder       *time.Time `json:"date_order" db:"date_order"`
	IDDevice        *string    `json:"id_device" db:"id_device"`
	Niche           *string    `json:"niche" db:"niche"`
	ProspectName    *string    `json:"prospect_name" db:"prospect_name"`
	ProspectNum     *string    `json:"prospect_num" db:"prospect_num"`
	Stage           *string    `json:"stage" db:"stage"`
	ConvLast        *string    `json:"conv_last" db:"conv_last"`
	ConvCurrent     *string    `json:"conv_current" db:"conv_current"`
	ExecutionStatus *string    `json:"execution_status" db:"execution_status"`
	FlowID          *string    `json:"flow_id" db:"flow_id"`
	CurrentNodeID   *string    `json:"current_node_id" db:"current_node_id"`
	LastNodeID      *string    `json:"last_node_id" db:"last_node_id"`
	WaitingForReply *bool      `json:"waiting_for_reply" db:"waiting_for_reply"`
	Human           *int       `json:"human" db:"human"`
	UserID          *string    `json:"user_id" db:"user_id"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// ConversationLog represents a logged conversation message
type ConversationLog struct {
	ID          string                 `json:"id" db:"id"`
	ProspectNum string                 `json:"prospect_num" db:"prospect_num"`
	Sender      string                 `json:"sender" db:"sender"`
	Message     string                 `json:"message" db:"message"`
	MessageType *string                `json:"message_type" db:"message_type"`
	Stage       *string                `json:"stage" db:"stage"`
	AIResponse  map[string]interface{} `json:"ai_response" db:"ai_response"`
	DeviceID    *string                `json:"device_id" db:"device_id"`
	UserID      *string                `json:"user_id" db:"user_id"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
}

// Order represents a billing order
type Order struct {
	ID           int       `json:"id" db:"id"`
	Amount       float64   `json:"amount" db:"amount"`
	CollectionID *string   `json:"collection_id" db:"collection_id"`
	Status       *string   `json:"status" db:"status"`
	BillID       *string   `json:"bill_id" db:"bill_id"`
	URL          *string   `json:"url" db:"url"`
	Product      string    `json:"product" db:"product"`
	Method       *string   `json:"method" db:"method"`
	UserID       *string   `json:"user_id" db:"user_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// WasapBot represents WasapBot flow processing data
type WasapBot struct {
	IDProspect      int       `json:"id_prospect" db:"id_prospect"`
	FlowReference   *string   `json:"flow_reference" db:"flow_reference"`
	ExecutionID     *string   `json:"execution_id" db:"execution_id"`
	ExecutionStatus *string   `json:"execution_status" db:"execution_status"`
	FlowID          *string   `json:"flow_id" db:"flow_id"`
	CurrentNodeID   *string   `json:"current_node_id" db:"current_node_id"`
	LastNodeID      *string   `json:"last_node_id" db:"last_node_id"`
	WaitingForReply *bool     `json:"waiting_for_reply" db:"waiting_for_reply"`
	ProspectNum     *string   `json:"prospect_num" db:"prospect_num"`
	Niche           *string   `json:"niche" db:"niche"`
	Instance        *string   `json:"instance" db:"instance"`
	Nama            *string   `json:"nama" db:"nama"`
	Stage           *string   `json:"stage" db:"stage"`
	ConvLast        *string   `json:"conv_last" db:"conv_last"`
	Status          *string   `json:"status" db:"status"`
	UserInput       *string   `json:"user_input" db:"user_input"`
	UserID          *string   `json:"user_id" db:"user_id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// WhatsAppMessage represents an incoming WhatsApp message
type WhatsAppMessage struct {
	From       string                 `json:"from"`
	To         string                 `json:"to"`
	Body       string                 `json:"body"`
	Type       string                 `json:"type"`
	Timestamp  int64                  `json:"timestamp"`
	DeviceID   string                 `json:"device_id"`
	InstanceID string                 `json:"instance_id"`
	MessageID  string                 `json:"message_id"`
	MediaURL   string                 `json:"media_url,omitempty"`
	Caption    string                 `json:"caption,omitempty"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

// AIResponse represents an AI service response
type AIResponse struct {
	Stage    string      `json:"Stage"`
	Response []AIMessage `json:"Response"`
}

// AIMessage represents a single AI message
type AIMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Type      string                 `json:"type"`
	UserID    string                 `json:"user_id"`
	DeviceID  string                 `json:"device_id,omitempty"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// Analytics represents analytics data
type Analytics struct {
	TotalConversations     int `json:"total_conversations"`
	ActiveConversations    int `json:"active_conversations"`
	CompletedConversations int `json:"completed_conversations"`
	FailedConversations    int `json:"failed_conversations"`
	TotalMessages          int `json:"total_messages"`
	TotalDevices           int `json:"total_devices"`
	TotalFlows             int `json:"total_flows"`
}

// ExecutionProcess represents flow execution state
type ExecutionProcess struct {
	ExecutionID     string                 `json:"execution_id"`
	FlowID          string                 `json:"flow_id"`
	CurrentNodeID   string                 `json:"current_node_id"`
	ProspectNum     string                 `json:"prospect_num"`
	Variables       map[string]interface{} `json:"variables"`
	WaitingForReply bool                   `json:"waiting_for_reply"`
	Status          string                 `json:"status"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// Helper methods for JSON marshaling/unmarshaling

func (cf *ChatbotFlow) MarshalNodes() ([]byte, error) {
	return json.Marshal(cf.Nodes)
}

func (cf *ChatbotFlow) UnmarshalNodes(data []byte) error {
	return json.Unmarshal(data, &cf.Nodes)
}

func (cf *ChatbotFlow) MarshalEdges() ([]byte, error) {
	return json.Marshal(cf.Edges)
}

func (cf *ChatbotFlow) UnmarshalEdges(data []byte) error {
	return json.Unmarshal(data, &cf.Edges)
}

func (cl *ConversationLog) MarshalAIResponse() ([]byte, error) {
	return json.Marshal(cl.AIResponse)
}

func (cl *ConversationLog) UnmarshalAIResponse(data []byte) error {
	return json.Unmarshal(data, &cl.AIResponse)
}
