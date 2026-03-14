import "../index-style.css";

type RegisterSuccessResponse = {
    message: string;
    player_id: string;
};

type RegisterErrorResponse = {
    error: string;
};

const app = document.querySelector<HTMLDivElement>("#app");

if (!app) {
    throw new Error("App container not found.");
}

function getPasswordStrength(password: string): {
    label: string;
    className: string;
    width: string;
} {
    if (password.length === 0) {
        return {
            label: "Enter a password",
            className: "strength-empty",
            width: "0%",
        };
    }

    if (password.length < 8) {
        return {
            label: "Too short",
            className: "strength-weak",
            width: "35%",
        };
    }

    return {
        label: "Valid length",
        className: "strength-strong",
        width: "100%",
    };
}

export function renderRegisterPage(container: HTMLDivElement): void {
    container.innerHTML = "";

    const page = document.createElement("main");
    page.className = "auth-page";

    const card = document.createElement("section");
    card.className = "auth-card";

    const backLink = document.createElement("a");
    backLink.className = "back-link";
    backLink.href = "/";
    backLink.textContent = "← Back to home";

    const title = document.createElement("h1");
    title.className = "auth-title";
    title.textContent = "Create Account";

    const subtitle = document.createElement("p");
    subtitle.className = "auth-subtitle";
    subtitle.textContent =
        "Claim your name, enter the board, and stop losing as a guest.";

    const form = document.createElement("form");
    form.className = "auth-form";

    const usernameGroup = document.createElement("div");
    usernameGroup.className = "form-group";

    const usernameLabel = document.createElement("label");
    usernameLabel.htmlFor = "username";
    usernameLabel.textContent = "Username";

    const usernameInput = document.createElement("input");
    usernameInput.id = "username";
    usernameInput.name = "username";
    usernameInput.type = "text";
    usernameInput.placeholder = "Enter your username";
    usernameInput.autocomplete = "username";
    usernameInput.required = true;
    usernameInput.className = "form-input";

    usernameGroup.appendChild(usernameLabel);
    usernameGroup.appendChild(usernameInput);

    const passwordGroup = document.createElement("div");
    passwordGroup.className = "form-group";

    const passwordLabel = document.createElement("label");
    passwordLabel.htmlFor = "password";
    passwordLabel.textContent = "Password";

    const passwordInput = document.createElement("input");
    passwordInput.id = "password";
    passwordInput.name = "password";
    passwordInput.type = "password";
    passwordInput.placeholder = "Enter your password";
    passwordInput.autocomplete = "new-password";
    passwordInput.required = true;
    passwordInput.className = "form-input";

    const strengthWrapper = document.createElement("div");
    strengthWrapper.className = "password-strength-wrapper";

    const strengthBarBackground = document.createElement("div");
    strengthBarBackground.className = "password-strength-bar-bg";

    const strengthBarFill = document.createElement("div");
    strengthBarFill.className = "password-strength-bar-fill strength-empty";
    strengthBarFill.style.width = "0%";

    const strengthText = document.createElement("p");
    strengthText.className = "password-strength-text";
    strengthText.textContent = "Enter a password";

    strengthBarBackground.appendChild(strengthBarFill);
    strengthWrapper.appendChild(strengthBarBackground);
    strengthWrapper.appendChild(strengthText);

    passwordGroup.appendChild(passwordLabel);
    passwordGroup.appendChild(passwordInput);
    passwordGroup.appendChild(strengthWrapper);

    const confirmPasswordGroup = document.createElement("div");
    confirmPasswordGroup.className = "form-group";

    const confirmPasswordLabel = document.createElement("label");
    confirmPasswordLabel.htmlFor = "confirm-password";
    confirmPasswordLabel.textContent = "Confirm Password";

    const confirmPasswordInput = document.createElement("input");
    confirmPasswordInput.id = "confirm-password";
    confirmPasswordInput.name = "confirm-password";
    confirmPasswordInput.type = "password";
    confirmPasswordInput.placeholder = "Confirm your password";
    confirmPasswordInput.autocomplete = "new-password";
    confirmPasswordInput.required = true;
    confirmPasswordInput.className = "form-input";

    confirmPasswordGroup.appendChild(confirmPasswordLabel);
    confirmPasswordGroup.appendChild(confirmPasswordInput);

    const feedback = document.createElement("div");
    feedback.className = "form-feedback hidden";

    const submitButton = document.createElement("button");
    submitButton.type = "submit";
    submitButton.className = "btn btn-primary auth-submit";
    submitButton.textContent = "Create Account";

    const loginHint = document.createElement("p");
    loginHint.className = "auth-footer-text";
    loginHint.innerHTML = `Already have an account? <a href="/login.html">Login</a>`;

    function setFeedback(message: string, type: "error" | "success"): void {
        feedback.className = `form-feedback ${type}`;
        feedback.textContent = message;
    }

    function clearFeedback(): void {
        feedback.className = "form-feedback hidden";
        feedback.textContent = "";
    }

    function updateStrengthBar(): void {
        const strength = getPasswordStrength(passwordInput.value);

        strengthBarFill.className = `password-strength-bar-fill ${strength.className}`;
        strengthBarFill.style.width = strength.width;
        strengthText.textContent = strength.label;
    }

    passwordInput.addEventListener("input", updateStrengthBar);

    form.addEventListener("submit", async (event) => {
        event.preventDefault();
        clearFeedback();

        const username = usernameInput.value.trim();
        const password = passwordInput.value;
        const confirmPassword = confirmPasswordInput.value;

        if (username.length === 0) {
            setFeedback("Username is required.", "error");
            return;
        }

        if (password.length < 8) {
            setFeedback("Password must be at least 8 characters long.", "error");
            return;
        }

        if (password !== confirmPassword) {
            setFeedback("Passwords do not match.", "error");
            return;
        }

        submitButton.disabled = true;
        submitButton.textContent = "Creating account...";

        try {
            const response = await fetch("/api/v1/user/register", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    username,
                    password,
                    confirm_password: confirmPassword,
                }),
            });

            const data = (await response.json()) as
                | RegisterSuccessResponse
                | RegisterErrorResponse;

            if (!response.ok) {
                const errorMessage =
                    "error" in data ? data.error : "Could not create account.";
                setFeedback(errorMessage, "error");
                return;
            }

            if (!("message" in data) || !("player_id" in data)) {
                setFeedback("Unexpected server response.", "error");
                return;
            }

            setFeedback(
                `${data.message} Welcome, ${data.player_id}. You can now log in.`,
                "success",
            );

            form.reset();
            updateStrengthBar();

            setTimeout(() => {
                window.location.href = `/login.html?username=${encodeURIComponent(data.player_id)}`;
            }, 1400);
        } catch (error) {
            console.error(error);
            setFeedback("Network error. Please try again.", "error");
        } finally {
            submitButton.disabled = false;
            submitButton.textContent = "Create Account";
        }
    });

    form.appendChild(usernameGroup);
    form.appendChild(passwordGroup);
    form.appendChild(confirmPasswordGroup);
    form.appendChild(feedback);
    form.appendChild(submitButton);

    card.appendChild(backLink);
    card.appendChild(title);
    card.appendChild(subtitle);
    card.appendChild(form);
    card.appendChild(loginHint);

    page.appendChild(card);
    container.appendChild(page);

    updateStrengthBar();
}

renderRegisterPage(app);
