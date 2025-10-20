export type Json =
  | string
  | number
  | boolean
  | null
  | { [key: string]: Json | undefined }
  | Json[]

export type Database = {
  // Allows to automatically instantiate createClient with right options
  // instead of createClient<Database, { PostgrestVersion: 'XX' }>(URL, KEY)
  __InternalSupabase: {
    PostgrestVersion: "13.0.5"
  }
  public: {
    Tables: {
      ai_whatsapp: {
        Row: {
          conv_current: string | null
          conv_last: string | null
          created_at: string | null
          current_node_id: string | null
          date_order: string | null
          execution_id: string | null
          execution_status:
            | Database["public"]["Enums"]["execution_status"]
            | null
          flow_id: string | null
          flow_reference: string | null
          human: number | null
          id_device: string | null
          id_prospect: number
          last_node_id: string | null
          niche: string | null
          prospect_name: string | null
          prospect_num: string | null
          stage: string | null
          updated_at: string | null
          user_id: string | null
          waiting_for_reply: boolean | null
        }
        Insert: {
          conv_current?: string | null
          conv_last?: string | null
          created_at?: string | null
          current_node_id?: string | null
          date_order?: string | null
          execution_id?: string | null
          execution_status?:
            | Database["public"]["Enums"]["execution_status"]
            | null
          flow_id?: string | null
          flow_reference?: string | null
          human?: number | null
          id_device?: string | null
          id_prospect?: number
          last_node_id?: string | null
          niche?: string | null
          prospect_name?: string | null
          prospect_num?: string | null
          stage?: string | null
          updated_at?: string | null
          user_id?: string | null
          waiting_for_reply?: boolean | null
        }
        Update: {
          conv_current?: string | null
          conv_last?: string | null
          created_at?: string | null
          current_node_id?: string | null
          date_order?: string | null
          execution_id?: string | null
          execution_status?:
            | Database["public"]["Enums"]["execution_status"]
            | null
          flow_id?: string | null
          flow_reference?: string | null
          human?: number | null
          id_device?: string | null
          id_prospect?: number
          last_node_id?: string | null
          niche?: string | null
          prospect_name?: string | null
          prospect_num?: string | null
          stage?: string | null
          updated_at?: string | null
          user_id?: string | null
          waiting_for_reply?: boolean | null
        }
        Relationships: []
      }
      chatbot_flows: {
        Row: {
          created_at: string | null
          description: string | null
          edges: Json | null
          id: string
          id_device: string | null
          name: string
          niche: string | null
          nodes: Json | null
          updated_at: string | null
          user_id: string | null
        }
        Insert: {
          created_at?: string | null
          description?: string | null
          edges?: Json | null
          id: string
          id_device?: string | null
          name: string
          niche?: string | null
          nodes?: Json | null
          updated_at?: string | null
          user_id?: string | null
        }
        Update: {
          created_at?: string | null
          description?: string | null
          edges?: Json | null
          id?: string
          id_device?: string | null
          name?: string
          niche?: string | null
          nodes?: Json | null
          updated_at?: string | null
          user_id?: string | null
        }
        Relationships: []
      }
      conversation_log: {
        Row: {
          ai_response: Json | null
          created_at: string | null
          device_id: string | null
          id: string
          message: string
          message_type: Database["public"]["Enums"]["message_type"] | null
          prospect_num: string
          sender: Database["public"]["Enums"]["message_sender"]
          stage: string | null
          user_id: string | null
        }
        Insert: {
          ai_response?: Json | null
          created_at?: string | null
          device_id?: string | null
          id: string
          message: string
          message_type?: Database["public"]["Enums"]["message_type"] | null
          prospect_num: string
          sender: Database["public"]["Enums"]["message_sender"]
          stage?: string | null
          user_id?: string | null
        }
        Update: {
          ai_response?: Json | null
          created_at?: string | null
          device_id?: string | null
          id?: string
          message?: string
          message_type?: Database["public"]["Enums"]["message_type"] | null
          prospect_num?: string
          sender?: Database["public"]["Enums"]["message_sender"]
          stage?: string | null
          user_id?: string | null
        }
        Relationships: []
      }
      device_settings: {
        Row: {
          api_key: string | null
          api_key_option: Database["public"]["Enums"]["ai_model"] | null
          created_at: string | null
          device_id: string | null
          id: string
          id_device: string | null
          instance: string | null
          phone_number: string | null
          provider: Database["public"]["Enums"]["provider_type"] | null
          updated_at: string | null
          user_id: string | null
          webhook_id: string | null
        }
        Insert: {
          api_key?: string | null
          api_key_option?: Database["public"]["Enums"]["ai_model"] | null
          created_at?: string | null
          device_id?: string | null
          id: string
          id_device?: string | null
          instance?: string | null
          phone_number?: string | null
          provider?: Database["public"]["Enums"]["provider_type"] | null
          updated_at?: string | null
          user_id?: string | null
          webhook_id?: string | null
        }
        Update: {
          api_key?: string | null
          api_key_option?: Database["public"]["Enums"]["ai_model"] | null
          created_at?: string | null
          device_id?: string | null
          id?: string
          id_device?: string | null
          instance?: string | null
          phone_number?: string | null
          provider?: Database["public"]["Enums"]["provider_type"] | null
          updated_at?: string | null
          user_id?: string | null
          webhook_id?: string | null
        }
        Relationships: []
      }
      orders: {
        Row: {
          amount: number
          bill_id: string | null
          collection_id: string | null
          created_at: string | null
          id: number
          method: string | null
          product: string
          status: Database["public"]["Enums"]["order_status"] | null
          updated_at: string | null
          url: string | null
          user_id: string | null
        }
        Insert: {
          amount: number
          bill_id?: string | null
          collection_id?: string | null
          created_at?: string | null
          id?: number
          method?: string | null
          product: string
          status?: Database["public"]["Enums"]["order_status"] | null
          updated_at?: string | null
          url?: string | null
          user_id?: string | null
        }
        Update: {
          amount?: number
          bill_id?: string | null
          collection_id?: string | null
          created_at?: string | null
          id?: number
          method?: string | null
          product?: string
          status?: Database["public"]["Enums"]["order_status"] | null
          updated_at?: string | null
          url?: string | null
          user_id?: string | null
        }
        Relationships: []
      }
      profiles: {
        Row: {
          created_at: string | null
          email: string
          full_name: string
          id: string
          is_active: boolean | null
          last_login: string | null
          updated_at: string | null
        }
        Insert: {
          created_at?: string | null
          email: string
          full_name: string
          id: string
          is_active?: boolean | null
          last_login?: string | null
          updated_at?: string | null
        }
        Update: {
          created_at?: string | null
          email?: string
          full_name?: string
          id?: string
          is_active?: boolean | null
          last_login?: string | null
          updated_at?: string | null
        }
        Relationships: []
      }
      user_roles: {
        Row: {
          created_at: string | null
          id: string
          role: Database["public"]["Enums"]["app_role"]
          user_id: string
        }
        Insert: {
          created_at?: string | null
          id?: string
          role: Database["public"]["Enums"]["app_role"]
          user_id: string
        }
        Update: {
          created_at?: string | null
          id?: string
          role?: Database["public"]["Enums"]["app_role"]
          user_id?: string
        }
        Relationships: []
      }
      wasapbot: {
        Row: {
          conv_last: string | null
          created_at: string | null
          current_node_id: string | null
          execution_id: string | null
          execution_status:
            | Database["public"]["Enums"]["execution_status"]
            | null
          flow_id: string | null
          flow_reference: string | null
          id_prospect: number
          instance: string | null
          last_node_id: string | null
          nama: string | null
          niche: string | null
          prospect_num: string | null
          stage: string | null
          status: string | null
          user_id: string | null
          user_input: string | null
          waiting_for_reply: boolean | null
        }
        Insert: {
          conv_last?: string | null
          created_at?: string | null
          current_node_id?: string | null
          execution_id?: string | null
          execution_status?:
            | Database["public"]["Enums"]["execution_status"]
            | null
          flow_id?: string | null
          flow_reference?: string | null
          id_prospect?: number
          instance?: string | null
          last_node_id?: string | null
          nama?: string | null
          niche?: string | null
          prospect_num?: string | null
          stage?: string | null
          status?: string | null
          user_id?: string | null
          user_input?: string | null
          waiting_for_reply?: boolean | null
        }
        Update: {
          conv_last?: string | null
          created_at?: string | null
          current_node_id?: string | null
          execution_id?: string | null
          execution_status?:
            | Database["public"]["Enums"]["execution_status"]
            | null
          flow_id?: string | null
          flow_reference?: string | null
          id_prospect?: number
          instance?: string | null
          last_node_id?: string | null
          nama?: string | null
          niche?: string | null
          prospect_num?: string | null
          stage?: string | null
          status?: string | null
          user_id?: string | null
          user_input?: string | null
          waiting_for_reply?: boolean | null
        }
        Relationships: []
      }
    }
    Views: {
      [_ in never]: never
    }
    Functions: {
      has_role: {
        Args: {
          _role: Database["public"]["Enums"]["app_role"]
          _user_id: string
        }
        Returns: boolean
      }
    }
    Enums: {
      ai_model:
        | "openai/gpt-5-chat"
        | "openai/gpt-5-mini"
        | "openai/chatgpt-4o-latest"
        | "openai/gpt-4.1"
        | "google/gemini-2.5-pro"
        | "google/gemini-pro-1.5"
      app_role: "admin" | "user"
      execution_status: "active" | "completed" | "failed"
      message_sender: "user" | "bot" | "staff"
      message_type: "text" | "image" | "document" | "audio" | "video"
      order_status: "Pending" | "Processing" | "Success" | "Failed"
      provider_type: "whacenter" | "wablas" | "waha"
    }
    CompositeTypes: {
      [_ in never]: never
    }
  }
}

