import React, { useCallback, useRef, useState } from 'react';
import {
  ReactFlow,
  Background,
  Controls,
  MiniMap,
  addEdge,
  useNodesState,
  useEdgesState,
  Connection,
  Edge,
  Node,
  ReactFlowProvider,
  Panel,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { 
  MessageSquare, 
  Image, 
  Volume2, 
  Video, 
  Clock, 
  GitBranch, 
  Play, 
  User, 
  Bot,
  Settings
} from 'lucide-react';

import MessageNode from './nodes/MessageNode';
import ImageNode from './nodes/ImageNode';
import AudioNode from './nodes/AudioNode';
import VideoNode from './nodes/VideoNode';
import DelayNode from './nodes/DelayNode';
import ConditionNode from './nodes/ConditionNode';
import UserReplyNode from './nodes/UserReplyNode';
import AIPromptNode from './nodes/AIPromptNode';
import StageNode from './nodes/StageNode';
import StartNode from './nodes/StartNode';

const nodeTypes = {
  start: StartNode,
  message: MessageNode,
  image: ImageNode,
  audio: AudioNode,
  video: VideoNode,
  delay: DelayNode,
  condition: ConditionNode,
  user_reply: UserReplyNode,
  ai_prompt: AIPromptNode,
  stage: StageNode,
};

const nodeTemplates = [
  { type: 'start', label: 'Start', icon: Play, color: 'bg-green-100 text-green-800' },
  { type: 'message', label: 'Message', icon: MessageSquare, color: 'bg-blue-100 text-blue-800' },
  { type: 'image', label: 'Image', icon: Image, color: 'bg-purple-100 text-purple-800' },
  { type: 'audio', label: 'Audio', icon: Volume2, color: 'bg-orange-100 text-orange-800' },
  { type: 'video', label: 'Video', icon: Video, color: 'bg-red-100 text-red-800' },
  { type: 'delay', label: 'Delay', icon: Clock, color: 'bg-yellow-100 text-yellow-800' },
  { type: 'condition', label: 'Condition', icon: GitBranch, color: 'bg-indigo-100 text-indigo-800' },
  { type: 'user_reply', label: 'User Reply', icon: User, color: 'bg-cyan-100 text-cyan-800' },
  { type: 'ai_prompt', label: 'AI Prompt', icon: Bot, color: 'bg-emerald-100 text-emerald-800' },
  { type: 'stage', label: 'Stage', icon: Settings, color: 'bg-gray-100 text-gray-800' },
];

interface ChatbotBuilderProps {
  flowId?: string;
  initialNodes?: Node[];
  initialEdges?: Edge[];
  onSave?: (nodes: Node[], edges: Edge[]) => void;
  onTest?: (nodes: Node[], edges: Edge[]) => void;
}

export default function ChatbotBuilder({ 
  flowId, 
  initialNodes = [], 
  initialEdges = [], 
  onSave, 
  onTest 
}: ChatbotBuilderProps) {
  const reactFlowWrapper = useRef<HTMLDivElement>(null);
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
  const [reactFlowInstance, setReactFlowInstance] = useState<any>(null);
  const [selectedNodeType, setSelectedNodeType] = useState<string | null>(null);

  const onConnect = useCallback(
    (params: Edge | Connection) => setEdges((eds) => addEdge(params, eds)),
    [setEdges]
  );

  const onDragOver = useCallback((event: React.DragEvent) => {
    event.preventDefault();
    event.dataTransfer.dropEffect = 'move';
  }, []);

  const onDrop = useCallback(
    (event: React.DragEvent) => {
      event.preventDefault();

      const type = event.dataTransfer.getData('application/reactflow');

      if (typeof type === 'undefined' || !type) {
        return;
      }

      if (!reactFlowWrapper.current) return;

      const reactFlowBounds = reactFlowWrapper.current.getBoundingClientRect();
      const position = reactFlowInstance.screenToFlowPosition({
        x: event.clientX - reactFlowBounds.left,
        y: event.clientY - reactFlowBounds.top,
      });

      const newNode: Node = {
        id: `${type}_${Date.now()}`,
        type,
        position,
        data: {
          label: `${type} node`,
          // Default data based on node type
          ...(type === 'message' && { message: 'Hello! How can I help you?' }),
          ...(type === 'delay' && { delay: 1000 }),
          ...(type === 'condition' && { condition: 'user_input contains "yes"' }),
          ...(type === 'stage' && { stage: 'Welcome' }),
          ...(type === 'ai_prompt' && { prompt: 'Please respond helpfully to the user' }),
        },
      };

      setNodes((nds) => nds.concat(newNode));
    },
    [reactFlowInstance, setNodes]
  );

  const onDragStart = (event: React.DragEvent, nodeType: string) => {
    event.dataTransfer.setData('application/reactflow', nodeType);
    event.dataTransfer.effectAllowed = 'move';
    setSelectedNodeType(nodeType);
  };

  const handleSave = () => {
    if (onSave) {
      onSave(nodes, edges);
    }
  };

  const handleTest = () => {
    if (onTest) {
      onTest(nodes, edges);
    }
  };

  const clearFlow = () => {
    setNodes([]);
    setEdges([]);
  };

  const autoLayout = () => {
    // Simple auto-layout - arrange nodes in a grid
    const updatedNodes = nodes.map((node, index) => ({
      ...node,
      position: {
        x: (index % 3) * 300,
        y: Math.floor(index / 3) * 150,
      },
    }));
    setNodes(updatedNodes);
  };

  return (
    <div className="h-screen flex bg-gray-50">
      {/* Sidebar - Node Palette */}
      <div className="w-80 bg-white border-r border-gray-200 p-4 overflow-y-auto">
        <h3 className="text-lg font-semibold mb-4">Flow Nodes</h3>
        <div className="space-y-2">
          {nodeTemplates.map((template) => (
            <Card
              key={template.type}
              className={`cursor-grab hover:shadow-md transition-shadow ${
                selectedNodeType === template.type ? 'ring-2 ring-blue-500' : ''
              }`}
              draggable
              onDragStart={(event) => onDragStart(event, template.type)}
            >
              <CardContent className="p-3">
                <div className="flex items-center space-x-3">
                  <div className={`p-2 rounded-lg ${template.color}`}>
                    <template.icon className="h-5 w-5" />
                  </div>
                  <div>
                    <h4 className="font-medium">{template.label}</h4>
                    <p className="text-sm text-gray-500 capitalize">{template.type.replace('_', ' ')}</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Flow Actions */}
        <div className="mt-6 space-y-2">
          <Button onClick={handleSave} className="w-full" variant="default">
            Save Flow
          </Button>
          <Button onClick={handleTest} className="w-full" variant="outline">
            Test Flow
          </Button>
          <Button onClick={autoLayout} className="w-full" variant="outline">
            Auto Layout
          </Button>
          <Button onClick={clearFlow} className="w-full" variant="destructive">
            Clear All
          </Button>
        </div>

        {/* Flow Stats */}
        <div className="mt-6">
          <h4 className="font-medium mb-2">Flow Statistics</h4>
          <div className="space-y-1 text-sm">
            <div className="flex justify-between">
              <span>Nodes:</span>
              <Badge variant="secondary">{nodes.length}</Badge>
            </div>
            <div className="flex justify-between">
              <span>Connections:</span>
              <Badge variant="secondary">{edges.length}</Badge>
            </div>
          </div>
        </div>
      </div>

      {/* Main Flow Canvas */}
      <div className="flex-1 relative">
        <ReactFlowProvider>
          <div className="h-full" ref={reactFlowWrapper}>
            <ReactFlow
              nodes={nodes}
              edges={edges}
              onNodesChange={onNodesChange}
              onEdgesChange={onEdgesChange}
              onConnect={onConnect}
              onInit={setReactFlowInstance}
              onDrop={onDrop}
              onDragOver={onDragOver}
              nodeTypes={nodeTypes}
              fitView
              deleteKeyCode={['Backspace', 'Delete']}
              multiSelectionKeyCode={['Meta', 'Ctrl']}
              snapToGrid
              snapGrid={[15, 15]}
            >
              <Panel position="top-right" className="bg-white p-2 rounded-lg shadow-lg border">
                <div className="flex items-center space-x-2">
                  <Badge variant="outline">
                    {nodes.length} nodes
                  </Badge>
                  <Badge variant="outline">
                    {edges.length} edges
                  </Badge>
                  {flowId && (
                    <Badge variant="default">
                      Flow: {flowId}
                    </Badge>
                  )}
                </div>
              </Panel>

              <Background color="#f3f4f6" gap={20} />
              <Controls />
              <MiniMap 
                nodeColor={(node) => {
                  const template = nodeTemplates.find(t => t.type === node.type);
                  return template ? '#8b5cf6' : '#64748b';
                }}
                style={{
                  height: 120,
                  backgroundColor: '#f8fafc',
                }}
              />
            </ReactFlow>
          </div>
        </ReactFlowProvider>
      </div>
    </div>
  );
}

// Wrapper component to provide ReactFlow context
export function ChatbotBuilderWrapper(props: ChatbotBuilderProps) {
  return (
    <ReactFlowProvider>
      <ChatbotBuilder {...props} />
    </ReactFlowProvider>
  );
}