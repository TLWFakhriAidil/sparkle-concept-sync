import { useEffect, useState } from 'react';
import { TopBar } from '@/components/TopBar';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { supabase } from '@/integrations/supabase/client';
import { useAuth } from '@/contexts/AuthContext';
import { useToast } from '@/hooks/use-toast';
import { Plus, Trash2, Edit } from 'lucide-react';
import { DeviceSettings as DeviceSettingsType } from '@/types/chatbot';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';

const DeviceSettings = () => {
  const { user } = useAuth();
  const { toast } = useToast();
  const [devices, setDevices] = useState<DeviceSettingsType[]>([]);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingDevice, setEditingDevice] = useState<DeviceSettingsType | null>(null);
  const [formData, setFormData] = useState<{
    device_id: string;
    phone_number: string;
    provider: 'whacenter' | 'wablas' | 'waha';
    api_key: string;
    api_key_option: 'openai/gpt-5-chat' | 'openai/gpt-5-mini' | 'openai/chatgpt-4o-latest' | 'openai/gpt-4.1' | 'google/gemini-2.5-pro' | 'google/gemini-pro-1.5';
    webhook_id: string;
    instance: string;
  }>({
    device_id: '',
    phone_number: '',
    provider: 'wablas',
    api_key: '',
    api_key_option: 'openai/gpt-4.1',
    webhook_id: '',
    instance: '',
  });

  useEffect(() => {
    fetchDevices();
  }, [user]);

  const fetchDevices = async () => {
    if (!user) return;

    const { data, error } = await supabase
      .from('device_settings')
      .select('*')
      .eq('user_id', user.id)
      .order('created_at', { ascending: false });

    if (error) {
      toast({
        title: "Error",
        description: "Failed to fetch devices",
        variant: "destructive",
      });
    } else {
      setDevices(data || []);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user) return;

    const deviceData = {
      ...formData,
      id: editingDevice?.id || `device_${Date.now()}`,
      id_device: formData.device_id,
      user_id: user.id,
    };

    const { error } = editingDevice
      ? await supabase.from('device_settings').update(deviceData).eq('id', editingDevice.id)
      : await supabase.from('device_settings').insert([deviceData]);

    if (error) {
      toast({
        title: "Error",
        description: error.message,
        variant: "destructive",
      });
    } else {
      toast({
        title: "Success",
        description: `Device ${editingDevice ? 'updated' : 'created'} successfully`,
      });
      setIsDialogOpen(false);
      resetForm();
      fetchDevices();
    }
  };

  const handleDelete = async (id: string) => {
    const { error } = await supabase.from('device_settings').delete().eq('id', id);

    if (error) {
      toast({
        title: "Error",
        description: "Failed to delete device",
        variant: "destructive",
      });
    } else {
      toast({
        title: "Success",
        description: "Device deleted successfully",
      });
      fetchDevices();
    }
  };

  const resetForm = () => {
    setFormData({
      device_id: '',
      phone_number: '',
      provider: 'wablas',
      api_key: '',
      api_key_option: 'openai/gpt-4.1',
      webhook_id: '',
      instance: '',
    });
    setEditingDevice(null);
  };

  const openEditDialog = (device: DeviceSettingsType) => {
    setEditingDevice(device);
    setFormData({
      device_id: device.id_device || '',
      phone_number: device.phone_number || '',
      provider: device.provider,
      api_key: device.api_key || '',
      api_key_option: device.api_key_option as any,
      webhook_id: device.webhook_id || '',
      instance: device.instance || '',
    });
    setIsDialogOpen(true);
  };

  return (
    <div className="min-h-screen bg-background">
      <TopBar />
      <main className="container mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold mb-2">Device Settings</h1>
            <p className="text-muted-foreground">Manage your WhatsApp devices</p>
          </div>

          <Dialog open={isDialogOpen} onOpenChange={(open) => {
            setIsDialogOpen(open);
            if (!open) resetForm();
          }}>
            <DialogTrigger asChild>
              <Button>
                <Plus className="mr-2 h-4 w-4" />
                Add Device
              </Button>
            </DialogTrigger>
            <DialogContent className="max-w-2xl">
              <DialogHeader>
                <DialogTitle>{editingDevice ? 'Edit Device' : 'Add New Device'}</DialogTitle>
              </DialogHeader>
              <form onSubmit={handleSubmit} className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="device_id">Device ID</Label>
                    <Input
                      id="device_id"
                      value={formData.device_id}
                      onChange={(e) => setFormData({ ...formData, device_id: e.target.value })}
                      required
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="phone_number">Phone Number</Label>
                    <Input
                      id="phone_number"
                      value={formData.phone_number}
                      onChange={(e) => setFormData({ ...formData, phone_number: e.target.value })}
                    />
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="provider">Provider</Label>
                    <Select value={formData.provider} onValueChange={(value: any) => setFormData({ ...formData, provider: value })}>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="wablas">Wablas</SelectItem>
                        <SelectItem value="whacenter">Whacenter</SelectItem>
                        <SelectItem value="waha">WAHA</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="ai_model">AI Model</Label>
                    <Select value={formData.api_key_option} onValueChange={(value: any) => setFormData({ ...formData, api_key_option: value })}>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="openai/gpt-5-chat">GPT-5 Chat</SelectItem>
                        <SelectItem value="openai/gpt-5-mini">GPT-5 Mini</SelectItem>
                        <SelectItem value="openai/chatgpt-4o-latest">GPT-4o Latest</SelectItem>
                        <SelectItem value="openai/gpt-4.1">GPT-4.1</SelectItem>
                        <SelectItem value="google/gemini-2.5-pro">Gemini 2.5 Pro</SelectItem>
                        <SelectItem value="google/gemini-pro-1.5">Gemini Pro 1.5</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="api_key">API Key</Label>
                  <Input
                    id="api_key"
                    type="password"
                    value={formData.api_key}
                    onChange={(e) => setFormData({ ...formData, api_key: e.target.value })}
                    required
                  />
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="webhook_id">Webhook ID</Label>
                    <Input
                      id="webhook_id"
                      value={formData.webhook_id}
                      onChange={(e) => setFormData({ ...formData, webhook_id: e.target.value })}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="instance">Instance</Label>
                    <Input
                      id="instance"
                      value={formData.instance}
                      onChange={(e) => setFormData({ ...formData, instance: e.target.value })}
                    />
                  </div>
                </div>

                <div className="flex justify-end gap-2">
                  <Button type="button" variant="outline" onClick={() => setIsDialogOpen(false)}>
                    Cancel
                  </Button>
                  <Button type="submit">
                    {editingDevice ? 'Update' : 'Create'} Device
                  </Button>
                </div>
              </form>
            </DialogContent>
          </Dialog>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {devices.map((device) => (
            <Card key={device.id}>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  <span>{device.phone_number || device.id_device}</span>
                  <div className="flex gap-2">
                    <Button size="icon" variant="ghost" onClick={() => openEditDialog(device)}>
                      <Edit className="h-4 w-4" />
                    </Button>
                    <Button size="icon" variant="ghost" onClick={() => handleDelete(device.id)}>
                      <Trash2 className="h-4 w-4 text-destructive" />
                    </Button>
                  </div>
                </CardTitle>
                <CardDescription>{device.provider}</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-2 text-sm">
                  <div>
                    <span className="font-medium">Device ID:</span> {device.id_device}
                  </div>
                  <div>
                    <span className="font-medium">AI Model:</span> {device.api_key_option}
                  </div>
                  {device.webhook_id && (
                    <div>
                      <span className="font-medium">Webhook:</span> {device.webhook_id}
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          ))}

          {devices.length === 0 && (
            <Card className="col-span-full">
              <CardContent className="flex flex-col items-center justify-center py-12">
                <p className="text-muted-foreground mb-4">No devices configured yet</p>
                <Button onClick={() => setIsDialogOpen(true)}>
                  <Plus className="mr-2 h-4 w-4" />
                  Add Your First Device
                </Button>
              </CardContent>
            </Card>
          )}
        </div>
      </main>
    </div>
  );
};

export default DeviceSettings;