type DatabaseWithoutInternals = Omit<Database, "__InternalSupabase">

type DefaultSchema = DatabaseWithoutInternals[Extract<keyof Database, "public">]

export type Tables<
  DefaultSchemaTableNameOrOptions extends
    | keyof (DefaultSchema["Tables"] & DefaultSchema["Views"])
    | { schema: keyof DatabaseWithoutInternals },
  TableName extends DefaultSchemaTableNameOrOptions extends {
    schema: keyof DatabaseWithoutInternals
  }
    ? keyof (DatabaseWithoutInternals[DefaultSchemaTableNameOrOptions["schema"]]["Tables"] &
        DatabaseWithoutInternals[DefaultSchemaTableNameOrOptions["schema"]]["Views"])
    : never = never,
> = DefaultSchemaTableNameOrOptions extends {
  schema: keyof DatabaseWithoutInternals
}
  ? (DatabaseWithoutInternals[DefaultSchemaTableNameOrOptions["schema"]]["Tables"] &
      DatabaseWithoutInternals[DefaultSchemaTableNameOrOptions["schema"]]["Views"])[TableName] extends {
      Row: infer R
    }
    ? R
    : never
  : DefaultSchemaTableNameOrOptions extends keyof (DefaultSchema["Tables"] &
        DefaultSchema["Views"])
    ? (DefaultSchema["Tables"] &
        DefaultSchema["Views"])[DefaultSchemaTableNameOrOptions] extends {
        Row: infer R
      }
      ? R
      : never
    : never

