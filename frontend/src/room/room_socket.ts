import type {
    GameState,
    PieceColor,
    Position,
} from "../engine/types";
import type { WsMessage } from "../ws/types";
import { buildRoomWsUrl, isGameStateMessage, isMoveResultMessage, isRoomStatusMessage } from "../ws/utils";

export type SocketHandlers = {
    onGameState: (state: GameState, playerColor?: PieceColor) => void;
    onRoomStatus: (success: boolean, gameStarted: boolean) => void;
    onMoveResult: (moved: boolean, state: GameState) => void;
    onOpen: () => void;
    onError: (event: Event) => void;
    onClose: () => void;
};

export class RoomSocket {
    private roomId: string;
    private playerId: string;
    private handlers: SocketHandlers;
    private socket: WebSocket | null = null;

    constructor(roomId: string, playerId: string, handlers: SocketHandlers) {
        this.roomId = roomId;
        this.playerId = playerId;
        this.handlers = handlers;
    }

    /**
     * Initializes the connection and attaches event listeners.
     */
    public connect(): void {
        const url = buildRoomWsUrl(this.roomId, this.playerId);
        this.socket = new WebSocket(url);

        this.socket.addEventListener("open", () => this.handlers.onOpen());
        this.socket.addEventListener("message", (event) => this.handleMessage(event));
        this.socket.addEventListener("error", (event) => this.handlers.onError(event));
        this.socket.addEventListener("close", () => this.handlers.onClose());
    }

    /**
     * Parses incoming WebSocket messages and dispatches them to the appropriate handler.
     */
    private handleMessage(event: MessageEvent<string>): void {
        try {
            const message: WsMessage = JSON.parse(event.data);

            if (isGameStateMessage(message)) {
                this.handlers.onGameState(message.game_state, message.player_color);
            } else if (isRoomStatusMessage(message)) {
                this.handlers.onRoomStatus(message.success, message.game_started);
            } else if (isMoveResultMessage(message)) {
                this.handlers.onMoveResult(message.moved, message.game_state);
            }
        } catch (error) {
            console.error("Failed to parse websocket payload:", error);
        }
    }

    /**
     * Sends a move request to the server.
     */
    public sendMove(requesteeColor: PieceColor, from: Position, to: Position): void {
        if (!this.socket || this.socket.readyState !== WebSocket.OPEN) return;

        const payload = {
            type: "make_move",
            requestee_color: requesteeColor,
            from_pos: { row: from.row, col: from.col },
            to_pos: { row: to.row, col: to.col },
        };

        this.socket.send(JSON.stringify(payload));
    }

    public disconnect(): void {
        this.socket?.close();
    }
}
