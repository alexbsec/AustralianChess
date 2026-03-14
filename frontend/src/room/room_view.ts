import type { 
    PieceColor, 
    GameState,
} from "../engine/types";
import type { UIState } from "../ui/types";
import { BOARD_SIZE, PIECE_ASSETS } from "../engine/constants";
import { positionsEqual, pieceColorToSide, pieceKindToName, displayToBoardPosition } from "../engine/coords";

export class RoomView {
    private container: HTMLDivElement;
    public boardElement!: HTMLDivElement;
    public turnBanner!: HTMLElement;
    public turnBadge!: HTMLElement;
    public turnTitle!: HTMLElement;
    public turnSubtitle!: HTMLElement;
    public statusText!: HTMLElement;
    public playerText!: HTMLElement;
    public activityText!: HTMLElement;

    constructor(container: HTMLDivElement) {
        this.container = container;
        this.initStaticLayout();
    }

    /**
     * Builds the initial skeleton of the page.
     */
    private initStaticLayout(): void {
        this.container.innerHTML = `
            <main class="room-page">
                <section class="room-header">
                    <div class="room-header-left">
                        <div class="hero-badge">Room</div>
                        <h1 class="room-title">Australian Chess</h1>
                        <p class="room-subtitle" id="room-id-label"></p>
                    </div>
                    <div class="room-header-actions">
                        <button class="btn btn-secondary" id="back-btn">Back</button>
                    </div>
                </section>
                <section class="room-layout">
                    <div class="room-board-panel" id="board-panel">
                        <section class="room-turn-banner room-turn-banner--waiting" id="turn-banner">
                            <span class="room-turn-badge" id="turn-badge">Waiting</span>
                            <h2 class="room-turn-title" id="turn-title">Connecting...</h2>
                            <p class="room-turn-subtitle" id="turn-subtitle">Initializing socket...</p>
                        </section>
                        <div class="room-board" id="game-board"></div>
                    </div>
                    <aside class="room-side-panel">
                        ${this.renderInfoCard("Match", "connecting-status")}
                        ${this.renderInfoCard("You", "player-info")}
                        ${this.renderInfoCard("Latest", "activity-log")}
                    </aside>
                </section>
            </main>
        `;

        // Cache elements
        this.boardElement = this.container.querySelector("#game-board")!;
        this.turnBanner = this.container.querySelector("#turn-banner")!;
        this.turnBadge = this.container.querySelector("#turn-badge")!;
        this.turnTitle = this.container.querySelector("#turn-title")!;
        this.turnSubtitle = this.container.querySelector("#turn-subtitle")!;
        this.statusText = this.container.querySelector("#connecting-status")!;
        this.playerText = this.container.querySelector("#player-info")!;
        this.activityText = this.container.querySelector("#activity-log")!;
    }

    private renderInfoCard(title: string, contentId: string): string {
        return `
            <article class="room-info-card">
                <h2 class="room-info-title">${title}</h2>
                <p class="room-info-text" id="${contentId}">...</p>
            </article>
        `;
    }

    /**
     * Completely re-renders the board squares based on UI and Game state.
     */
    public renderBoard(gameState: GameState, playerColor: PieceColor, uiState: UIState): void {
        this.boardElement.innerHTML = "";
        
        for (let row = 0; row < BOARD_SIZE; row++) {
            for (let col = 0; col < BOARD_SIZE; col++) {
                const square = this.createSquare(row, col, gameState, playerColor, uiState);
                this.boardElement.appendChild(square);
            }
        }
    }

    private createSquare(
        displayRow: number, 
        displayCol: number, 
        gameState: GameState, 
        playerColor: PieceColor, 
        uiState: UIState
    ): HTMLDivElement {
        const boardPos = displayToBoardPosition({ row: displayRow, col: displayCol }, playerColor);
        const squareData = gameState.board.data[boardPos.row][boardPos.col];
        
        const square = document.createElement("div");
        square.className = `room-board-square ${(displayRow + displayCol) % 2 === 0 ? "light" : "dark"}`;
        square.dataset.row = String(displayRow);
        square.dataset.col = String(displayCol);

        if (positionsEqual(uiState.selected, boardPos)) {
            square.classList.add("selected");
        }

        const isPossibleMove = uiState.possibleMoves.some(m => positionsEqual(m, boardPos));
        if (isPossibleMove) {
            square.classList.add("possible-move");
            const dot = document.createElement("div");
            dot.className = "move-dot";
            square.appendChild(dot);
        }

        if (squareData?.piece) {
            const name = pieceKindToName(squareData.piece.kind);
            const side = pieceColorToSide(squareData.piece.color);
            
            const img = document.createElement("img");
            img.className = "room-piece";
            img.src = PIECE_ASSETS[side][name];
            img.draggable = false;
            square.appendChild(img);
        }

        return square;
    }

    /**
     * Updates the banner text without needing a full board re-render.
     */
    public updateStatus(title: string, subtitle: string, badge: string, className: string): void {
        this.turnTitle.textContent = title;
        this.turnSubtitle.textContent = subtitle;
        this.turnBadge.textContent = badge;
        this.turnBanner.className = `room-turn-banner ${className}`;
    }

    public updateRoomId(roomId: string): void {
        const label = this.container.querySelector("#room-id-label");
        if (label) label.textContent = `Room ID: ${roomId}`;
    }
}
