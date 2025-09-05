interface Island {
  id: string;
  name: string;
  x: number;
  y: number;
  width: number;
  height: number;
  iconAsset: string;
}

interface Edge {
  from: string;
  to: string;
}

interface Territory {
  id: string;
  name: string;
  backgroundAsset: string;
  startIsland: string;
  islands: Island[];
  edges: Edge[];
  refuelIslands: { id: string }[];
  terminalIslands: { id: string }[];
  islandPrerequisites: { [key: string]: string[] };
}

export type {Island, Edge, Territory};