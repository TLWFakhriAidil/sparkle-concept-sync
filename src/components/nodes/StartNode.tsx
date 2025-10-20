import React from 'react';
import { Handle, Position, NodeProps } from '@xyflow/react';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Play } from 'lucide-react';

interface StartNodeData {
  label: string;
}

export default function StartNode({ data }: NodeProps<StartNodeData>) {
  return (
    <Card className="min-w-[200px] shadow-lg border-2 border-green-200 bg-green-50">
      <CardHeader className="pb-2">
        <div className="flex items-center space-x-2">
          <div className="p-1 bg-green-100 rounded">
            <Play className="h-4 w-4 text-green-600" />
          </div>
          <span className="text-sm font-medium text-green-800">Start</span>
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        <div className="text-xs text-green-600 bg-white p-2 rounded border">
          Flow starts here
        </div>
      </CardContent>

      {/* Only output handle for start node */}
      <Handle
        type="source"
        position={Position.Bottom}
        className="w-3 h-3 bg-green-500"
      />
    </Card>
  );
}