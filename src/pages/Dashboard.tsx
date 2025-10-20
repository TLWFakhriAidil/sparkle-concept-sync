import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { TopBar } from '@/components/TopBar';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { supabase } from '@/integrations/supabase/client';
import { useAuth } from '@/contexts/AuthContext';
import { MessageSquare, Workflow, Settings, BarChart3, Plus } from 'lucide-react';

const Dashboard = () => {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [stats, setStats] = useState({
    flows: 0,
    devices: 0,
    conversations: 0,
    activeChats: 0,
  });

  useEffect(() => {
    const fetchStats = async () => {
      if (!user) return;

      const [flowsRes, devicesRes, conversationsRes] = await Promise.all([
        supabase.from('chatbot_flows').select('id', { count: 'exact', head: true }).eq('user_id', user.id),
        supabase.from('device_settings').select('id', { count: 'exact', head: true }).eq('user_id', user.id),
        supabase.from('ai_whatsapp').select('id', { count: 'exact', head: true }).eq('user_id', user.id),
      ]);

      setStats({
        flows: flowsRes.count || 0,
        devices: devicesRes.count || 0,
        conversations: conversationsRes.count || 0,
        activeChats: 0,
      });
    };

    fetchStats();
  }, [user]);

  const quickActions = [
    {
      title: 'Create Flow',
      description: 'Build a new chatbot flow',
      icon: Workflow,
      action: () => navigate('/flow-builder'),
    },
    {
      title: 'Add Device',
      description: 'Connect a WhatsApp device',
      icon: Settings,
      action: () => navigate('/devices'),
    },
    {
      title: 'View Analytics',
      description: 'Check your performance',
      icon: BarChart3,
      action: () => navigate('/analytics'),
    },
  ];

  return (
    <div className="min-h-screen bg-background">
      <TopBar />
      <main className="container mx-auto px-4 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-2">Dashboard</h1>
          <p className="text-muted-foreground">Manage your WhatsApp AI chatbots</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Flows</CardTitle>
              <Workflow className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.flows}</div>
              <p className="text-xs text-muted-foreground">Active chatbot flows</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Devices</CardTitle>
              <Settings className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.devices}</div>
              <p className="text-xs text-muted-foreground">Connected devices</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Conversations</CardTitle>
              <MessageSquare className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.conversations}</div>
              <p className="text-xs text-muted-foreground">Total conversations</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Active Chats</CardTitle>
              <BarChart3 className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.activeChats}</div>
              <p className="text-xs text-muted-foreground">Currently active</p>
            </CardContent>
          </Card>
        </div>

        <div className="mb-8">
          <h2 className="text-2xl font-bold mb-4">Quick Actions</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {quickActions.map((action, index) => (
              <Card key={index} className="hover:shadow-lg transition-shadow cursor-pointer" onClick={action.action}>
                <CardHeader>
                  <div className="p-3 bg-primary/10 rounded-full w-fit mb-2">
                    <action.icon className="h-6 w-6 text-primary" />
                  </div>
                  <CardTitle>{action.title}</CardTitle>
                  <CardDescription>{action.description}</CardDescription>
                </CardHeader>
                <CardContent>
                  <Button variant="ghost" className="w-full">
                    <Plus className="mr-2 h-4 w-4" />
                    Get Started
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </main>
    </div>
  );
};

export default Dashboard;
