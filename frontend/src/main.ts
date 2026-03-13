import "./style.css";

const app = document.querySelector<HTMLDivElement>("#app");

if (!app) {
    throw new Error("App container not found.");
}

function renderLobby(container: HTMLDivElement): void {
    const title = document.createElement("h1");
    title.textContent = "Australian Chess";

    const createButton = document.createElement("button");
    createButton.textContent = "Create Room";

    createButton.addEventListener("click", async () => {
        try {
            const response = await fetch("/api/v1/room/create", {
                method: "GET",
            });

            if (!response.ok) {
                throw new Error("Failed to create room");
            }

            const data = await response.json();
            const roomId = data.roomId;

            window.location.href =
                `/room.html?roomId=${encodeURIComponent(roomId)}&playerId=noob`;

        } catch (error) {
            console.error(error);
            alert("Could not create room");
        }
    });

    container.appendChild(title);
    container.appendChild(createButton);
}

renderLobby(app);
