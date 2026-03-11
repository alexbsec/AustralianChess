import "./style.css";
import { InitializeGame } from "./game/game";

const app = document.querySelector<HTMLDivElement>("#app");

if (!app) {
    throw new Error("App container not found.");
}

const boardElement = document.createElement("div");
boardElement.classList.add("board");
app.appendChild(boardElement);

document.addEventListener("click", (event) => {
    console.log("DOCUMENT CLICK", event.target);
}, true);

InitializeGame(boardElement);