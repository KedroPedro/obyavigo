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
    this.initHeader();
    this.initCategories();
    this.initSearch();
    this.initLoved();
    this.initMobileMenu();
    this.initTheme();
    this.loadUserAvatar();
  }
  setupEventListeners() {
    document.addEventListener("keydown", this.handleGlobalKeydown.bind(this));
    window.addEventListener("scroll", this.throttle(this.handleScroll.bind(this), 16));
    window.addEventListener("resize", this.debounce(this.handleResize.bind(this), 250));
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
    const modalIds = ["searchOverlay", "mobileMenuOverlay", "chatOverlay", "lovedOverlay"];
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
  initHeader() {}
  navigateToAds(category, subcategory) {
    const params = new URLSearchParams();
    if (category) {
      const categoryName = this.getCategoryName(category);
      params.set("category", categoryName);
    }
    if (subcategory) params.set("subcategory", subcategory);
    window.location.href = `/ads/?${params.toString()}`;
  }
  initCategories() {
    const categoriesBtn = document.getElementById("categoriesBtn");
    const categoriesMenu = document.getElementById("categoriesMenu");
    if (!categoriesBtn || !categoriesMenu) return;
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
      services: ["Ремонт техники", "Красота и здоровье", "Образование", "Транспортные услуги", "Ремонт и строительство"],
      auto: ["Легковые автомобили", "Мотоциклы", "Грузовики", "Спецтехника", "Запчасти и аксессуары"],
      realty: ["Квартиры", "Дома и коттеджи", "Коммерческая недвижимость", "Земельные участки", "Гаражи и стоянки"],
      electronics: ["Смартфоны и планшеты", "Ноутбуки и компьютеры", "Телевизоры и аудио", "Фото и видео", "Бытовая техника"],
      work: ["Вакансии", "Резюме", "Фриланс", "Удаленная работа", "Подработка"],
      fashion: ["Одежда", "Обувь", "Аксессуары", "Часы и украшения", "Косметика"],
      home: ["Мебель", "Интерьер", "Посуда", "Текстиль", "Хозяйственные товары"],
      hobby: ["Спорт и отдых", "Книги и журналы", "Коллекционирование", "Музыкальные инструменты", "Туризм и рыбалка"],
    };
    const categoriesHtml = Object.keys(categoriesData)
      .map((category) => {
        const categoryName = categoryNameMap[category] || category;
        const subcategoriesHtml = categoriesData[category]
          .map((sub) => `<button class="subcategory-btn" data-subcategory="${sub}">${sub}</button>`)
          .join("");
        return `
          <div class="category-item-overlay" data-category="${category}">
            <div class="category-header"><strong>${categoryName}</strong></div>
            <div class="subcategories-list">${subcategoriesHtml}</div>
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
        <div class="categories-grid-overlay">${categoriesHtml}</div>
      </div>
    `;
    categoriesMenu.style.display = "none";
    categoriesBtn.setAttribute("aria-expanded", "false");
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
    document.getElementById("closeCategories")?.addEventListener("click", closeCategories);
    categoriesMenu.addEventListener("click", (e) => {
      if (e.target === categoriesMenu) closeCategories();
    });
    categoriesMenu.addEventListener("click", (e) => {
      const target = e.target;
      if (!(target instanceof HTMLElement)) return;
      if (target.classList.contains("subcategory-btn")) {
        e.preventDefault();
        const subcategory = target.dataset.subcategory;
        const parent = target.closest(".category-item-overlay");
        const category = parent?.getAttribute("data-category") || "";
        closeCategories();
        this.navigateToAds(category, subcategory);
        return;
      }
      const categoryItem = target.closest(".category-item-overlay");
      if (categoryItem && (target.classList.contains("category-header") || target.closest(".category-header"))) {
        const category = categoryItem.getAttribute("data-category");
        closeCategories();
        this.navigateToAds(category);
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
    const searchSuggestions = document.getElementById("searchSuggestions");
    if (!searchBtn || !searchOverlay) return;
    
    searchBtn.addEventListener("click", () => {
      this.openModal("searchOverlay");
      setTimeout(() => {
        if (searchInput) searchInput.focus();
      }, 300);
    });
    
    closeSearch?.addEventListener("click", () => {
      this.closeModal("searchOverlay");
      if (searchInput) searchInput.value = "";
      if (searchSuggestions) searchSuggestions.innerHTML = "";
    });
    
    searchOverlay.addEventListener("click", (e) => {
      if (e.target === searchOverlay) {
        this.closeModal("searchOverlay");
        if (searchInput) searchInput.value = "";
        if (searchSuggestions) searchSuggestions.innerHTML = "";
      }
    });
    
    
    if (searchInput) {
      let searchTimeout;
      searchInput.addEventListener("input", (e) => {
        const query = e.target.value.trim();
        
        clearTimeout(searchTimeout);
        
        if (query.length < 2) {
          if (searchSuggestions) searchSuggestions.innerHTML = "";
          return;
        }
        
        searchTimeout = setTimeout(() => {
          this.performSearch(query);
        }, 300);
      });
      
      
      searchInput.addEventListener("keydown", (e) => {
        if (e.key === "Enter") {
          const query = e.target.value.trim();
          if (query.length > 0) {
            this.navigateToSearch(query);
          }
        }
      });
    }
  }
  
  async performSearch(query) {
    const searchSuggestions = document.getElementById("searchSuggestions");
    if (!searchSuggestions) return;
    
    try {
      searchSuggestions.innerHTML = '<div class="search-loading">Поиск...</div>';
      
      const response = await fetch(`/api/ads/?q=${encodeURIComponent(query)}&limit=3`);
      
      if (!response.ok) {
        throw new Error("Ошибка поиска");
      }
      
      const result = await response.json();
      
      if (!result.success || !result.data || result.data.length === 0) {
        searchSuggestions.innerHTML = '<div class="search-no-results">Ничего не найдено</div>';
        return;
      }
      
      // Ограничиваем до 3 результатов
      const limitedResults = result.data.slice(0, 3);
      this.displaySearchResults(limitedResults);
    } catch (error) {
      console.error("Ошибка при поиске:", error);
      searchSuggestions.innerHTML = '<div class="search-error">Ошибка при поиске. Попробуйте позже.</div>';
    }
  }
  
  displaySearchResults(ads) {
    const searchSuggestions = document.getElementById("searchSuggestions");
    if (!searchSuggestions) return;
    
    const resultsHtml = ads.map(ad => {
      const imageUrl = ad.imageID ? `/api/images/${ad.imageID}/` : '/static/pictures/logo.png';
      const price = ad.price > 0 ? `${ad.price} BYN` : 'Договорная';
      
      return `
        <div class="search-result-item" data-ad-id="${ad.adID}" role="option">
          <img src="${imageUrl}" alt="${ad.title}" class="search-result-image" onerror="this.src='/static/pictures/logo.png'">
          <div class="search-result-info">
            <h4 class="search-result-title">${ad.title}</h4>
            <p class="search-result-price">${price}</p>
            <p class="search-result-location">${ad.locationName || ''}</p>
          </div>
        </div>
      `;
    }).join("");
    
    searchSuggestions.innerHTML = resultsHtml + 
      '<div class="search-show-all">Нажмите Enter для полного поиска</div>';
    
    
    const resultItems = searchSuggestions.querySelectorAll('.search-result-item');
    resultItems.forEach(item => {
      item.addEventListener('click', () => {
        const adId = item.getAttribute('data-ad-id');
        if (adId) {
          window.location.href = `/ads/${adId}/`;
        }
      });
    });
  }
  
  navigateToSearch(query) {
    window.location.href = `/ads/?q=${encodeURIComponent(query)}`;
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
    const mobileNav = mobileMenuOverlay.querySelector(".mobile-nav");
    if (mobileNav) {
      mobileNav.addEventListener("click", (e) => {
        const target = e.target;
        if (!(target instanceof HTMLElement)) return;
        if (target.matches("button[data-category]")) {
          const category = target.getAttribute("data-category");
          this.closeModal("mobileMenuOverlay");
          this.navigateToAds(category);
        }
      });
    }
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
  async loadUserAvatar() {
    try {
      const response = await fetch('/api/user/profile/', {
        credentials: 'include'
      });

      if (!response.ok) return;

      const data = await response.json();
      
      if (data.profile_picture_id) {
        
        const avatarImages = document.querySelectorAll('.profile-img, .profile-avatar img, #profileAvatar');
        avatarImages.forEach(img => {
          if (img) {
            img.src = `/api/avatars/${data.profile_picture_id}/`;
          }
        });
      }
    } catch (error) {
      console.error('Ошибка при загрузке аватара:', error);
    }
  }
}
const headerManager = new HeaderManager();
