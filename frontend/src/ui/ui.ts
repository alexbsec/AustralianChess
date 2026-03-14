import { canMovePiece } from "../engine/movements";
import type { GameState, Position } from "../engine/types";

const BOARD_SIZE = 12;

export function getPseudoLegalMoves(
    gameState: GameState,
    from: Position,
): Position[] {

    const moves: Position[] = [];

    for (let row = 0; row < BOARD_SIZE; row++) {
        for (let col = 0; col < BOARD_SIZE; col++) {

            if (
                canMovePiece(
                    from.row,
                    from.col,
                    row,
                    col,
                    gameState
                )
            ) {
                moves.push({ row, col });
            }

        }
    }

    return moves;
}
