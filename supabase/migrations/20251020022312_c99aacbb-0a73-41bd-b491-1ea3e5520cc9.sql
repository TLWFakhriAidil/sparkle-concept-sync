-- Create enum types
CREATE TYPE public.app_role AS ENUM ('admin', 'user');
CREATE TYPE public.execution_status AS ENUM ('active', 'completed', 'failed');
CREATE TYPE public.message_sender AS ENUM ('user', 'bot', 'staff');
CREATE TYPE public.message_type AS ENUM ('text', 'image', 'document', 'audio', 'video');
CREATE TYPE public.order_status AS ENUM ('Pending', 'Processing', 'Success', 'Failed');
CREATE TYPE public.provider_type AS ENUM ('whacenter', 'wablas', 'waha');
CREATE TYPE public.ai_model AS ENUM ('openai/gpt-5-chat', 'openai/gpt-5-mini', 'openai/chatgpt-4o-latest', 'openai/gpt-4.1', 'google/gemini-2.5-pro', 'google/gemini-pro-1.5');

-- User profiles table
CREATE TABLE public.profiles (
  id UUID PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
  email VARCHAR(255) NOT NULL,
  full_name VARCHAR(255) NOT NULL,
  is_active BOOLEAN DEFAULT true,
  last_login TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User roles table
CREATE TABLE public.user_roles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE NOT NULL,
  role app_role NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  UNIQUE (user_id, role)
);

