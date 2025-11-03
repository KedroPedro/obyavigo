document.addEventListener("DOMContentLoaded", function () {
  initTabs();
  initProfileForm();
  initPasswordModal();
  initLogout();
});
function initHeader() {
  const header = document.getElementById("header");
  let lastScrollY = window.scrollY;
  window.addEventListener("scroll", () => {
    if (window.scrollY > lastScrollY && window.scrollY > 100) {
      header.classList.add("hidden");
    } else {
      header.classList.remove("hidden");
    }
    lastScrollY = window.scrollY;
  });
}
function initCategories() {
  const categoriesBtn = document.getElementById("categoriesBtn");
  const categoriesMenu = document.getElementById("categoriesMenu");
  const categoriesData = {
    services: [
      "Ремонт техники",
      "Красота и здоровье",
      "Образование",
      "Транспортные услуги",
      "Ремонт и строительство",
    ],
    auto: [
      "Легковые автомобили",
      "Мотоциклы",
      "Грузовики",
      "Спецтехника",
      "Запчасти и аксессуары",
    ],
    realty: [
      "Квартиры",
      "Дома и коттеджи",
      "Коммерческая недвижимость",
      "Земельные участки",
      "Гаражи и стоянки",
    ],
    electronics: [
      "Смартфоны и планшеты",
      "Ноутбуки и компьютеры",
      "Телевизоры и аудио",
      "Фото и видео",
      "Бытовая техника",
    ],
    work: ["Вакансии", "Резюме", "Фриланс", "Удаленная работа", "Подработка"],
    fashion: ["Одежда", "Обувь", "Аксессуары", "Часы и украшения", "Косметика"],
    home: ["Мебель", "Интерьер", "Посуда", "Текстиль", "Хозяйственные товары"],
    hobby: [
      "Спорт и отдых",
      "Книги и журналы",
      "Коллекционирование",
      "Музыкальные инструменты",
      "Туризм и рыбалка",
    ],
  };
  categoriesMenu.innerHTML = Object.keys(categoriesData)
    .map(
      (category) => `
        <div class="category-item" data-category="${category}">
            <span>${getCategoryName(category)}</span>
            <span class="sub-arrow">›</span>
            <div class="subcategories-menu">
                ${categoriesData[category]
                  .map(
                    (sub) => `
                    <div class="subcategory-item">${sub}</div>
                `,
                  )
                  .join("")}
            </div>
        </div>
    `,
    )
    .join("");
  categoriesBtn?.addEventListener("click", (e) => {
    e.stopPropagation();
    categoriesMenu.classList.toggle("show");
  });
  document.addEventListener("click", () => {
    categoriesMenu.classList.remove("show");
  });
  categoriesMenu?.addEventListener("click", (e) => {
    e.stopPropagation();
  });
}
function getCategoryName(key) {
  const names = {
    services: "Услуги",
    auto: "Авто",
    realty: "Недвижимость",
    electronics: "Электроника",
    work: "Работа",
    fashion: "Мода",
    home: "Для дома",
    hobby: "Хобби",
  };
  return names[key] || key;
}
function initSearch() {
  const searchBtn = document.getElementById("searchBtn");
  const searchOverlay = document.getElementById("searchOverlay");
  const closeSearch = document.getElementById("closeSearch");
  searchBtn?.addEventListener("click", () => {
    searchOverlay.style.display = "flex";
    document.getElementById("searchInput").focus();
    document.body.style.overflow = "hidden";
  });
  closeSearch?.addEventListener("click", () => {
    searchOverlay.style.display = "none";
    document.body.style.overflow = "";
  });
  searchOverlay?.addEventListener("click", (e) => {
    if (e.target === searchOverlay) {
      searchOverlay.style.display = "none";
      document.body.style.overflow = "";
    }
  });
}
function initChat() {
  const messagesBtn = document.getElementById("messagesBtn");
  const chatOverlay = document.getElementById("chatOverlay");
  const closeChatBtn = document.getElementById("closeChatBtn");
  messagesBtn?.addEventListener("click", () => {
    chatOverlay.style.display = "flex";
    document.body.style.overflow = "hidden";
  });
  closeChatBtn?.addEventListener("click", () => {
    chatOverlay.style.display = "none";
    document.body.style.overflow = "";
  });
  chatOverlay?.addEventListener("click", (e) => {
    if (e.target === chatOverlay) {
      chatOverlay.style.display = "none";
      document.body.style.overflow = "";
    }
  });
}
function initLoved() {
  const lovedBtn = document.getElementById("lovedBtn");
  const lovedOverlay = document.getElementById("lovedOverlay");
  const closeLovedBtn = document.getElementById("closeLovedBtn");
  lovedBtn?.addEventListener("click", () => {
    lovedOverlay.style.display = "flex";
    document.body.style.overflow = "hidden";
  });
  closeLovedBtn?.addEventListener("click", () => {
    lovedOverlay.style.display = "none";
    document.body.style.overflow = "";
  });
  lovedOverlay?.addEventListener("click", (e) => {
    if (e.target === lovedOverlay) {
      lovedOverlay.style.display = "none";
      document.body.style.overflow = "";
    }
  });
}
function initMobileMenu() {
  const mobileMenuBtn = document.getElementById("mobileMenuBtn");
  const mobileMenuOverlay = document.getElementById("mobileMenuOverlay");
  const closeMobileMenu = document.getElementById("closeMobileMenu");
  mobileMenuBtn?.addEventListener("click", () => {
    mobileMenuOverlay.style.display = "flex";
    document.body.style.overflow = "hidden";
  });
  closeMobileMenu?.addEventListener("click", () => {
    mobileMenuOverlay.style.display = "none";
    document.body.style.overflow = "";
  });
  mobileMenuOverlay?.addEventListener("click", (e) => {
    if (e.target === mobileMenuOverlay) {
      mobileMenuOverlay.style.display = "none";
      document.body.style.overflow = "";
    }
  });
}
function initProfileMenu() {
  const profileBtn = document.getElementById("profileBtn");
  const profileOverlay = document.getElementById("profileMenuOverlay");
  profileBtn?.addEventListener("click", (e) => {
    e.stopPropagation();
    profileOverlay.style.display = "flex";
    document.body.style.overflow = "hidden";
  });
  profileOverlay?.addEventListener("click", (e) => {
    if (e.target === profileOverlay) {
      profileOverlay.style.display = "none";
      document.body.style.overflow = "";
    }
  });
  document.querySelector(".logout")?.addEventListener("click", (e) => {
    e.preventDefault();
    if (confirm("Вы уверены?")) {
      profileOverlay.style.display = "none";
      document.body.style.overflow = "";
    }
  });
}
function initTheme() {
  const themeToggleBtn = document.getElementById("themeToggleBtn");
  if (!themeToggleBtn) return;
  const savedTheme = localStorage.getItem("theme") || "light";
  document.documentElement.setAttribute("data-theme", savedTheme);
  themeToggleBtn.querySelector(".theme-icon").textContent =
    savedTheme === "dark" ? "☀️" : "🌙";
  themeToggleBtn.addEventListener("click", () => {
    const current = document.documentElement.getAttribute("data-theme");
    const newTheme = current === "dark" ? "light" : "dark";
    document.documentElement.setAttribute("data-theme", newTheme);
    localStorage.setItem("theme", newTheme);
    themeToggleBtn.querySelector(".theme-icon").textContent =
      newTheme === "dark" ? "☀️" : "🌙";
  });
}
function initTabs() {
  const tabButtons = document.querySelectorAll(".menu-item");
  tabButtons.forEach((button) => {
    button.addEventListener("click", () => {
      tabButtons.forEach((btn) => btn.classList.remove("active"));
      document
        .querySelectorAll(".tab-content")
        .forEach((tab) => tab.classList.remove("active"));
      button.classList.add("active");
      const tabId = button.dataset.tab + "-tab";
      document.getElementById(tabId).classList.add("active");
      if (tabId === "ads-tab") {
        loadUserAds();
      }
      if (tabId === "loved-tab") {
        loadLikedAds();
      }
    });
  });
}
function initProfileForm() {
  const form = document.getElementById("profileForm");
  if (!form) return;
  form.addEventListener("submit", async function (e) {
    e.preventDefault();
    showError("nameError", "");
    showError("phoneError", "");
    const username = document.getElementById("userName").value.trim();
    const phone = document.getElementById("userPhone").value.trim();
    if (!username) {
      showError("nameError", "Имя не может быть пустым");
      return;
    }
    if (phone && !isValidBelarusPhone(phone)) {
      showError("phoneError", "Неверный формат белорусского номера");
      return;
    }
    const data = {
      username: username,
      phone_number: phone || null
    };
    try {
      const response = await fetch("/api/profile/update/", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });
      const result = await response.json();
      if (response.ok) {
        alert("Профиль успешно обновлён");
        const profileName = document.getElementById("profileName");
        if (profileName) {
          profileName.textContent = username;
        }
      } else {
        showError("nameError", result.message || "Ошибка при обновлении профиля");
      }
    } catch (error) {
      console.error("Ошибка при обновлении профиля:", error);
      showError("nameError", "Ошибка подключения к серверу");
    }
  });
}
function isValidBelarusPhone(phone) {
  if (!phone) return true;
  const cleaned = phone.replace(/\D/g, "");
  if (cleaned.length === 12 && cleaned.startsWith("375")) {
    return /^375(25|29|33|44)\d{7}$/.test(cleaned);
  }
  if (cleaned.length === 9) {
    return /^(25|29|33|44)\d{7}$/.test(cleaned);
  }
  return false;
}
function showError(elementId, message) {
  const el = document.getElementById(elementId);
  if (el) el.textContent = message;
}
function initPasswordModal() {
  const modal = document.getElementById("passwordModal");
  if (!modal) return;
  document
    .getElementById("changePasswordBtn")
    ?.addEventListener("click", () => {
      modal.style.display = "flex";
    });
  document
    .getElementById("closePasswordModal")
    ?.addEventListener("click", () => {
      modal.style.display = "none";
    });
  modal.addEventListener("click", (e) => {
    if (e.target === modal) modal.style.display = "none";
  });
  document.getElementById("savePasswordBtn")?.addEventListener("click", async () => {
    const oldPass = document.getElementById("oldPassword")?.value || "";
    const newPass = document.getElementById("newPassword").value;
    const confirmPass = document.getElementById("confirmPassword").value;
    showError("oldPasswordError", "");
    showError("newPasswordError", "");
    showError("confirmPasswordError", "");
    if (newPass !== confirmPass) {
      showError("confirmPasswordError", "Пароли не совпадают");
      return;
    }
    if (!isValidPassword(newPass)) {
      showError(
        "newPasswordError",
        "Пароль: минимум 8 символов, буквы и цифры",
      );
      return;
    }
    try {
      const response = await fetch("/api/profile/change-password/", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          old_password: oldPass,
          new_password: newPass,
        }),
      });
      const data = await response.json();
      if (response.ok) {
        alert("Пароль успешно изменён");
        modal.style.display = "none";
        if (document.getElementById("oldPassword")) document.getElementById("oldPassword").value = "";
        document.getElementById("newPassword").value = "";
        document.getElementById("confirmPassword").value = "";
      } else {
        showError("newPasswordError", data.message || "Ошибка при изменении пароля");
      }
    } catch (error) {
      console.error("Ошибка при смене пароля:", error);
      showError("newPasswordError", "Ошибка подключения к серверу");
    }
  });
}
function isValidPassword(password) {
  return (
    password.length >= 8 && /[a-zA-Z]/.test(password) && /\d/.test(password)
  );
}
function initLogout() {
  document.querySelector(".logout")?.addEventListener("click", (e) => {
    e.preventDefault();
    if (confirm("Выйти из аккаунта?")) {
      window.location.href = "/auth.html";
    }
  });
}
document.getElementById("deleteAccountBtn")?.addEventListener("click", async () => {
  if (confirm("Вы уверены, что хотите удалить аккаунт. Это действие нельзя отменить")) {
    try {
      const response = await fetch("/api/profile/delete-account/", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (response.ok) {
        alert("Аккаунт удалён");
        window.location.href = "/auth/";
      } else {
        const data = await response.json();
        alert(data.message || "Ошибка при удалении аккаунта");
      }
    } catch (error) {
      console.error("Ошибка при удалении аккаунта:", error);
      alert("Ошибка подключения к серверу");
    }
  }
});
function deleteAllCookies() {
  document.cookie.split(";").forEach((cookie) => {
    const eqPos = cookie.indexOf("=");
    const name = eqPos > -1 ? cookie.substr(0, eqPos).trim() : cookie.trim();
    document.cookie = name + "=; Max-Age=-1; path=/;";
  });
}
document.getElementById("logoutBtn").onclick = async () => {
  try {
    const response = await fetch("/api/logout/", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
    });
    if (response.ok) {
      window.location.href = "/auth/";
    } else {
      alert("Ошибка выхода из аккаунта");
    }
  } catch (e) {
    console.error("Ошибка запроса выхода", e);
  }
};
let userAdsLoaded = false;
let likedAdsLoaded = false;
async function loadUserAds() {
  if (userAdsLoaded) return;
  const adsList = document.getElementById("userAdsList");
  try {
    const response = await fetch("/api/user/ads/", {
      credentials: "include",
    });
    if (!response.ok) {
      throw new Error("Failed to load ads");
    }
    const data = await response.json();
    if (data.success && data.data && data.data.length > 0) {
      adsList.innerHTML = data.data.map(ad => `
        <div class="ad-card" onclick="window.location.href='/ads/${ad.adID}/'" style="cursor: pointer;">
          <div class="ad-image">
            ${ad.imageID ?
              `<img src="/api/images/${ad.imageID}/" alt="${ad.title}" loading="lazy" />` :
              `<img src="/static/pictures/logo.png" alt="${ad.title}" loading="lazy" />`
            }
            <span class="ad-status ${ad.status}">${getStatusText(ad.status)}</span>
          </div>
          <div class="ad-content">
            <h3 class="ad-title">${ad.title}</h3>
            <p class="ad-description">${ad.desc}</p>
            <div class="ad-meta">
              <span class="ad-price">${ad.price} BYN</span>
              <span class="ad-location">📍 ${ad.locationName}</span>
              <span class="ad-category">${ad.categoryName}</span>
            </div>
            <div class="ad-footer">
              <span class="ad-date">${new Date(ad.createdAt).toLocaleDateString('ru-RU')}</span>
              <span class="ad-views">👁️ ${ad.viewsCount}</span>
            </div>
          </div>
        </div>
      `).join("");
    } else {
      adsList.innerHTML = `
        <div class="no-ads">
          <div class="no-ads-content">
            <div class="no-ads-icon">📦</div>
            <h3>У вас пока нет объявлений</h3>
            <p>Создайте свое первое объявление</p>
            <a href="/create-ad/" class="btn-primary">Создать объявление</a>
          </div>
        </div>
      `;
    }
    userAdsLoaded = true;
  } catch (error) {
    console.error("Error loading user ads:", error);
    adsList.innerHTML = `
      <div class="no-ads">
        <div class="no-ads-content">
          <div class="no-ads-icon">⚠️</div>
          <h3>Ошибка загрузки объявлений</h3>
          <p>Попробуйте обновить страницу</p>
        </div>
      </div>
    `;
  }
}
async function loadLikedAds() {
  if (likedAdsLoaded) return;
  const adsList = document.getElementById("lovedAdsList");
  try {
    const response = await fetch("/api/ads/?limit=100", {
      credentials: "include",
    });
    if (!response.ok) {
      throw new Error("Failed to load liked ads");
    }
    const data = await response.json();
    adsList.innerHTML = `
      <div class="no-ads">
        <div class="no-ads-content">
          <div class="no-ads-icon">❤️</div>
          <h3>Избранные объявления</h3>
          <p>Функция в разработке. Используйте страницу "Избранное" из меню</p>
        </div>
      </div>
    `;
    likedAdsLoaded = true;
  } catch (error) {
    console.error("Error loading liked ads:", error);
    adsList.innerHTML = `
      <div class="no-ads">
        <div class="no-ads-content">
          <div class="no-ads-icon">⚠️</div>
          <h3>Ошибка загрузки избранного</h3>
          <p>Попробуйте обновить страницу</p>
        </div>
      </div>
    `;
  }
}
function getStatusText(status) {
  const statusMap = {
    'public': 'Активно',
    'moderation': 'На модерации',
    'draft': 'Черновик',
    'archived': 'В архиве',
    'blocked': 'Заблокировано'
  };
  return statusMap[status] || status;
}
