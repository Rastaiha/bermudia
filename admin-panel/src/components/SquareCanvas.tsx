// SquareCanvas.tsx
import React, { useRef } from "react";
import type { Island } from '../models/Territory';

const DEFAULT_COLORS = ["#3b82f6", "#22c55e", "#ef4444", "#f59e0b", "#8b5cf6"];

const SquareCanvas : React.FC<{
    allIslands: Island[];
    onUpdateIsland: (id: string, updates: Partial<Island>) => void;
    onSelectIsland: (id: string) => void;
    selectedIsland: string | null
}> = ({allIslands, onUpdateIsland, onSelectIsland, selectedIsland}) => {

    const width = 900;
    const height = 600;
    const size = 60;

    const mapX = (x: number) => (x*Math.max(0, width - size));
    const mapY = (y: number) => (y*Math.max(0, height - size));

    const containerRef = useRef<HTMLDivElement | null>(null);
    // holds { id, offsetX, offsetY } where offsets are pointer-within-square measured
    const dragRef = useRef<{ id: string; offsetX: number; offsetY: number } | null>(null);

    const handleMouseDown = (e: React.MouseEvent, id: string) => {
        e.preventDefault();
        const container = containerRef.current;
        if (!container) return;
        const rect = container.getBoundingClientRect();
        const sq = allIslands.find((s) => s.id === id);
        if (!sq) return;

        // pointer position inside the square (measured relative to container)
        const offsetX = e.clientX - rect.left - mapX(sq.x);
        const offsetY = e.clientY - rect.top - mapY(sq.y);
        dragRef.current = { id, offsetX, offsetY };

        // move handler (attached to window so dragging stays smooth)
        const onMove = (ev: MouseEvent) => {
            if (!dragRef.current || !containerRef.current) return;
            const r = containerRef.current.getBoundingClientRect();
            const { id: dragId, offsetX: ox, offsetY: oy } = dragRef.current;
            let nx = ev.clientX - r.left - ox;
            let ny = ev.clientY - r.top - oy;
            nx = Math.max(0, Math.min(nx, r.width - size));
            ny = Math.max(0, Math.min(ny, r.height - size));

            onUpdateIsland(dragId, {x: nx/Math.max(0, width - size), y: ny/Math.max(0, height - size)})
        };

        const onUp = () => {
            window.removeEventListener("mousemove", onMove);
            window.removeEventListener("mouseup", onUp);
            dragRef.current = null;
        };

        window.addEventListener("mousemove", onMove);
        window.addEventListener("mouseup", onUp);

        onSelectIsland(id);
    };

    return (
        <div
            className="m-4"
            ref={containerRef}
            style={{
                width: typeof width === "number" ? `${width}px` : width,
                height: typeof height === "number" ? `${height}px` : height,
                background: '#ffffff',
                border: "1px solid rgba(0,0,0,0.12)",
                borderRadius: 8,
                position: "relative",
                overflow: "hidden",
                userSelect: "none",
            }}
        >
            {allIslands.map((island, i) => (
                <div
                    key={island.id}
                    onMouseDown={(e) => handleMouseDown(e, island.id)}
                    style={{
                        position: "absolute",
                        left: mapX(island.x),
                        top: mapY(island.y),
                        width: size,
                        height: size,
                        background: DEFAULT_COLORS[i % DEFAULT_COLORS.length],
                        borderRadius: 6,
                        boxShadow: "0 2px 6px rgba(0,0,0,0.15)",
                        cursor: "grab",
                    }}
                    className={`text-center content-center ${island.id === selectedIsland ? 'border-4 border-gray-300' : ''}`}
                >
                    {island.id}
                </div>
            ))}
        </div>
    );
}

export default SquareCanvas ;