-- Device settings table
CREATE TABLE public.device_settings (
  id VARCHAR(255) PRIMARY KEY,
  device_id VARCHAR(255),
  api_key_option ai_model DEFAULT 'openai/gpt-4.1',
  webhook_id VARCHAR(500),
  provider provider_type DEFAULT 'wablas',
  phone_number VARCHAR(20),
  api_key TEXT,
  id_device VARCHAR(255),
  user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
  instance TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Chatbot flows table
CREATE TABLE public.chatbot_flows (
  id VARCHAR(255) PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  niche TEXT,
  id_device VARCHAR(255),
  nodes JSONB,
  edges JSONB,
  user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- AI WhatsApp conversations table
CREATE TABLE public.ai_whatsapp (
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
  execution_status execution_status DEFAULT NULL,
  flow_id VARCHAR(255) DEFAULT NULL,
  current_node_id VARCHAR(255) DEFAULT NULL,
  last_node_id VARCHAR(255) DEFAULT NULL,
  waiting_for_reply BOOLEAN DEFAULT false,
  human INTEGER DEFAULT 0,
  user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Conversation log table
CREATE TABLE public.conversation_log (
  id VARCHAR(255) PRIMARY KEY,
  prospect_num VARCHAR(20) NOT NULL,
  sender message_sender NOT NULL,
  message TEXT NOT NULL,
  message_type message_type DEFAULT 'text',
  stage VARCHAR(255),
  ai_response JSONB,
  device_id VARCHAR(255),
  user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Orders table
CREATE TABLE public.orders (
  id SERIAL PRIMARY KEY,
  amount DECIMAL(10,2) NOT NULL,
  collection_id VARCHAR(255),
  status order_status DEFAULT 'Pending',
  bill_id VARCHAR(255),
  url TEXT,
  product VARCHAR(255) NOT NULL,
  method VARCHAR(50) DEFAULT 'billplz',
  user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- WasapBot table
CREATE TABLE public.wasapbot (
  id_prospect SERIAL PRIMARY KEY,
  flow_reference VARCHAR(255) DEFAULT NULL,
  execution_id VARCHAR(255) DEFAULT NULL,
  execution_status execution_status DEFAULT NULL,
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
  user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Enable RLS on all tables
ALTER TABLE public.profiles ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.user_roles ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.device_settings ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.chatbot_flows ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.ai_whatsapp ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.conversation_log ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.wasapbot ENABLE ROW LEVEL SECURITY;

-- Security definer function to check roles
CREATE OR REPLACE FUNCTION public.has_role(_user_id uuid, _role app_role)
RETURNS boolean
LANGUAGE sql
STABLE
SECURITY DEFINER
SET search_path = public
AS $$
  SELECT EXISTS (
    SELECT 1
    FROM public.user_roles
    WHERE user_id = _user_id AND role = _role
  )
$$;

-- RLS Policies for profiles
CREATE POLICY "Users can view their own profile" ON public.profiles
  FOR SELECT USING (auth.uid() = id);

CREATE POLICY "Users can update their own profile" ON public.profiles
  FOR UPDATE USING (auth.uid() = id);

CREATE POLICY "Users can insert their own profile" ON public.profiles
  FOR INSERT WITH CHECK (auth.uid() = id);

-- RLS Policies for user_roles
CREATE POLICY "Users can view their own roles" ON public.user_roles
  FOR SELECT USING (auth.uid() = user_id);

CREATE POLICY "Admins can view all roles" ON public.user_roles
  FOR SELECT USING (public.has_role(auth.uid(), 'admin'));

CREATE POLICY "Admins can insert roles" ON public.user_roles
  FOR INSERT WITH CHECK (public.has_role(auth.uid(), 'admin'));

CREATE POLICY "Admins can update roles" ON public.user_roles
  FOR UPDATE USING (public.has_role(auth.uid(), 'admin'));

CREATE POLICY "Admins can delete roles" ON public.user_roles
  FOR DELETE USING (public.has_role(auth.uid(), 'admin'));

-- RLS Policies for device_settings
CREATE POLICY "Users can view their own devices" ON public.device_settings
  FOR SELECT USING (auth.uid() = user_id);

CREATE POLICY "Users can create their own devices" ON public.device_settings
  FOR INSERT WITH CHECK (auth.uid() = user_id);

CREATE POLICY "Users can update their own devices" ON public.device_settings
  FOR UPDATE USING (auth.uid() = user_id);

CREATE POLICY "Users can delete their own devices" ON public.device_settings
  FOR DELETE USING (auth.uid() = user_id);

-- RLS Policies for chatbot_flows
CREATE POLICY "Users can view their own flows" ON public.chatbot_flows
  FOR SELECT USING (auth.uid() = user_id);

CREATE POLICY "Users can create their own flows" ON public.chatbot_flows
  FOR INSERT WITH CHECK (auth.uid() = user_id);

CREATE POLICY "Users can update their own flows" ON public.chatbot_flows
  FOR UPDATE USING (auth.uid() = user_id);

CREATE POLICY "Users can delete their own flows" ON public.chatbot_flows
  FOR DELETE USING (auth.uid() = user_id);

-- RLS Policies for ai_whatsapp
CREATE POLICY "Users can view their own conversations" ON public.ai_whatsapp
  FOR SELECT USING (auth.uid() = user_id);

CREATE POLICY "Users can create their own conversations" ON public.ai_whatsapp
  FOR INSERT WITH CHECK (auth.uid() = user_id);

CREATE POLICY "Users can update their own conversations" ON public.ai_whatsapp
  FOR UPDATE USING (auth.uid() = user_id);

CREATE POLICY "Users can delete their own conversations" ON public.ai_whatsapp
  FOR DELETE USING (auth.uid() = user_id);

-- RLS Policies for conversation_log
CREATE POLICY "Users can view their own conversation logs" ON public.conversation_log
  FOR SELECT USING (auth.uid() = user_id);

CREATE POLICY "Users can create their own conversation logs" ON public.conversation_log
  FOR INSERT WITH CHECK (auth.uid() = user_id);

-- RLS Policies for orders
CREATE POLICY "Users can view their own orders" ON public.orders
  FOR SELECT USING (auth.uid() = user_id);

CREATE POLICY "Users can create their own orders" ON public.orders
  FOR INSERT WITH CHECK (auth.uid() = user_id);

CREATE POLICY "Users can update their own orders" ON public.orders
  FOR UPDATE USING (auth.uid() = user_id);

-- RLS Policies for wasapbot
CREATE POLICY "Users can view their own wasapbot data" ON public.wasapbot
  FOR SELECT USING (auth.uid() = user_id);

CREATE POLICY "Users can create their own wasapbot data" ON public.wasapbot
  FOR INSERT WITH CHECK (auth.uid() = user_id);

CREATE POLICY "Users can update their own wasapbot data" ON public.wasapbot
  FOR UPDATE USING (auth.uid() = user_id);

-- Create indexes for performance
CREATE INDEX idx_device_settings_user_id ON public.device_settings(user_id);
CREATE INDEX idx_chatbot_flows_user_id ON public.chatbot_flows(user_id);
CREATE INDEX idx_chatbot_flows_id_device ON public.chatbot_flows(id_device);
CREATE INDEX idx_ai_whatsapp_user_id ON public.ai_whatsapp(user_id);
CREATE INDEX idx_ai_whatsapp_prospect_num ON public.ai_whatsapp(prospect_num);
CREATE INDEX idx_ai_whatsapp_execution_status ON public.ai_whatsapp(execution_status);
CREATE INDEX idx_conversation_log_prospect_num ON public.conversation_log(prospect_num);
CREATE INDEX idx_conversation_log_user_id ON public.conversation_log(user_id);
CREATE INDEX idx_orders_user_id ON public.orders(user_id);
CREATE INDEX idx_wasapbot_prospect_num ON public.wasapbot(prospect_num);
CREATE INDEX idx_wasapbot_user_id ON public.wasapbot(user_id);

-- Trigger to create profile on user signup
CREATE OR REPLACE FUNCTION public.handle_new_user()
RETURNS TRIGGER
LANGUAGE plpgsql
SECURITY DEFINER SET search_path = public
AS $$
BEGIN
  INSERT INTO public.profiles (id, email, full_name)
  VALUES (
    new.id,
    new.email,
    COALESCE(new.raw_user_meta_data->>'full_name', new.email)
  );
  
  -- Assign default 'user' role
  INSERT INTO public.user_roles (user_id, role)
  VALUES (new.id, 'user');
  
  RETURN new;
END;
$$;

CREATE TRIGGER on_auth_user_created
  AFTER INSERT ON auth.users
  FOR EACH ROW EXECUTE FUNCTION public.handle_new_user();

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION public.update_updated_at_column()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$;

-- Triggers for updated_at
CREATE TRIGGER update_profiles_updated_at
  BEFORE UPDATE ON public.profiles
  FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

CREATE TRIGGER update_device_settings_updated_at
  BEFORE UPDATE ON public.device_settings
  FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

CREATE TRIGGER update_chatbot_flows_updated_at
  BEFORE UPDATE ON public.chatbot_flows
  FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

CREATE TRIGGER update_ai_whatsapp_updated_at
  BEFORE UPDATE ON public.ai_whatsapp
  FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

CREATE TRIGGER update_orders_updated_at
  BEFORE UPDATE ON public.orders
  FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();