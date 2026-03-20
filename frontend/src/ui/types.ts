import type { Position } from "../engine/types";

export type UIState = {
    selected: Position | null;
    possibleMoves: Position[];
    isMouseDown: boolean;
    draggingPos: Position | null;
    lastMove: { from: Position; to: Position} | null;
};
