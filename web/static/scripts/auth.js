class AuthManager {
  constructor() {
    this.currentForm = "login";
    this.init();
  }

  init() {
    document.addEventListener("DOMContentLoaded", () => {
      this.initTheme();
      this.initAuthForms();
      this.initInputRestrictions();
      this.setupKeyboardNavigation();
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

  initInputRestrictions() {
    const emailInputs = document.querySelectorAll(
      'input[type="email"], #loginEmail, #regEmail, #forgotEmail',
    );
    emailInputs.forEach((input) => {
      input.addEventListener("keypress", (e) => {
        const char = e.key;
        if (!/^[a-zA-Z0-9@._%+-]$/.test(char) && !this.isControlKey(e)) {
          e.preventDefault();
        }
      });
    });
  }

  setupKeyboardNavigation() {
    document.addEventListener("keydown", (e) => {
      if (e.key === "Escape") {
        this.clearAllErrors();
      }
    });
  }

  isControlKey(e) {
    return (
      e.ctrlKey ||
      e.metaKey ||
      e.altKey ||
      [
        "Backspace",
        "Delete",
        "ArrowLeft",
        "ArrowRight",
        "Tab",
        "Enter",
      ].includes(e.key)
    );
  }
  hasCyrillic(str) {
    return /[а-яА-ЯёЁ]/.test(str);
  }

  isValidName(name) {
    if (!name) return false;
    const re = /^[a-zA-Zа-яА-ЯёЁ\s\-']+$/;
    return re.test(name) && name.trim().length >= 2;
  }

  isValidEmail(email) {
    if (!email) return false;
    if (this.hasCyrillic(email)) return false;
    const re = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    return re.test(email);
  }

  isValidBelarusPhone(phone) {
    if (!phone) return true;
    if (this.hasCyrillic(phone)) return false;
    const cleaned = phone.replace(/\D/g, "");
    if (cleaned.length === 12 && cleaned.startsWith("375")) {
      return /^375(25|29|33|44)\d{7}$/.test(cleaned);
    }
    if (cleaned.length === 9) {
      return /^(25|29|33|44)\d{7}$/.test(cleaned);
    }
    return false;
  }

  isValidPassword(password) {
    return (
      password.length >= 8 && /[a-zA-Z]/.test(password) && /\d/.test(password)
    );
  }

  clearErrors(formId) {
    const form = document.getElementById(formId);
    if (!form) return;
    const errorElements = form.querySelectorAll(".error-message");
    errorElements.forEach((el) => (el.textContent = ""));
    const inputs = form.querySelectorAll("input");
    inputs.forEach((input) => {
      input.classList.remove("invalid", "valid");
    });
  }

  clearAllErrors() {
    ["loginForm", "registerForm", "forgotPasswordForm"].forEach((formId) => {
      this.clearErrors(formId);
    });
  }

  showError(elementId, message) {
    const element = document.getElementById(elementId);
    if (element) {
      element.textContent = message;
    }
  }

  initAuthForms() {
    const loginBtn = document.getElementById("loginBtn");
    const registerBtn = document.getElementById("registerBtn");
    const loginForm = document.getElementById("loginForm");
    const registerForm = document.getElementById("registerForm");
    const forgotPasswordForm = document.getElementById("forgotPasswordForm");
    const forgotPasswordLink = document.getElementById("forgotPasswordLink");
    const backToLoginBtn = document.getElementById("backToLoginBtn");

    if (loginBtn) {
      loginBtn.addEventListener("click", () => this.switchToForm("login"));
    }

    if (registerBtn) {
      registerBtn.addEventListener("click", () =>
        this.switchToForm("register"),
      );
    }

    if (forgotPasswordLink) {
      forgotPasswordLink.addEventListener("click", (e) => {
        e.preventDefault();
        this.switchToForm("forgot");
      });
    }

    if (backToLoginBtn) {
      backToLoginBtn.addEventListener("click", () =>
        this.switchToForm("login"),
      );
    }

    if (loginForm) {
      loginForm.addEventListener("submit", (e) => this.handleLogin(e));
    }

    if (registerForm) {
      registerForm.addEventListener("submit", (e) => this.handleRegister(e));
    }

    if (forgotPasswordForm) {
      forgotPasswordForm.addEventListener("submit", (e) =>
        this.handleForgotPassword(e),
      );
    }
  }

  switchToForm(formType) {
    const loginBtn = document.getElementById("loginBtn");
    const registerBtn = document.getElementById("registerBtn");
    const loginForm = document.getElementById("loginForm");
    const registerForm = document.getElementById("registerForm");
    const forgotPasswordForm = document.getElementById("forgotPasswordForm");
    [loginForm, registerForm, forgotPasswordForm].forEach((form) => {
      if (form) form.classList.remove("active");
    });
    [loginBtn, registerBtn].forEach((btn) => {
      if (btn) btn.classList.remove("active");
    });
    switch (formType) {
      case "login":
        if (loginBtn) loginBtn.classList.add("active");
        if (loginForm) loginForm.classList.add("active");
        if (loginBtn) loginBtn.setAttribute("aria-selected", "true");
        if (registerBtn) registerBtn.setAttribute("aria-selected", "false");
        break;
      case "register":
        if (registerBtn) registerBtn.classList.add("active");
        if (registerForm) registerForm.classList.add("active");
        if (registerBtn) registerBtn.setAttribute("aria-selected", "true");
        if (loginBtn) loginBtn.setAttribute("aria-selected", "false");
        break;
      case "forgot":
        if (forgotPasswordForm) forgotPasswordForm.classList.add("active");
        break;
    }

    this.clearAllErrors();
  }

  handleLogin(e) {
    e.preventDefault();
    this.clearErrors("loginForm");

    const emailInput = document.getElementById("loginEmail");
    const passwordInput = document.getElementById("loginPassword");
    let isValid = true;

    const email = emailInput.value.trim();
    const password = passwordInput.value;

    if (!email) {
      this.showError("loginEmailError", "Введите email или телефон");
      emailInput.classList.add("invalid");
      isValid = false;
    } else if (this.hasCyrillic(email)) {
      this.showError(
        "loginEmailError",
        "Email или телефон не должен содержать русские буквы",
      );
      emailInput.classList.add("invalid");
      isValid = false;
    } else if (this.isValidEmail(email)) {
      emailInput.classList.add("valid");
    } else if (this.isValidBelarusPhone(email)) {
      emailInput.classList.add("valid");
    } else {
      this.showError("loginEmailError", "Неверный email или номер телефона");
      emailInput.classList.add("invalid");
      isValid = false;
    }

    if (!password) {
      this.showError("loginPasswordError", "Введите пароль");
      passwordInput.classList.add("invalid");
      isValid = false;
    } else if (password.length < 6) {
      this.showError(
        "loginPasswordError",
        "Пароль должен быть не менее 6 символов",
      );
      passwordInput.classList.add("invalid");
      isValid = false;
    } else {
      passwordInput.classList.add("valid");
    }

    if (isValid) {
      this.loginUser(email, password);
    }
  }

  async loginUser(email, password) {
    try {
      const response = await fetch("/api/auth/login/", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: email,
          password: password,
        }),
      });

      const passwordInput = document.getElementById("loginPassword");
      const emailInput = document.getElementById("loginEmail");

      const result = await response.json();
      if (response.ok) {
        alert("Вход выполнен успешно");
        window.location.href = "/";
      } else if (response.status === 401) {
        alert("Ошибка: неверный логин или пароль");
        passwordInput.classList.remove("valid");
        passwordInput.classList.add("invalid");
        emailInput.classList.remove("valid");
        emailInput.classList.add("invalid");
        passwordInput.value = "";
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

  handleRegister(e) {
    e.preventDefault();
    this.clearErrors("registerForm");

    const nameInput = document.getElementById("regName");
    const emailInput = document.getElementById("regEmail");
    const passwordInput = document.getElementById("regPassword");
    const passwordConfirmInput = document.getElementById("regPasswordConfirm");
    const termsCheckbox = document.getElementById("agreeTerms");

    let isValid = true;

    const name = nameInput.value.trim();
    if (!name) {
      this.showError("regNameError", "Введите ваше имя");
      nameInput.classList.add("invalid");
      isValid = false;
    } else if (!this.isValidName(name)) {
      this.showError(
        "regNameError",
        "Имя должно содержать только буквы, пробелы, дефисы или апострофы",
      );
      nameInput.classList.add("invalid");
      isValid = false;
    } else {
      nameInput.classList.add("valid");
    }

    const email = emailInput.value.trim();
    if (!email) {
      this.showError("regEmailError", "Введите email");
      emailInput.classList.add("invalid");
      isValid = false;
    } else if (this.hasCyrillic(email)) {
      this.showError("regEmailError", "Email должен быть на английском языке");
      emailInput.classList.add("invalid");
      isValid = false;
    } else if (!this.isValidEmail(email)) {
      this.showError("regEmailError", "Неверный формат email");
      emailInput.classList.add("invalid");
      isValid = false;
    } else {
      emailInput.classList.add("valid");
    }

    const password = passwordInput.value;
    if (!password) {
      this.showError("regPasswordError", "Введите пароль");
      passwordInput.classList.add("invalid");
      isValid = false;
    } else if (!this.isValidPassword(password)) {
      this.showError(
        "regPasswordError",
        "Пароль: минимум 8 символов, буквы и цифры",
      );
      passwordInput.classList.add("invalid");
      isValid = false;
    } else {
      passwordInput.classList.add("valid");
    }

    const passwordConfirm = passwordConfirmInput.value;
    if (!passwordConfirm) {
      this.showError("regPasswordConfirmError", "Повторите пароль");
      passwordConfirmInput.classList.add("invalid");
      isValid = false;
    } else if (password !== passwordConfirm) {
      this.showError("regPasswordConfirmError", "Пароли не совпадают");
      passwordConfirmInput.classList.add("invalid");
      isValid = false;
    } else {
      passwordConfirmInput.classList.add("valid");
    }

    if (!termsCheckbox || !termsCheckbox.checked) {
      this.showError(
        "regPasswordConfirmError",
        "Необходимо согласиться с условиями",
      );
      isValid = false;
    }

    if (isValid) {
      this.registerUser(name, email, password);
    }
  }

  async registerUser(name, email, password) {
    try {
      const response = await fetch("/api/auth/register/", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username: name,
          email: email,
          password: password,
        }),
      });

      const result = await response.json();
      const emailInput = document.getElementById("regEmail");

      if (response.ok) {
        alert("Регистрация успешна");
        window.location.href = "/";
      } else if (response.status === 409) {
        this.showError("regEmailError", "Ошибка: email уже зарегистрирован");
        emailInput.classList.remove("valid");
        emailInput.classList.add("invalid");
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

  handleForgotPassword(e) {
    e.preventDefault();
    this.clearErrors("forgotPasswordForm");

    const emailInput = document.getElementById("forgotEmail");
    const email = emailInput.value.trim();

    if (!email) {
      this.showError("forgotEmailError", "Введите email");
      emailInput.classList.add("invalid");
    } else if (this.hasCyrillic(email)) {
      this.showError(
        "forgotEmailError",
        "Email должен быть на английском языке",
      );
      emailInput.classList.add("invalid");
    } else if (!this.isValidEmail(email)) {
      this.showError("forgotEmailError", "Неверный формат email");
      emailInput.classList.add("invalid");
    } else {
      emailInput.classList.add("valid");
      this.forgotPassword(email);
    }
  }

  forgotPassword(email) {
    console.log("Восстановление пароля:", { email });
    alert("Ссылка для восстановления пароля отправлена на ваш email");
  }
}
new AuthManager();
