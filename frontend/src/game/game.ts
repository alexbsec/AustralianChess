import { GameController } from "./controller";

export function InitializeGame(container: HTMLElement): void {
    let gameController = new GameController(container);
    gameController.start();
}