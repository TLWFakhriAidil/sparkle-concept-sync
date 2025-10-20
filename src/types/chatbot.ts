export type NodeType =
  | 'start'
  | 'message'
  | 'image'
  | 'audio'
  | 'video'
  | 'delay'
  | 'condition'
  | 'stage'
  | 'user_reply'
  | 'ai_prompt'
  | 'advanced_ai_prompt'
  | 'manual';

export interface FlowNode {
  id: string;
  type: NodeType;
  position: { x: number; y: number };
  data: {
    label?: string;
    message?: string;
    mediaUrl?: string;
    delay?: number;
    condition?: string;
    stage?: string;
    timeout?: number;
    prompt?: string;
    [key: string]: any;
  };
}

export interface FlowEdge {
  id: string;
  source: string;
  target: string;
  sourceHandle?: string;
  targetHandle?: string;
  label?: string;
}

export interface ChatbotFlow {
  id: string;
  name: string;
  description?: string;
  niche?: string;
  id_device?: string;
  nodes: FlowNode[];
  edges: FlowEdge[];
  user_id?: string;
  created_at?: string;
  updated_at?: string;
}

export interface DeviceSettings {
  id: string;
  device_id?: string;
  api_key_option: string;
  webhook_id?: string;
  provider: 'whacenter' | 'wablas' | 'waha';
  phone_number?: string;
  api_key?: string;
  id_device?: string;
  user_id?: string;
  instance?: string;
  created_at?: string;
  updated_at?: string;
}

export interface Conversation {
  id_prospect: number;
  flow_reference?: string;
  execution_id?: string;
  date_order?: string;
  id_device?: string;
  niche?: string;
  prospect_name?: string;
  prospect_num?: string;
  stage?: string;
  conv_last?: string;
  conv_current?: string;
  execution_status?: 'active' | 'completed' | 'failed';
  flow_id?: string;
  current_node_id?: string;
  last_node_id?: string;
  waiting_for_reply?: boolean;
  human?: number;
  user_id?: string;
  created_at?: string;
  updated_at?: string;
}
