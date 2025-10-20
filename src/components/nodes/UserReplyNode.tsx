import React, { useState } from 'react';
import { Handle, Position, NodeProps } from '@xyflow/react';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { User, Edit3 } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface UserReplyNodeData {
  label: string;
  timeout: number;
}

export default function UserReplyNode({ data, id }: NodeProps<UserReplyNodeData>) {
  const [isEditing, setIsEditing] = useState(false);
  const [timeout, setTimeout] = useState(data.timeout || 30000);
  const [label, setLabel] = useState(data.label || 'User Reply');

  const handleSave = () => {
    data.timeout = timeout;
    data.label = label;
    setIsEditing(false);
  };

  return (
    <Card className="min-w-[200px] shadow-lg border-2 border-cyan-200 bg-cyan-50">
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <div className="p-1 bg-cyan-100 rounded">
              <User className="h-4 w-4 text-cyan-600" />
            </div>
            {isEditing ? (
              <Input value={label} onChange={(e) => setLabel(e.target.value)} className="h-6 text-sm font-medium" placeholder="Node label" />
            ) : (
              <span className="text-sm font-medium text-cyan-800">{label}</span>
            )}
          </div>
          <Button size="sm" variant="ghost" onClick={() => setIsEditing(!isEditing)}>
            <Edit3 className="h-3 w-3" />
          </Button>
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        {isEditing ? (
          <div className="space-y-2">
            <div>
              <label className="text-xs font-medium text-gray-700">Timeout (ms)</label>
              <Input type="number" value={timeout} onChange={(e) => setTimeout(Number(e.target.value))} className="text-sm" />
            </div>
            <div className="flex space-x-2">
              <Button size="sm" onClick={handleSave}>Save</Button>
              <Button size="sm" variant="outline" onClick={() => setIsEditing(false)}>Cancel</Button>
            </div>
          </div>
        ) : (
          <div className="text-sm text-gray-700 bg-white p-2 rounded border text-center">
            Wait for user input ({timeout/1000}s)
          </div>
        )}
      </CardContent>
      <Handle type="target" position={Position.Top} className="w-3 h-3 bg-cyan-500" />
      <Handle type="source" position={Position.Bottom} className="w-3 h-3 bg-cyan-500" />
    </Card>
  );
}