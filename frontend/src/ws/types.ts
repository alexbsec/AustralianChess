import type { GameState, PieceColor } from "../engine/types";

export type GameStateMessage = {
    type: "game_state";
    game_state: GameState;
    player_color?: PieceColor;
};

export type RoomStatusMessage = {
    room_id: string;
    player_id: string;
    success: boolean;
    game_started: boolean;
};

export type MoveResultMessage = {
    room_id: string;
    moved: boolean;
    game_state: GameState;
};

export type WsMessage = GameStateMessage | RoomStatusMessage | MoveResultMessage;
