import React, { useState } from 'react';
import { Handle, Position, NodeProps } from '@xyflow/react';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Volume2, Edit3 } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface AudioNodeData {
  label: string;
  audioUrl: string;
}

export default function AudioNode({ data, id }: NodeProps<AudioNodeData>) {
  const [isEditing, setIsEditing] = useState(false);
  const [audioUrl, setAudioUrl] = useState(data.audioUrl || '');
  const [label, setLabel] = useState(data.label || 'Audio');

  const handleSave = () => {
    data.audioUrl = audioUrl;
    data.label = label;
    setIsEditing(false);
  };

  return (
    <Card className="min-w-[250px] shadow-lg border-2 border-orange-200 bg-orange-50">
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <div className="p-1 bg-orange-100 rounded">
              <Volume2 className="h-4 w-4 text-orange-600" />
            </div>
            {isEditing ? (
              <Input
                value={label}
                onChange={(e) => setLabel(e.target.value)}
                className="h-6 text-sm font-medium"
                placeholder="Node label"
              />
            ) : (
              <span className="text-sm font-medium text-orange-800">{label}</span>
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
            <Input
              value={audioUrl}
              onChange={(e) => setAudioUrl(e.target.value)}
              placeholder="Audio URL..."
              className="text-sm"
            />
            <div className="flex space-x-2">
              <Button size="sm" onClick={handleSave}>Save</Button>
              <Button size="sm" variant="outline" onClick={() => setIsEditing(false)}>Cancel</Button>
            </div>
          </div>
        ) : (
          <div className="text-sm text-gray-700 bg-white p-2 rounded border">
            {audioUrl || 'No audio URL set'}
          </div>
        )}
      </CardContent>
      <Handle type="target" position={Position.Top} className="w-3 h-3 bg-orange-500" />
      <Handle type="source" position={Position.Bottom} className="w-3 h-3 bg-orange-500" />
    </Card>
  );
}