import React, { useState } from 'react';
import { Handle, Position, NodeProps } from '@xyflow/react';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Clock, Edit3 } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface DelayNodeData {
  label: string;
  delay: number;
}

export default function DelayNode({ data, id }: NodeProps<DelayNodeData>) {
  const [isEditing, setIsEditing] = useState(false);
  const [delay, setDelay] = useState(data.delay || 1000);
  const [label, setLabel] = useState(data.label || 'Delay');

  const handleSave = () => {
    data.delay = delay;
    data.label = label;
    setIsEditing(false);
  };

  const formatDelay = (ms: number) => {
    if (ms < 1000) return `${ms}ms`;
    if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`;
    return `${(ms / 60000).toFixed(1)}m`;
  };

  return (
    <Card className="min-w-[200px] shadow-lg border-2 border-yellow-200 bg-yellow-50">
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <div className="p-1 bg-yellow-100 rounded">
              <Clock className="h-4 w-4 text-yellow-600" />
            </div>
            {isEditing ? (
              <Input
                value={label}
                onChange={(e) => setLabel(e.target.value)}
                className="h-6 text-sm font-medium"
                placeholder="Node label"
              />
            ) : (
              <span className="text-sm font-medium text-yellow-800">{label}</span>
            )}
          </div>
          <Button
            size="sm"
            variant="ghost"
            onClick={() => setIsEditing(!isEditing)}
          >
            <Edit3 className="h-3 w-3" />
          </Button>
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        {isEditing ? (
          <div className="space-y-2">
            <div>
              <label className="text-xs font-medium text-gray-700">Delay (milliseconds)</label>
              <Input
                type="number"
                value={delay}
                onChange={(e) => setDelay(Number(e.target.value))}
                placeholder="1000"
                className="text-sm"
                min="0"
                step="100"
              />
            </div>
            <div className="text-xs text-gray-500">
              Preview: {formatDelay(delay)}
            </div>
            <div className="flex space-x-2">
              <Button size="sm" onClick={handleSave}>
                Save
              </Button>
              <Button 
                size="sm" 
                variant="outline"
                onClick={() => setIsEditing(false)}
              >
                Cancel
              </Button>
            </div>
          </div>
        ) : (
          <div className="text-sm text-gray-700 bg-white p-2 rounded border text-center">
            Wait {formatDelay(delay)}
          </div>
        )}
      </CardContent>

      <Handle
        type="target"
        position={Position.Top}
        className="w-3 h-3 bg-yellow-500"
      />
      
      <Handle
        type="source"
        position={Position.Bottom}
        className="w-3 h-3 bg-yellow-500"
      />
    </Card>
  );
}