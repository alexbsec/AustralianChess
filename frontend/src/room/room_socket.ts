import type {
    GameState,
    PieceColor,
    PieceKind,
    Position,
} from "../engine/types";
import type { WsMessage } from "../ws/types";
import { buildRoomWsUrl, isGameStateMessage, isMoveResultMessage, isPromotionResultMessage, isRoomStatusMessage } from "../ws/utils";

export type SocketHandlers = {
    onGameState: (state: GameState, playerColor?: PieceColor) => void;
    onRoomStatus: (success: boolean, gameStarted: boolean) => void;
    onMoveResult: (moved: boolean, state: GameState) => void;
    onPromoteResult: (promoted: boolean, state: GameState) => void;
    onOpen: () => void;
    onError: (event: Event) => void;
    onClose: () => void;
};

export class RoomSocket {
    private roomId: string;
    private handlers: SocketHandlers;
    private socket: WebSocket | null = null;

    constructor(roomId: string, handlers: SocketHandlers) {
        this.roomId = roomId;
        this.handlers = handlers;
    }

    /**
     * Initializes the connection and attaches event listeners.
     */
    public connect(): void {
        const url = buildRoomWsUrl(this.roomId);
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
            } else if (isPromotionResultMessage(message)) {
                this.handlers.onPromoteResult(message.promoted, message.game_state);
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

    public sendPromote(requesteeColor: PieceColor, promoteTo: PieceKind, pawnPosition: Position, destination: Position): void {
        if (!this.socket || this.socket.readyState !== WebSocket.OPEN) return;

        const payload = {
            type: "pawn_promote",
            requestee_color: requesteeColor,
            promote_to: promoteTo,
            pawn_position: { row: pawnPosition.row, col: pawnPosition.col },
            destination_position: { row: destination.row, col: destination.col },
        };
        
        this.socket.send(JSON.stringify(payload));
    }

    public disconnect(): void {
        this.socket?.close();
    }
}
