import React, { useState } from 'react';
import { Handle, Position, NodeProps } from '@xyflow/react';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Video, Edit3 } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface VideoNodeData {
  label: string;
  videoUrl: string;
}

export default function VideoNode({ data, id }: NodeProps<VideoNodeData>) {
  const [isEditing, setIsEditing] = useState(false);
  const [videoUrl, setVideoUrl] = useState(data.videoUrl || '');
  const [label, setLabel] = useState(data.label || 'Video');

  const handleSave = () => {
    data.videoUrl = videoUrl;
    data.label = label;
    setIsEditing(false);
  };

  return (
    <Card className="min-w-[250px] shadow-lg border-2 border-red-200 bg-red-50">
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <div className="p-1 bg-red-100 rounded">
              <Video className="h-4 w-4 text-red-600" />
            </div>
            {isEditing ? (
              <Input value={label} onChange={(e) => setLabel(e.target.value)} className="h-6 text-sm font-medium" placeholder="Node label" />
            ) : (
              <span className="text-sm font-medium text-red-800">{label}</span>
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
            <Input value={videoUrl} onChange={(e) => setVideoUrl(e.target.value)} placeholder="Video URL..." className="text-sm" />
            <div className="flex space-x-2">
              <Button size="sm" onClick={handleSave}>Save</Button>
              <Button size="sm" variant="outline" onClick={() => setIsEditing(false)}>Cancel</Button>
            </div>
          </div>
        ) : (
          <div className="text-sm text-gray-700 bg-white p-2 rounded border">{videoUrl || 'No video URL set'}</div>
        )}
      </CardContent>
      <Handle type="target" position={Position.Top} className="w-3 h-3 bg-red-500" />
      <Handle type="source" position={Position.Bottom} className="w-3 h-3 bg-red-500" />
    </Card>
  );
}