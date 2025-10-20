import React, { useState } from 'react';
import { Handle, Position, NodeProps } from '@xyflow/react';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { MessageSquare, Edit3 } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface MessageNodeData {
  label: string;
  message: string;
}

export default function MessageNode({ data, id }: NodeProps<MessageNodeData>) {
  const [isEditing, setIsEditing] = useState(false);
  const [message, setMessage] = useState(data.message || 'Hello! How can I help you?');
  const [label, setLabel] = useState(data.label || 'Message');

  const handleSave = () => {
    // Update node data
    data.message = message;
    data.label = label;
    setIsEditing(false);
  };

  return (
    <Card className="min-w-[250px] shadow-lg border-2 border-blue-200 bg-blue-50">
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <div className="p-1 bg-blue-100 rounded">
              <MessageSquare className="h-4 w-4 text-blue-600" />
            </div>
            {isEditing ? (
              <Input
                value={label}
                onChange={(e) => setLabel(e.target.value)}
                className="h-6 text-sm font-medium"
                placeholder="Node label"
              />
            ) : (
              <span className="text-sm font-medium text-blue-800">{label}</span>
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
            <Textarea
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              placeholder="Enter your message..."
              className="min-h-[60px] text-sm"
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
            {message}
          </div>
        )}
      </CardContent>

      {/* Input Handle */}
      <Handle
        type="target"
        position={Position.Top}
        className="w-3 h-3 bg-blue-500"
      />
      
      {/* Output Handle */}
      <Handle
        type="source"
        position={Position.Bottom}
        className="w-3 h-3 bg-blue-500"
      />
    </Card>
  );
}