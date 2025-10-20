import React, { useState } from 'react';
import { Handle, Position, NodeProps } from '@xyflow/react';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Bot, Edit3 } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface AIPromptNodeData {
  label: string;
  prompt: string;
  model: string;
  temperature?: number;
}

const AI_MODELS = [
  { value: 'openai/gpt-5-chat', label: 'GPT-5 Chat' },
  { value: 'openai/gpt-5-mini', label: 'GPT-5 Mini' },
  { value: 'openai/chatgpt-4o-latest', label: 'GPT-4o Latest' },
  { value: 'openai/gpt-4.1', label: 'GPT-4.1' },
  { value: 'google/gemini-2.5-pro', label: 'Gemini 2.5 Pro' },
  { value: 'google/gemini-pro-1.5', label: 'Gemini Pro 1.5' },
];

export default function AIPromptNode({ data, id }: NodeProps<AIPromptNodeData>) {
  const [isEditing, setIsEditing] = useState(false);
  const [prompt, setPrompt] = useState(data.prompt || 'Please respond helpfully to the user');
  const [label, setLabel] = useState(data.label || 'AI Prompt');
  const [model, setModel] = useState(data.model || 'openai/gpt-4.1');

  const handleSave = () => {
    data.prompt = prompt;
    data.label = label;
    data.model = model;
    setIsEditing(false);
  };

  return (
    <Card className="min-w-[280px] shadow-lg border-2 border-emerald-200 bg-emerald-50">
      <CardHeader className="pb-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <div className="p-1 bg-emerald-100 rounded">
              <Bot className="h-4 w-4 text-emerald-600" />
            </div>
            {isEditing ? (
              <Input
                value={label}
                onChange={(e) => setLabel(e.target.value)}
                className="h-6 text-sm font-medium"
                placeholder="Node label"
              />
            ) : (
              <span className="text-sm font-medium text-emerald-800">{label}</span>
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
          <div className="space-y-3">
            <div>
              <label className="text-xs font-medium text-gray-700">AI Model</label>
              <Select value={model} onValueChange={setModel}>
                <SelectTrigger className="text-sm">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {AI_MODELS.map((modelOption) => (
                    <SelectItem key={modelOption.value} value={modelOption.value}>
                      {modelOption.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div>
              <label className="text-xs font-medium text-gray-700">Prompt</label>
              <Textarea
                value={prompt}
                onChange={(e) => setPrompt(e.target.value)}
                placeholder="Enter AI prompt..."
                className="min-h-[80px] text-sm"
              />
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
          <div className="space-y-2">
            <div className="text-xs text-emerald-600 bg-white p-2 rounded border">
              <div className="font-medium">Model: {AI_MODELS.find(m => m.value === model)?.label}</div>
              <div className="mt-1 text-gray-700">{prompt}</div>
            </div>
          </div>
        )}
      </CardContent>

      {/* Input Handle */}
      <Handle
        type="target"
        position={Position.Top}
        className="w-3 h-3 bg-emerald-500"
      />
      
      {/* Output Handle */}
      <Handle
        type="source"
        position={Position.Bottom}
        className="w-3 h-3 bg-emerald-500"
      />
    </Card>
  );
}