export type TablesInsert<
  DefaultSchemaTableNameOrOptions extends
    | keyof DefaultSchema["Tables"]
    | { schema: keyof DatabaseWithoutInternals },
  TableName extends DefaultSchemaTableNameOrOptions extends {
    schema: keyof DatabaseWithoutInternals
  }
    ? keyof DatabaseWithoutInternals[DefaultSchemaTableNameOrOptions["schema"]]["Tables"]
    : never = never,
> = DefaultSchemaTableNameOrOptions extends {
  schema: keyof DatabaseWithoutInternals
}
  ? DatabaseWithoutInternals[DefaultSchemaTableNameOrOptions["schema"]]["Tables"][TableName] extends {
      Insert: infer I
    }
    ? I
    : never
  : DefaultSchemaTableNameOrOptions extends keyof DefaultSchema["Tables"]
    ? DefaultSchema["Tables"][DefaultSchemaTableNameOrOptions] extends {
        Insert: infer I
      }
      ? I
      : never
    : never

export type TablesUpdate<
  DefaultSchemaTableNameOrOptions extends
    | keyof DefaultSchema["Tables"]
    | { schema: keyof DatabaseWithoutInternals },
  TableName extends DefaultSchemaTableNameOrOptions extends {
    schema: keyof DatabaseWithoutInternals
  }
    ? keyof DatabaseWithoutInternals[DefaultSchemaTableNameOrOptions["schema"]]["Tables"]
    : never = never,
