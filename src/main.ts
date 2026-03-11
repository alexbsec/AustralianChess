import "./style.css";
import { InitializeBoard } from "./board/game";

const app = document.querySelector<HTMLDivElement>("#app");

if (!app) {
    throw new Error("App container not found.");
}

const boardElement = document.createElement("div");
boardElement.classList.add("board");
app.appendChild(boardElement);

InitializeBoard(boardElement);