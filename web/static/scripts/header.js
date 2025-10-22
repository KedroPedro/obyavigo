class HeaderManager {
  constructor() {
    this.isInitialized = false;
    this.modals = new Map();
    this.init();
  }

  init() {
    if (this.isInitialized) return;

    if (document.readyState === "loading") {
      document.addEventListener("DOMContentLoaded", () => {
        this.initializeComponents();
        this.setupEventListeners();
        this.isInitialized = true;
      });
    } else {
      this.initializeComponents();
      this.setupEventListeners();
      this.isInitialized = true;
    }
  }

  initializeComponents() {
    console.log("HeaderManager: initializing components...");
    this.initHeader();
    this.initCategories();
    this.initSearch();
    this.initLoved();
    this.initMobileMenu();
    this.initTheme();
    console.log("HeaderManager: all components initialized");
  }

  setupEventListeners() {
    document.addEventListener("keydown", this.handleGlobalKeydown.bind(this));
    window.addEventListener(
      "scroll",
      this.throttle(this.handleScroll.bind(this), 16),
    );
    window.addEventListener(
      "resize",
      this.debounce(this.handleResize.bind(this), 250),
    );
  }

  throttle(func, limit) {
    let inThrottle;
    return function () {
      const args = arguments;
      const context = this;
      if (!inThrottle) {
        func.apply(context, args);
        inThrottle = true;
        setTimeout(() => (inThrottle = false), limit);
      }
    };
  }

  debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
      const later = () => {
        clearTimeout(timeout);
        func(...args);
      };
      clearTimeout(timeout);
      timeout = setTimeout(later, wait);
    };
  }

  handleGlobalKeydown(e) {
    if (e.key === "Escape") {
      this.closeAllModals();
    }
  }

  handleScroll() {
    const header = document.getElementById("header");
    if (!header) return;

    const scrollY = window.scrollY;
    const lastScrollY = this.lastScrollY || 0;

    if (scrollY > 100) {
      header.classList.add("scrolled");
      if (scrollY > lastScrollY && scrollY > 200) {
        header.classList.add("hidden");
      } else {
        header.classList.remove("hidden");
      }
    } else {
      header.classList.remove("scrolled", "hidden");
    }

    this.lastScrollY = scrollY;
  }

  handleResize() {
    this.updateMobileMenuVisibility();
  }

  closeAllModals() {
    const modalIds = [
      "searchOverlay",
      "mobileMenuOverlay",
      "chatOverlay",
      "lovedOverlay",
    ];
    modalIds.forEach((id) => {
      const modal = document.getElementById(id);
      if (modal && modal.style.display === "flex") {
        this.closeModal(id);
      }
    });
  }

  closeModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
      modal.style.display = "none";
      document.body.style.overflow = "";
    }
  }

  openModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
      modal.style.display = "flex";
      document.body.style.overflow = "hidden";
    }
  }

  initHeader() {
    // Header initialization
  }

  initCategories() {
    const categoriesBtn = document.getElementById("categoriesBtn");
    const categoriesMenu = document.getElementById("categoriesMenu");

    if (!categoriesBtn || !categoriesMenu) return;

    // Локальный маппинг имён — без зависимости от this
    const categoryNameMap = {
      services: "Услуги",
      auto: "Авто",
      realty: "Недвижимость",
      electronics: "Электроника",
      work: "Работа",
      fashion: "Мода",
      home: "Для дома",
      hobby: "Хобби",
    };

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
      fashion: [
        "Одежда",
        "Обувь",
        "Аксессуары",
        "Часы и украшения",
        "Косметика",
      ],
      home: [
        "Мебель",
        "Интерьер",
        "Посуда",
        "Текстиль",
        "Хозяйственные товары",
      ],
      hobby: [
        "Спорт и отдых",
        "Книги и журналы",
        "Коллекционирование",
        "Музыкальные инструменты",
        "Туризм и рыбалка",
      ],
    };

    // Генерация HTML без использования this внутри шаблона
    const categoriesHtml = Object.keys(categoriesData)
      .map((category) => {
        const categoryName = categoryNameMap[category] || category;
        const subcategoriesHtml = categoriesData[category]
          .map(
            (sub) =>
              `<button class="subcategory-btn" data-subcategory="${sub}">${sub}</button>`,
          )
          .join("");

        return `
          <div class="category-item-overlay" data-category="${category}">
            <div class="category-header">
              <strong>${categoryName}</strong>
            </div>
            <div class="subcategories-list">
              ${subcategoriesHtml}
            </div>
          </div>
        `;
      })
      .join("");

    categoriesMenu.innerHTML = `
      <div class="categories-modal" role="dialog" aria-modal="true" aria-labelledby="categories-title">
        <div class="categories-modal-header">
          <h3 id="categories-title">Все категории</h3>
          <button class="close-search" id="closeCategories" aria-label="Закрыть">×</button>
        </div>
        <div class="categories-grid-overlay">
          ${categoriesHtml}
        </div>
      </div>
    `;

    // Явно скрываем меню при инициализации
    categoriesMenu.style.display = "none";
    categoriesBtn.setAttribute("aria-expanded", "false");

    // Обработчик открытия/закрытия
    categoriesBtn.addEventListener("click", (e) => {
      e.stopPropagation();
      const isExpanded = categoriesMenu.style.display === "flex";
      categoriesBtn.setAttribute("aria-expanded", (!isExpanded).toString());

      if (isExpanded) {
        categoriesMenu.style.display = "none";
        document.body.classList.remove("categories-open");
        document.body.style.overflow = "";
      } else {
        categoriesMenu.style.display = "flex";
        document.body.classList.add("categories-open");
        document.body.style.overflow = "hidden";
      }
    });

    const closeCategories = () => {
      categoriesMenu.style.display = "none";
      categoriesBtn.setAttribute("aria-expanded", "false");
      document.body.classList.remove("categories-open");
      document.body.style.overflow = "";
    };

    document
      .getElementById("closeCategories")
      ?.addEventListener("click", closeCategories);

    categoriesMenu.addEventListener("click", (e) => {
      if (e.target === categoriesMenu) {
        closeCategories();
      }
    });

    categoriesMenu.addEventListener("click", (e) => {
      if (e.target.classList.contains("subcategory-btn")) {
        e.preventDefault();
        const subcategory = e.target.dataset.subcategory;
        console.log("Выбрана подкатегория:", subcategory);
        closeCategories();
      }
    });
  }

  getCategoryName(key) {
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

  initSearch() {
    const searchBtn = document.getElementById("searchBtn");
    const searchOverlay = document.getElementById("searchOverlay");
    const closeSearch = document.getElementById("closeSearch");
    const searchInput = document.getElementById("searchInput");

    if (!searchBtn || !searchOverlay) return;

    searchBtn.addEventListener("click", () => {
      this.openModal("searchOverlay");
      setTimeout(() => {
        if (searchInput) searchInput.focus();
      }, 300);
    });

    closeSearch?.addEventListener("click", () => {
      this.closeModal("searchOverlay");
    });

    searchOverlay.addEventListener("click", (e) => {
      if (e.target === searchOverlay) {
        this.closeModal("searchOverlay");
      }
    });
  }

  initLoved() {
    const lovedBtn = document.getElementById("lovedBtn");
    const lovedOverlay = document.getElementById("lovedOverlay");
    const closeLovedBtn = document.getElementById("closeLovedBtn");

    if (!lovedBtn || !lovedOverlay) return;

    lovedBtn.addEventListener("click", () => {
      this.openModal("lovedOverlay");
    });

    closeLovedBtn?.addEventListener("click", () => {
      this.closeModal("lovedOverlay");
    });

    lovedOverlay.addEventListener("click", (e) => {
      if (e.target === lovedOverlay) {
        this.closeModal("lovedOverlay");
      }
    });
  }

  initMobileMenu() {
    const mobileMenuBtn = document.getElementById("mobileMenuBtn");
    const mobileMenuOverlay = document.getElementById("mobileMenuOverlay");
    const closeMobileMenu = document.getElementById("closeMobileMenu");

    if (!mobileMenuBtn || !mobileMenuOverlay) return;

    mobileMenuBtn.addEventListener("click", () => {
      this.openModal("mobileMenuOverlay");
    });

    closeMobileMenu?.addEventListener("click", () => {
      this.closeModal("mobileMenuOverlay");
    });

    mobileMenuOverlay.addEventListener("click", (e) => {
      if (e.target === mobileMenuOverlay) {
        this.closeModal("mobileMenuOverlay");
      }
    });
  }

  initTheme() {
    const themeToggleBtn = document.getElementById("themeToggleBtn");

    if (!themeToggleBtn) return;

    const savedTheme = localStorage.getItem("theme") || "light";
    document.documentElement.setAttribute("data-theme", savedTheme);

    const icon = themeToggleBtn.querySelector(".theme-icon");
    if (icon) {
      icon.textContent = savedTheme === "dark" ? "☀️" : "🌙";
    }

    themeToggleBtn.addEventListener("click", () => {
      const current = document.documentElement.getAttribute("data-theme");
      const newTheme = current === "dark" ? "light" : "dark";

      document.documentElement.setAttribute("data-theme", newTheme);
      localStorage.setItem("theme", newTheme);

      if (icon) {
        icon.textContent = newTheme === "dark" ? "☀️" : "🌙";
      }
    });
  }

  updateMobileMenuVisibility() {
    const mobileMenuBtn = document.getElementById("mobileMenuBtn");
    if (mobileMenuBtn) {
      mobileMenuBtn.style.display = window.innerWidth <= 768 ? "block" : "none";
    }
  }
}

// Initialize header manager
const headerManager = new HeaderManager();
console.log("HeaderManager initialized:", headerManager.isInitialized);

document.addEventListener("DOMContentLoaded", () => {
  console.log("HeaderManager: DOMContentLoaded fired");
  console.log("SearchBtn exists:", !!document.getElementById("searchBtn"));
  console.log(
    "CategoriesBtn exists:",
    !!document.getElementById("categoriesBtn"),
  );
});
