document.addEventListener("DOMContentLoaded", function () {
  initHeader();
  initCategories();
  initSearch();
  initProfileMenu();

  initLoved();
  initMobileMenu();
  initTheme();

  loadDynamicContent();
});

function loadDynamicContent() {
  showSkeletons();
  setTimeout(() => {
    fetchStats();
    fetchCategories();
    fetchFeaturedAds();
    fetchLovedItems(); // Имитация избранного
  }, 1000);
}

function showSkeletons() {
  const categoriesGrid = document.getElementById("categoriesGrid");
  const adsGrid = document.getElementById("adsGrid");

  categoriesGrid.innerHTML = `
        <div class="skeleton-card"></div>
        <div class="skeleton-card"></div>
        <div class="skeleton-card"></div>
        <div class="skeleton-card"></div>
    `;

  adsGrid.innerHTML = `
        <div class="skeleton-ad"></div>
        <div class="skeleton-ad"></div>
        <div class="skeleton-ad"></div>
    `;
}

function fetchStats() {
  const stats = {
    yearsOnMarket: "15 лет",
    citiesCovered: "85 городов",
    appRating: "4.9/5",
    dailyDeals: "3 200+",
  };

  document.getElementById("years-on-market").textContent = stats.yearsOnMarket;
  document.getElementById("cities-covered").textContent = stats.citiesCovered;
  document.getElementById("app-rating").textContent = stats.appRating;
  document.getElementById("daily-deals").textContent = stats.dailyDeals;
}

function fetchCategories() {
  const categoriesGrid = document.getElementById("categoriesGrid");
  const categories = [
    { name: "Автомобили", image: "auto" },
    { name: "Недвижимость", image: "realty" },
    { name: "Электроника", image: "electronics" },
    { name: "Услуги", image: "services" },
    { name: "Работа", image: "work" },
    { name: "Мода", image: "fashion" },
    { name: "Для дома", image: "home" },
    { name: "Хобби", image: "hobby" },
  ];

  categoriesGrid.innerHTML = categories
    .map(
      (cat) => `
        <div class="category-card">
            <div class="category-image">
                <img src="./pictures/category-${cat.image}.jpg" alt="${cat.name}" onerror="this.src='./pictures/placeholder.jpg'">
            </div>
            <div class="category-content">
                <h3>${cat.name}</h3>
                <p>Загружается...</p>
            </div>
        </div>
    `,
    )
    .join("");
}

function fetchFeaturedAds() {
  const adsGrid = document.getElementById("adsGrid");
  adsGrid.innerHTML = `
        <div class="ad-card">
            <h4>Загрузка объявлений...</h4>
            <p>Данные подгружаются с сервера</p>
        </div>
        <div class="ad-card">
            <h4>Загрузка объявлений...</h4>
            <p>Данные подгружаются с сервера</p>
        </div>
        <div class="ad-card">
            <h4>Загрузка объявлений...</h4>
            <p>Данные подгружаются с сервера</p>
        </div>
    `;
}

function fetchLovedItems() {
  const lovedList = document.getElementById("lovedList");
  lovedList.innerHTML = `
        <div class="no-loved">Вы ещё ничего не добавили в избранное</div>
    `;
}

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
        <div class="category-item" data-category="${category}" role="menuitem">
            <span>${getCategoryName(category)}</span>
            <span class="sub-arrow">›</span>
            <div class="subcategories-menu">
                ${categoriesData[category]
                  .map(
                    (sub) => `
                    <div class="subcategory-item" role="menuitem">${sub}</div>
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
    const isExpanded = categoriesMenu.classList.contains("show");
    categoriesBtn.setAttribute("aria-expanded", !isExpanded);
    categoriesMenu.classList.toggle("show");
  });

  document.addEventListener("click", () => {
    categoriesMenu.classList.remove("show");
    categoriesBtn.setAttribute("aria-expanded", "false");
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
  const searchModal = document.getElementById("searchModal");
  const closeSearch = document.getElementById("closeSearch");
  const searchInput = document.getElementById("searchInput");

  searchBtn?.addEventListener("click", () => {
    searchOverlay.style.display = "flex";
    document.body.style.overflow = "hidden";
    setTimeout(() => {
      searchInput.focus();
    }, 300);
  });

  function closeSearchModal() {
    searchOverlay.style.display = "none";
    document.body.style.overflow = "";
  }

  closeSearch?.addEventListener("click", closeSearchModal);
  searchOverlay?.addEventListener("click", (e) => {
    if (e.target === searchOverlay) closeSearchModal();
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

  function closeLoved() {
    lovedOverlay.style.display = "none";
    document.body.style.overflow = "";
  }

  closeLovedBtn?.addEventListener("click", closeLoved);
  lovedOverlay?.addEventListener("click", (e) => {
    if (e.target === lovedOverlay) closeLoved();
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
    if (confirm("Вы уверены, что хотите выйти?")) {
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
  const icon = themeToggleBtn.querySelector(".theme-icon");
  if (icon) icon.textContent = savedTheme === "dark" ? "☀️" : "🌙";

  themeToggleBtn.addEventListener("click", () => {
    const current = document.documentElement.getAttribute("data-theme");
    const next = current === "dark" ? "light" : "dark";
    document.documentElement.setAttribute("data-theme", next);
    localStorage.setItem("theme", next);
    if (icon) icon.textContent = next === "dark" ? "☀️" : "🌙";
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

  function closeMobileMenuFn() {
    mobileMenuOverlay.style.display = "none";
    document.body.style.overflow = "";
  }

  closeMobileMenu?.addEventListener("click", closeMobileMenuFn);
  mobileMenuOverlay?.addEventListener("click", (e) => {
    if (e.target === mobileMenuOverlay) closeMobileMenuFn();
  });
}

document.addEventListener("keydown", function (e) {
  if (e.key === "Escape") {
    const modals = [
      "searchOverlay",
      "mobileMenuOverlay",
      "profileMenuOverlay",
      "chatOverlay",
      "lovedOverlay",
    ];
    modals.forEach((id) => {
      const el = document.getElementById(id);
      if (el.style.display === "flex") {
        el.style.display = "none";
        document.body.style.overflow = "";
      }
    });
  }
});
if (document.querySelector(".ad-page")) {
  setTimeout(() => {
    if (typeof initAdPage === "function") {
      initAdPage();
    }
  }, 500);
}
