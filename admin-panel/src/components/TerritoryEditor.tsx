import React, { useState } from 'react';
import { Plus, Trash2, Save } from 'lucide-react';
import SquareCanvas from './SquareCanvas';
import type { Island, Territory, Edge } from '../models/Territory';

// Territory Info Component
const TerritoryInfo: React.FC<{
  territory: Territory;
  onUpdate: (updates: Partial<Territory>) => void;
  onExport: () => void;
}> = ({ territory, onUpdate, onExport }) => (
  <div className="p-4 border-b">
    <div className="flex items-center justify-between mb-4">
      <h1 className="text-xl font-bold text-gray-800">Territory Editor</h1>
      <button
        onClick={onExport}
        className="px-3 py-2 bg-green-500 text-white text-sm rounded hover:bg-green-600 flex items-center gap-2"
      >
        <Save className="w-4 h-4" />
        Export JSON
      </button>
    </div>
    <div className="space-y-3">
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Territory Name</label>
        <input
          type="text"
          placeholder="Territory Name"
          value={territory.name}
          onChange={(e) => onUpdate({ name: e.target.value })}
          className="w-full px-3 py-2 border rounded-md text-sm"
        />
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Background Asset</label>
        <input
          type="text"
          placeholder="background1.jpg"
          value={territory.backgroundAsset}
          onChange={(e) => onUpdate({ backgroundAsset: e.target.value })}
          className="w-full px-3 py-2 border rounded-md text-sm"
        />
      </div>
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Start Island</label>
        <select
          value={territory.startIsland}
          onChange={(e) => onUpdate({ startIsland: e.target.value })}
          className="w-full px-3 py-2 border rounded-md text-sm"
        >
          <option value="">Select start island</option>
          {territory.islands.map(island => (
            <option key={island.id} value={island.id}>
              {island.id}
            </option>
          ))}
        </select>
      </div>
    </div>
  </div>
);

