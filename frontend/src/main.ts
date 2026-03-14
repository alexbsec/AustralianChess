import "./index-style.css";

import { navigateTo } from "./router";
import { renderLandingPage } from "./pages/landing";
import { renderRegisterPage } from "./pages/register";
import { renderLoginPage } from "./pages/login";
import { renderRoomPage } from "./pages/room";

const app = document.querySelector<HTMLDivElement>('#app')!; // Added !

if (!app) {
    throw new Error("App container not found.");
}

function router(): void {
    const path = window.location.pathname;

    app.innerHTML = "";

    switch (path) {
        case "/":
            renderLandingPage(app);
            return;
        case "/register":
            renderRegisterPage(app);
            return;
        case "/login":
            renderLoginPage(app);
            return;
        case "/room":
            renderRoomPage(app);
            return;
        default:
            renderNotFoundPage(app);
    }
}

function renderNotFoundPage(container: HTMLDivElement): void {
    const wrapper = document.createElement("main");
    wrapper.className = "auth-page";

    const card = document.createElement("section");
    card.className = "auth-card";

    const title = document.createElement("h1");
    title.className = "auth-title";
    title.textContent = "Page not found";

    const text = document.createElement("p");
    text.className = "auth-subtitle";
    text.textContent = "That route does not exist.";

    const button = document.createElement("button");
    button.className = "btn btn-primary";
    button.textContent = "Go Home";
    button.addEventListener("click", () => navigateTo("/"));

    card.appendChild(title);
    card.appendChild(text);
    card.appendChild(button);
    wrapper.appendChild(card);
    container.appendChild(wrapper);
}

window.addEventListener("popstate", router);
window.addEventListener("app:navigate", router);

router();
