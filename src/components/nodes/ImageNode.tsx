import React, { useState } from 'react';
import { Handle, Position, NodeProps } from '@xyflow/react';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Image, Edit3 } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface ImageNodeData {
  label: string;
  imageUrl: string;
  caption?: string;
}

export default function ImageNode({ data, id }: NodeProps<ImageNodeData>) {
  const [isEditing, setIsEditing] = useState(false);
  const [imageUrl, setImageUrl] = useState(data.imageUrl || '');
  const [label, setLabel] = useState(data.label || 'Image');
  const [caption, setCaption] = useState(data.caption || '');

  const handleSave = () => {
    data.imageUrl = imageUrl;
    data.label = label;
    data.caption = caption;
    setIsEditing(false);
  };

  return (
    <Card className="min-w-[250px] shadow-lg border-2 border-purple-200 bg-purple-50">
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <div className="p-1 bg-purple-100 rounded">
              <Image className="h-4 w-4 text-purple-600" />
            </div>
            {isEditing ? (
              <Input
                value={label}
                onChange={(e) => setLabel(e.target.value)}
                className="h-6 text-sm font-medium"
                placeholder="Node label"
              />
            ) : (
              <span className="text-sm font-medium text-purple-800">{label}</span>
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
            <Input
              value={imageUrl}
              onChange={(e) => setImageUrl(e.target.value)}
              placeholder="Image URL..."
              className="text-sm"
            />
            <Input
              value={caption}
              onChange={(e) => setCaption(e.target.value)}
              placeholder="Caption (optional)"
              className="text-sm"
            />
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
          <div className="text-sm text-gray-700 bg-white p-2 rounded border">
            {imageUrl ? (
              <div>
                <div className="truncate">{imageUrl}</div>
                {caption && <div className="text-xs text-gray-500 mt-1">{caption}</div>}
              </div>
            ) : (
              <div className="text-gray-400">No image URL set</div>
            )}
          </div>
        )}
      </CardContent>

      <Handle type="target" position={Position.Top} className="w-3 h-3 bg-purple-500" />
      <Handle type="source" position={Position.Bottom} className="w-3 h-3 bg-purple-500" />
    </Card>
  );
}