import { useEffect, useState } from 'react';
import { TopBar } from '@/components/TopBar';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { supabase } from '@/integrations/supabase/client';
import { useAuth } from '@/contexts/AuthContext';
import { BarChart3, MessageSquare, Users, TrendingUp } from 'lucide-react';

const Analytics = () => {
  const { user } = useAuth();
  const [stats, setStats] = useState({
    totalConversations: 0,
    activeConversations: 0,
    completedConversations: 0,
    failedConversations: 0,
    totalMessages: 0,
  });

  useEffect(() => {
    const fetchAnalytics = async () => {
      if (!user) return;

      const [conversationsRes, messagesRes] = await Promise.all([
        supabase.from('ai_whatsapp').select('execution_status').eq('user_id', user.id),
        supabase.from('conversation_log').select('id', { count: 'exact', head: true }).eq('user_id', user.id),
      ]);

      if (conversationsRes.data) {
        const active = conversationsRes.data.filter(c => c.execution_status === 'active').length;
        const completed = conversationsRes.data.filter(c => c.execution_status === 'completed').length;
        const failed = conversationsRes.data.filter(c => c.execution_status === 'failed').length;

        setStats({
          totalConversations: conversationsRes.data.length,
          activeConversations: active,
          completedConversations: completed,
          failedConversations: failed,
          totalMessages: messagesRes.count || 0,
        });
      }
    };

    fetchAnalytics();
  }, [user]);

  const analyticsCards = [
    {
      title: 'Total Conversations',
      value: stats.totalConversations,
      description: 'All time conversations',
      icon: MessageSquare,
      color: 'text-blue-500',
    },
    {
      title: 'Active Conversations',
      value: stats.activeConversations,
      description: 'Currently in progress',
      icon: TrendingUp,
      color: 'text-green-500',
    },
    {
      title: 'Completed',
      value: stats.completedConversations,
      description: 'Successfully finished',
      icon: Users,
      color: 'text-purple-500',
    },
    {
      title: 'Total Messages',
      value: stats.totalMessages,
      description: 'Messages exchanged',
      icon: BarChart3,
      color: 'text-orange-500',
    },
  ];

  return (
    <div className="min-h-screen bg-background">
      <TopBar />
      <main className="container mx-auto px-4 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-2">Analytics</h1>
          <p className="text-muted-foreground">Monitor your chatbot performance</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          {analyticsCards.map((card, index) => (
            <Card key={index}>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">{card.title}</CardTitle>
                <card.icon className={`h-4 w-4 ${card.color}`} />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{card.value}</div>
                <p className="text-xs text-muted-foreground">{card.description}</p>
              </CardContent>
            </Card>
          ))}
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Conversation Breakdown</CardTitle>
            <CardDescription>Status distribution of your conversations</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <div className="w-3 h-3 rounded-full bg-green-500"></div>
                  <span>Active</span>
                </div>
                <span className="font-medium">{stats.activeConversations}</span>
              </div>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <div className="w-3 h-3 rounded-full bg-blue-500"></div>
                  <span>Completed</span>
                </div>
                <span className="font-medium">{stats.completedConversations}</span>
              </div>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <div className="w-3 h-3 rounded-full bg-red-500"></div>
                  <span>Failed</span>
                </div>
                <span className="font-medium">{stats.failedConversations}</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </main>
    </div>
  );
};

export default Analytics;