// Islands List Component
const IslandsList: React.FC<{
  islands: Island[];
  selectedIsland: string | null;
  onSelectIsland: (id: string) => void;
  onCreateIsland: (id: string, name: string) => void;
  onDeleteIsland: (id: string) => void;
}> = ({ islands, selectedIsland, onSelectIsland, onCreateIsland, onDeleteIsland }) => {
  const [newIslandName, setNewIslandName] = useState('');
  const [newIslandId, setNewIslandId] = useState('');

  const handleCreate = () => {
    if (!newIslandName.trim()) return;
    onCreateIsland(newIslandId, newIslandName);
    setNewIslandName('');
    setNewIslandId('');
  };

  const handleNewIslandId = (id: string) => {
    const pattern : RegExp = /^\w{1,}$/
    if (!pattern.test(id)) {
      return
    }
    setNewIslandId(id);
  }

  return (
    <div className="p-4 border-b">
      <h2 className="font-semibold text-gray-700 mb-3">Islands</h2>
      
      <div className="flex flex-col gap-2 mb-3">
        <input
          type="text"
          placeholder="New island id"
          value={newIslandId}
          onChange={(e) => handleNewIslandId(e.target.value)}
          className="flex-1 px-3 py-2 border rounded-md text-sm"
          onKeyDown={(e) => e.key === 'Enter' && handleCreate()}
          required
        />
        <input
          type="text"
          placeholder="New island name"
          value={newIslandName}
          onChange={(e) => setNewIslandName(e.target.value)}
          className="flex-1 px-3 py-2 border rounded-md text-sm"
          onKeyDown={(e) => e.key === 'Enter' && handleCreate()}
          required
        />
        <button
          onClick={handleCreate}
          className="flex flex-1 content-center px-3 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600"
        >
          <Plus className="flex-1 w-4 h-4" />
        </button>
      </div>

      <div className="space-y-2 max-h-64 overflow-y-auto">
        {islands.map((island) => (
          <div
            key={island.id}
            className={`p-3 border rounded-md cursor-pointer transition-colors ${
              selectedIsland === island.id ? 'bg-blue-100 border-blue-500' : 'hover:bg-gray-50'
            }`}
            onClick={() => onSelectIsland(island.id)}
          >
            <div className="flex items-center justify-between">
              <div className="flex-1">
                <div className="font-medium text-sm">{island.id}</div>
                <div className="text-xs text-gray-500">{island.name}</div>
              </div>
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  onDeleteIsland(island.id);
                }}
                className="p-1 text-red-500 hover:bg-red-100 rounded"
              >
                <Trash2 className="w-4 h-4" />
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

// Special Islands List Component
const SpecialIslandsList: React.FC<{
  title: string;
  islands: Island[];
  specialIslands: { id: string }[];
  onAdd: (islandId: string) => void;
  onRemove: (islandId: string) => void;
  bgColor: string;
}> = ({ title, islands, specialIslands, onAdd, onRemove, bgColor }) => (
  <div className="p-4 border-b">
    <h3 className="font-semibold text-gray-700 mb-3">{title}</h3>
    
    <div className="space-y-2 mb-4 max-h-32 overflow-y-auto">
      {specialIslands.map((specialIsland) => {
        const island = islands.find(i => i.id === specialIsland.id);
        return island ? (
          <div key={specialIsland.id} className={`flex items-center justify-between p-2 ${bgColor} rounded border`}>
            <span className="text-sm">{island.id}</span>
            <button
              onClick={() => onRemove(specialIsland.id)}
              className="text-red-500 hover:bg-red-100 p-1 rounded"
            >
              <Trash2 className="w-3 h-3" />
            </button>
          </div>
        ) : null;
      })}
    </div>
    
    <select
      onChange={(e) => e.target.value && onAdd(e.target.value)}
      value=""
      className="w-full px-3 py-2 border rounded-md text-sm"
    >
      <option value="">Add island to {title.toLowerCase()}</option>
      {islands
        .filter(island => !specialIslands.some(special => special.id === island.id))
        .map(island => (
          <option key={island.id} value={island.id}>
            {island.id}
          </option>
        ))}
    </select>
  </div>
);

// Island Properties Component
const IslandProperties: React.FC<{
  island: Island;
  allIslands: Island[];
  edges: Edge[];
  onUpdateIsland: (updates: Partial<Island>) => void;
  onToggleConnection: (fromId: string, toId: string) => void;
}> = ({ island, allIslands, edges, onUpdateIsland, onToggleConnection }) => (
  <div className="p-4 space-y-4">
    <h3 className="font-semibold text-gray-700">Island Properties</h3>
    
    <div className="space-y-3">
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Name</label>
        <input
          type="text"
          value={island.name}
          onChange={(e) => onUpdateIsland({ name: e.target.value })}
          className="w-full px-3 py-2 border rounded-md text-sm"
        />
      </div>
      
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Icon Asset</label>
        <input
          type="text"
          value={island.iconAsset}
          onChange={(e) => onUpdateIsland({ iconAsset: e.target.value })}
          className="w-full px-3 py-2 border rounded-md text-sm"
          placeholder="island1.png"
        />
      </div>

      <div className="border rounded-md p-2 bg-gray-200 border-gray-500">
        <label className="block text-sm font-medium text-sm mb-1">Position</label>
        <div className="font-medium text-sm text-gray-600">X: {island.x.toFixed(4)}</div>
        <div className="font-medium text-sm text-gray-600">Y: {island.y.toFixed(4)}</div>
        <div className="font-medium text-sm text-gray-600">Width: {island.width}</div>
        <div className="font-medium text-sm text-gray-600">Height: {island.height}</div>
      </div>
      
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">Connections To Other Islands</label>
        <div className="space-y-2 max-h-40 overflow-y-auto border rounded-md p-2">
          {allIslands
            .filter(targetIsland => targetIsland.id !== island.id)
            .map(targetIsland => {
              const isConnected = edges.some(edge => 
                edge.from === island.id && edge.to === targetIsland.id
                || edge.to === island.id && edge.from === targetIsland.id
              );
              
              return (
                <label key={targetIsland.id} className="flex items-center gap-2 p-1">
                  <input
                    type="checkbox"
                    checked={isConnected}
                    onChange={() => onToggleConnection(island.id, targetIsland.id)}
                    className="rounded"
                  />
                  <span className="text-sm">{targetIsland.id}</span>
                </label>
              );
            })}
        </div>
      </div>
    </div>
  </div>
);

// Main Territory Editor Component
const TerritoryEditor: React.FC = () => {
  const [territory, setTerritory] = useState<Territory>({
    id: 'territory1',
    name: 'جزایر اسرارآمیز',
    backgroundAsset: 'background1.jpg',
    startIsland: '',
    islands: [],
    edges: [],
    refuelIslands: [],
    terminalIslands: [],
    islandPrerequisites: {}
  });

  const [selectedIsland, setSelectedIsland] = useState<string | null>(null);

  const updateTerritory = (updates: Partial<Territory>) => {
    setTerritory(prev => ({ ...prev, ...updates }));
  };

  const createIsland = (id: string, name: string) => {
    const width = 0.12;
    const height = 0.18;
    const newIsland: Island = {
      id,
      name,
      x: Math.random() * (1-width),
      y: Math.random() * (1-height),
      width,
      height,
      iconAsset: 'island1.png'
    };

    setTerritory(prev => ({
      ...prev,
      islands: [...prev.islands, newIsland]
    }));
    setSelectedIsland(newIsland.id);
  };

  const deleteIsland = (islandId: string) => {
    setTerritory(prev => ({
      ...prev,
      islands: prev.islands.filter(island => island.id !== islandId),
      edges: prev.edges.filter(edge => edge.from !== islandId && edge.to !== islandId),
      refuelIslands: prev.refuelIslands.filter(ref => ref.id !== islandId),
      terminalIslands: prev.terminalIslands.filter(term => term.id !== islandId),
      startIsland: prev.startIsland === islandId ? '' : prev.startIsland
    }));
    
    if (selectedIsland === islandId) {
      setSelectedIsland(null);
    }
  };

  const updateIsland = (islandId: string, updates: Partial<Island>) => {
    setTerritory(prev => ({
      ...prev,
      islands: prev.islands.map(island => 
        island.id === islandId ? { ...island, ...updates } : island
      )
    }));
  };

  const addToSpecialList = (islandId: string, listType: 'refuel' | 'terminal') => {
    setTerritory(prev => {
      const targetList = listType === 'refuel' ? 'refuelIslands' : 'terminalIslands';
      const alreadyExists = prev[targetList].some(item => item.id === islandId);
      
      if (alreadyExists) return prev;
      
      return {
        ...prev,
        [targetList]: [...prev[targetList], { id: islandId }]
      };
    });
  };

  const removeFromSpecialList = (islandId: string, listType: 'refuel' | 'terminal') => {
    setTerritory(prev => {
      const targetList = listType === 'refuel' ? 'refuelIslands' : 'terminalIslands';
      return {
        ...prev,
        [targetList]: prev[targetList].filter(item => item.id !== islandId)
      };
    });
  };

  const toggleConnection = (fromId: string, toId: string) => {
    setTerritory(prev => {
      const existingEdge = prev.edges.find(edge => (edge.from === fromId && edge.to === toId) || (edge.from === toId && edge.to === fromId));
      
      if (existingEdge) {
        return {
          ...prev,
          edges: prev.edges.filter(edge => !(edge.from === fromId && edge.to === toId || edge.from === toId && edge.to === fromId))
        };
      } else {
        return {
          ...prev,
          edges: [...prev.edges, { from: fromId, to: toId }]
        };
      }
    });
  };

  const exportTerritory = () => {
    console.log('Territory JSON:', JSON.stringify(territory, null, 2));
    alert('Territory JSON has been logged to console');
  };

  const selectedIslandData = selectedIsland ? territory.islands.find(i => i.id === selectedIsland) : null;

  return (
    <div className="flex h-screen bg-gray-100">
      {/* Left Panel */}
      <div className="w-90 bg-white shadow-lg overflow-y-auto">
        <TerritoryInfo 
          territory={territory} 
          onUpdate={updateTerritory} 
          onExport={exportTerritory}
        />
        
        <IslandsList
          islands={territory.islands}
          selectedIsland={selectedIsland}
          onSelectIsland={setSelectedIsland}
          onCreateIsland={createIsland}
          onDeleteIsland={deleteIsland}
        />

        <SpecialIslandsList
          title="Refuel Islands"
          islands={territory.islands}
          specialIslands={territory.refuelIslands}
          onAdd={(id) => addToSpecialList(id, 'refuel')}
          onRemove={(id) => removeFromSpecialList(id, 'refuel')}
          bgColor="bg-yellow-50"
        />

        <SpecialIslandsList
          title="Terminal Islands"
          islands={territory.islands}
          specialIslands={territory.terminalIslands}
          onAdd={(id) => addToSpecialList(id, 'terminal')}
          onRemove={(id) => removeFromSpecialList(id, 'terminal')}
          bgColor="bg-purple-50"
        />
      </div>

      <SquareCanvas
        allIslands={territory.islands}
        onUpdateIsland={updateIsland}
        onSelectIsland={setSelectedIsland}
        selectedIsland={selectedIsland}
      />

      {/* Right Panel */}
      <div className="flex-1 bg-white shadow-lg">
        {selectedIslandData ? (
          <IslandProperties
            island={selectedIslandData}
            allIslands={territory.islands}
            edges={territory.edges}
            onUpdateIsland={(updates) => updateIsland(selectedIslandData.id, updates)}
            onToggleConnection={toggleConnection}
          />
        ) : (
          <div className="flex items-center justify-center h-full text-gray-500">
            <div className="text-center">
              <div className="text-lg font-medium mb-2">No Island Selected</div>
              <div className="text-sm">Select an island from the left panel to edit its properties</div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default TerritoryEditor;