> = DefaultSchemaTableNameOrOptions extends {
  schema: keyof DatabaseWithoutInternals
}
  ? DatabaseWithoutInternals[DefaultSchemaTableNameOrOptions["schema"]]["Tables"][TableName] extends {
      Update: infer U
    }
    ? U
    : never
  : DefaultSchemaTableNameOrOptions extends keyof DefaultSchema["Tables"]
    ? DefaultSchema["Tables"][DefaultSchemaTableNameOrOptions] extends {
        Update: infer U
      }
      ? U
      : never
    : never

export type Enums<
  DefaultSchemaEnumNameOrOptions extends
    | keyof DefaultSchema["Enums"]
    | { schema: keyof DatabaseWithoutInternals },
  EnumName extends DefaultSchemaEnumNameOrOptions extends {
    schema: keyof DatabaseWithoutInternals
  }
    ? keyof DatabaseWithoutInternals[DefaultSchemaEnumNameOrOptions["schema"]]["Enums"]
    : never = never,
> = DefaultSchemaEnumNameOrOptions extends {
  schema: keyof DatabaseWithoutInternals
}
  ? DatabaseWithoutInternals[DefaultSchemaEnumNameOrOptions["schema"]]["Enums"][EnumName]
  : DefaultSchemaEnumNameOrOptions extends keyof DefaultSchema["Enums"]
    ? DefaultSchema["Enums"][DefaultSchemaEnumNameOrOptions]
    : never

export type CompositeTypes<
  PublicCompositeTypeNameOrOptions extends
    | keyof DefaultSchema["CompositeTypes"]
    | { schema: keyof DatabaseWithoutInternals },
  CompositeTypeName extends PublicCompositeTypeNameOrOptions extends {
    schema: keyof DatabaseWithoutInternals
  }
    ? keyof DatabaseWithoutInternals[PublicCompositeTypeNameOrOptions["schema"]]["CompositeTypes"]
    : never = never,
> = PublicCompositeTypeNameOrOptions extends {
  schema: keyof DatabaseWithoutInternals
}
  ? DatabaseWithoutInternals[PublicCompositeTypeNameOrOptions["schema"]]["CompositeTypes"][CompositeTypeName]
  : PublicCompositeTypeNameOrOptions extends keyof DefaultSchema["CompositeTypes"]
    ? DefaultSchema["CompositeTypes"][PublicCompositeTypeNameOrOptions]
    : never

export const Constants = {
  public: {
    Enums: {
      ai_model: [
        "openai/gpt-5-chat",
        "openai/gpt-5-mini",
        "openai/chatgpt-4o-latest",
        "openai/gpt-4.1",
        "google/gemini-2.5-pro",
        "google/gemini-pro-1.5",
      ],
      app_role: ["admin", "user"],
      execution_status: ["active", "completed", "failed"],
      message_sender: ["user", "bot", "staff"],
      message_type: ["text", "image", "document", "audio", "video"],
      order_status: ["Pending", "Processing", "Success", "Failed"],
      provider_type: ["whacenter", "wablas", "waha"],
    },
  },
} as const
