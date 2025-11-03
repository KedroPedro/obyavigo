class ResetPasswordManager {
  constructor() {
    this.token = null;
    this.init();
  }

  init() {
    document.addEventListener("DOMContentLoaded", () => {
      this.initTheme();
      this.extractToken();
      this.initForm();
    });
  }

  initTheme() {
    const themeToggleBtn = document.getElementById("themeToggleBtn");
    if (!themeToggleBtn) return;
    const themeIcon = themeToggleBtn.querySelector(".theme-icon");
    const savedTheme = localStorage.getItem("theme") || "light";
    document.documentElement.setAttribute("data-theme", savedTheme);
    if (themeIcon) {
      themeIcon.textContent = savedTheme === "dark" ? "☀️" : "🌙";
    }
    themeToggleBtn.addEventListener("click", () => {
      const currentTheme = document.documentElement.getAttribute("data-theme");
      const newTheme = currentTheme === "dark" ? "light" : "dark";
      document.documentElement.setAttribute("data-theme", newTheme);
      localStorage.setItem("theme", newTheme);
      if (themeIcon) {
        themeIcon.textContent = newTheme === "dark" ? "☀️" : "🌙";
      }
    });
  }

  extractToken() {
    const urlParams = new URLSearchParams(window.location.search);
    this.token = urlParams.get("token");
    if (!this.token) {
      alert("Неверная ссылка для восстановления пароля");
      window.location.href = "/auth/";
    }
  }

  initForm() {
    const form = document.getElementById("resetPasswordForm");
    if (form) {
      form.addEventListener("submit", (e) => this.handleSubmit(e));
    }
  }

  isValidPassword(password) {
    return (
      password.length >= 8 && /[a-zA-Z]/.test(password) && /\d/.test(password)
    );
  }

  clearErrors() {
    document.getElementById("newPasswordError").textContent = "";
    document.getElementById("confirmPasswordError").textContent = "";
    document.getElementById("newPassword").classList.remove("invalid", "valid");
    document
      .getElementById("confirmPassword")
      .classList.remove("invalid", "valid");
  }

  showError(elementId, message) {
    const element = document.getElementById(elementId);
    if (element) {
      element.textContent = message;
    }
  }

  handleSubmit(e) {
    e.preventDefault();
    this.clearErrors();

    const newPasswordInput = document.getElementById("newPassword");
    const confirmPasswordInput = document.getElementById("confirmPassword");

    const newPassword = newPasswordInput.value;
    const confirmPassword = confirmPasswordInput.value;

    let isValid = true;

    if (!newPassword) {
      this.showError("newPasswordError", "Введите пароль");
      newPasswordInput.classList.add("invalid");
      isValid = false;
    } else if (!this.isValidPassword(newPassword)) {
      this.showError(
        "newPasswordError",
        "Пароль: минимум 8 символов, буквы и цифры",
      );
      newPasswordInput.classList.add("invalid");
      isValid = false;
    } else {
      newPasswordInput.classList.add("valid");
    }

    if (!confirmPassword) {
      this.showError("confirmPasswordError", "Повторите пароль");
      confirmPasswordInput.classList.add("invalid");
      isValid = false;
    } else if (newPassword !== confirmPassword) {
      this.showError("confirmPasswordError", "Пароли не совпадают");
      confirmPasswordInput.classList.add("invalid");
      isValid = false;
    } else {
      confirmPasswordInput.classList.add("valid");
    }

    if (isValid) {
      this.resetPassword(newPassword);
    }
  }

  async resetPassword(newPassword) {
    try {
      const response = await fetch("/api/auth/reset-password/", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          token: this.token,
          new_password: newPassword,
        }),
      });

      const result = await response.json();

      if (response.ok) {
        alert("Пароль успешно изменен. Теперь вы можете войти с новым паролем");
        window.location.href = "/auth/";
      } else {
        const errorMessage =
          result.error || result.message || "Неизвестная ошибка";
        alert("Ошибка: " + errorMessage);
      }
    } catch (error) {
      console.error("Ошибка сети:", error);
      alert("Не удалось подключиться к серверу. Проверьте соединение.");
    }
  }
}

new ResetPasswordManager();
