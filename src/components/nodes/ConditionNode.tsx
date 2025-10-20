import React, { useState } from 'react';
import { Handle, Position, NodeProps } from '@xyflow/react';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { GitBranch, Edit3 } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface ConditionNodeData {
  label: string;
  condition: string;
}

export default function ConditionNode({ data, id }: NodeProps<ConditionNodeData>) {
  const [isEditing, setIsEditing] = useState(false);
  const [condition, setCondition] = useState(data.condition || 'user_input contains "yes"');
  const [label, setLabel] = useState(data.label || 'Condition');

  const handleSave = () => {
    data.condition = condition;
    data.label = label;
    setIsEditing(false);
  };

  return (
    <Card className="min-w-[250px] shadow-lg border-2 border-indigo-200 bg-indigo-50">
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <div className="p-1 bg-indigo-100 rounded">
              <GitBranch className="h-4 w-4 text-indigo-600" />
            </div>
            {isEditing ? (
              <Input
                value={label}
                onChange={(e) => setLabel(e.target.value)}
                className="h-6 text-sm font-medium"
                placeholder="Node label"
              />
            ) : (
              <span className="text-sm font-medium text-indigo-800">{label}</span>
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
              value={condition}
              onChange={(e) => setCondition(e.target.value)}
              placeholder='e.g., user_input contains "yes"'
              className="min-h-[60px] text-sm font-mono"
            />
            <div className="text-xs text-gray-500">
              Available variables: user_input, stage, previous_messages
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
          <div className="text-sm text-gray-700 bg-white p-2 rounded border font-mono">
            {condition}
          </div>
        )}
      </CardContent>

      {/* Input Handle */}
      <Handle
        type="target"
        position={Position.Top}
        className="w-3 h-3 bg-indigo-500"
      />
      
      {/* Output Handles - True/False */}
      <Handle
        type="source"
        position={Position.Right}
        id="true"
        className="w-3 h-3 bg-green-500"
        style={{ top: '60%' }}
      />
      <Handle
        type="source"
        position={Position.Left}
        id="false"
        className="w-3 h-3 bg-red-500"
        style={{ top: '60%' }}
      />
      
      {/* Labels for outputs */}
      <div className="absolute -right-8 top-1/2 transform -translate-y-1/2 text-xs text-green-600 font-medium">
        Yes
      </div>
      <div className="absolute -left-8 top-1/2 transform -translate-y-1/2 text-xs text-red-600 font-medium">
        No
      </div>
    </Card>
  );
}