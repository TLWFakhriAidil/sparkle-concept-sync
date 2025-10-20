import React, { useState } from 'react';
import { Handle, Position, NodeProps } from '@xyflow/react';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Settings, Edit3 } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface StageNodeData {
  label: string;
  stage: string;
}

export default function StageNode({ data, id }: NodeProps<StageNodeData>) {
  const [isEditing, setIsEditing] = useState(false);
  const [stage, setStage] = useState(data.stage || 'Welcome');
  const [label, setLabel] = useState(data.label || 'Stage');

  const handleSave = () => {
    data.stage = stage;
    data.label = label;
    setIsEditing(false);
  };

  return (
    <Card className="min-w-[200px] shadow-lg border-2 border-gray-200 bg-gray-50">
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <div className="p-1 bg-gray-100 rounded">
              <Settings className="h-4 w-4 text-gray-600" />
            </div>
            {isEditing ? (
              <Input value={label} onChange={(e) => setLabel(e.target.value)} className="h-6 text-sm font-medium" placeholder="Node label" />
            ) : (
              <span className="text-sm font-medium text-gray-800">{label}</span>
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
            <Input value={stage} onChange={(e) => setStage(e.target.value)} placeholder="Stage name..." className="text-sm" />
            <div className="flex space-x-2">
              <Button size="sm" onClick={handleSave}>Save</Button>
              <Button size="sm" variant="outline" onClick={() => setIsEditing(false)}>Cancel</Button>
            </div>
          </div>
        ) : (
          <div className="text-sm text-gray-700 bg-white p-2 rounded border text-center">
            Set stage: {stage}
          </div>
        )}
      </CardContent>
      <Handle type="target" position={Position.Top} className="w-3 h-3 bg-gray-500" />
      <Handle type="source" position={Position.Bottom} className="w-3 h-3 bg-gray-500" />
    </Card>
  );
}