import {
    type PieceColor,
    type GameState,
    type PieceKind,
    colorToPieceSide,
} from "../engine/types";
import type { UIState } from "../ui/types";
import { BISHOP_PIECE, BOARD_SIZE, KANGAROO_PIECE, KNIGHT_PIECE, OLIGARCH_PIECE, PIECE_ASSETS, QUEEN_PIECE, ROOK_PIECE } from "../engine/constants";
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
                        <button class="btn btn-secondary" id="copy-link-btn" style="width: fit-content; padding: 6px 12px;">Copy room link to share with friends!</button>
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

    public showInactivityWarning(seconds: number): void {
        let banner = document.getElementById("inactivity-warning");
        if (!banner) {
            banner = document.createElement("div");
            banner.id = "inactivity-warning";
            Object.assign(banner.style, {
                position: "fixed",
                bottom: "24px",
                left: "50%",
                transform: "translateX(-50%)",
                background: "#b45309",
                color: "white",
                padding: "12px 24px",
                borderRadius: "8px",
                fontWeight: "600",
                zIndex: "9999",
                boxShadow: "0 4px 12px rgba(0,0,0,0.3)",
                whiteSpace: "nowrap",
            });
            document.body.appendChild(banner);
        }
        banner.textContent = `A move happen! Room closes in ${seconds}s`;
    }

    public hideInactivityWarning(): void {
        document.getElementById("inactivity-warning")?.remove();
    }

    public showGameOverModal(reason: string, onNewGame: () => void, newGameLabel = "Create New Room"): void {
        if (document.getElementById("game-over-modal")) return;

        const modalOverlay = document.createElement("div");
        modalOverlay.id = "game-over-modal";
        modalOverlay.className = "modal-overlay";

        modalOverlay.style.zIndex = "99999";
        modalOverlay.removeAttribute('hidden');
        modalOverlay.style.display = 'flex';

        modalOverlay.innerHTML = `
            <div class="modal-content" style="position: relative; z-index: 100000;">
                <h2 style="margin-top: 0; color: var(--primary);">Game Over</h2>
                <p style="margin-bottom: 24px; color: var(--text);">${reason}</p>
                <div class="modal-actions">
                    <button id="modal-new-room-btn" class="btn btn-primary">${newGameLabel}</button>
                    <button id="modal-close-btn" class="btn btn-secondary">Close</button>
                </div>
            </div>
        `;

        document.body.appendChild(modalOverlay);

        document.getElementById("modal-new-room-btn")?.addEventListener("click", () => {
            modalOverlay.remove();
            onNewGame();
        });

        document.getElementById("modal-close-btn")?.addEventListener("click", () => {
            modalOverlay.remove();
        });
    }

    public showPromotionPicker(color: PieceColor, onPick: (piece: PieceKind) => void): void {
        const existing = document.getElementById("promotion-overlay");
        if (existing) existing.remove();

        const side = colorToPieceSide(color);
        const options: { kind: PieceKind; name: string }[] = [
            { kind: QUEEN_PIECE, name: "queen" },
            { kind: ROOK_PIECE, name: "rook" },
            { kind: BISHOP_PIECE, name: "bishop" },
            { kind: KNIGHT_PIECE, name: "knight" },
            { kind: OLIGARCH_PIECE, name: "oligarch" },
            { kind: KANGAROO_PIECE, name: "kangaroo" },
        ];

        const overlay = document.createElement("div");
        overlay.id = "promotion-overlay";
        overlay.className = "promotion-overlay";

        const picker = document.createElement("div");
        picker.className = "promotion-picker";

        const label = document.createElement("p");
        label.className = "promotion-label";
        label.textContent = "Promote pawn to:";
        picker.appendChild(label);

        const grid = document.createElement("div");
        grid.className = "promotion-grid";

        options.forEach(({ kind, name }) => {
            const btn = document.createElement("button");
            btn.className = "promotion-option";
            btn.title = name

            const img = document.createElement("img");
            img.src = PIECE_ASSETS[side][name]
            img.alt = name;
            img.draggable = false;

            btn.appendChild(img);
            btn.addEventListener("click", () => {
                overlay.remove();
                onPick(kind);
            });
            grid.appendChild(btn);
        });

        picker.appendChild(grid);
        overlay.appendChild(picker);
        document.body.appendChild(overlay);
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
            const isBeingDragged = uiState.draggingPos && positionsEqual(uiState.draggingPos, boardPos);

            if (!isBeingDragged) {
                const name = pieceKindToName(squareData.piece.kind);
                const side = pieceColorToSide(squareData.piece.color);

                const img = document.createElement("img");
                img.className = "room-piece";
                img.src = PIECE_ASSETS[side][name];
                img.draggable = false;
                square.appendChild(img);
            }
        }

        if (uiState.lastMove) {
            const { from, to } = uiState.lastMove;
            if (positionsEqual(boardPos, from) || positionsEqual(boardPos, to)) {
                square.classList.add("last-move");
            }
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

        const copyBtn = this.container.querySelector<HTMLButtonElement>("#copy-link-btn");
        if (copyBtn) {
            const copyIconSvg = `<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="flex-shrink:0"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>`;
            const checkIconSvg = `<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="flex-shrink:0"><polyline points="20 6 9 16 4 11"/></svg>`;

            copyBtn.style.display = "flex";
            copyBtn.style.alignItems = "center";
            copyBtn.style.gap = "6px";
            copyBtn.innerHTML = `${copyIconSvg} Copy room link`;

            copyBtn.addEventListener("click", () => {
                const url = `${window.location.origin}/room?id=${roomId}`;
                navigator.clipboard.writeText(url).then(() => {
                    copyBtn.innerHTML = `${checkIconSvg} Room link copied to clipboard`;
                    setTimeout(() => {
                        copyBtn.innerHTML = `${copyIconSvg} Copy room link`;
                    }, 2000);
                }).catch(() => {
                    prompt("Copy this link:", url);
                });
            });
        }
    }
}
