const CATEGORY_SUBCATEGORIES = {
  "Услуги": [
    "Ремонт техники",
    "Красота и здоровье",
    "Образование",
    "Транспортные услуги",
    "Ремонт и строительство",
  ],
  "Авто": [
    "Легковые автомобили",
    "Мотоциклы",
    "Грузовики",
    "Спецтехника",
    "Запчасти и аксессуары",
  ],
  "Недвижимость": [
    "Квартиры",
    "Дома и коттеджи",
    "Коммерческая недвижимость",
    "Земельные участки",
    "Гаражи и стоянки",
  ],
  "Электроника": [
    "Смартфоны и планшеты",
    "Ноутбуки и компьютеры",
    "Телевизоры и аудио",
    "Фото и видео",
    "Бытовая техника",
  ],
  "Работа": ["Вакансии", "Резюме", "Фриланс", "Удаленная работа", "Подработка"],
  "Мода": ["Одежда", "Обувь", "Аксессуары", "Часы и украшения", "Косметика"],
  "Для дома": ["Мебель", "Интерьер", "Посуда", "Текстиль", "Хозяйственные товары"],
  "Хобби": [
    "Спорт и отдых",
    "Книги и журналы",
    "Коллекционирование",
    "Музыкальные инструменты",
    "Туризм и рыбалка",
  ],
};
const REGIONS = {
  "Минская область": ["Минск", "Борисов", "Солигорск", "Молодечно", "Слуцк"],
  "Гомельская область": ["Гомель", "Мозырь", "Жлобин", "Речица", "Петриков"],
  "Могилёвская область": ["Могилёв", "Кричев", "Славгород", "Шклов"],
  "Витебская область": ["Витебск", "Новополоцк", "Полоцк", "Браслав", "Орша"],
  "Гродненская область": ["Гродно", "Лида", "Слоним", "Щучин", "Волковыск"],
  "Брестская область": ["Брест", "Пинск", "Барановичи", "Кобрин"],
};
(function () {
  const start = () => {
    console.log("[ADS] Initializing filters...");
    initFilters();
    initLoadMore();
    initSearchQuery();
  };
  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", start);
  } else {
    setTimeout(start, 0);
  }
})();
function populateSubcategories(selectedCategory) {
  const subSelect = document.getElementById("subcategoryFilter");
  if (!subSelect) return;
  subSelect.innerHTML = '<option value="">Все подкатегории</option>';
  if (!selectedCategory || !CATEGORY_SUBCATEGORIES[selectedCategory]) {
    subSelect.disabled = true;
    subSelect.setAttribute("disabled", "");
    return;
  }
  CATEGORY_SUBCATEGORIES[selectedCategory].forEach((name) => {
    const opt = document.createElement("option");
    opt.value = name;
    opt.textContent = name;
    subSelect.appendChild(opt);
  });
  subSelect.disabled = false;
  subSelect.removeAttribute("disabled");
}
function populateCities(region) {
  const citySelect = document.getElementById("cityFilter");
  if (!citySelect) return;
  const current = citySelect.value;
  citySelect.innerHTML = '<option value="">Все города</option>';
  if (!region || !REGIONS[region]) return;
  REGIONS[region].forEach((city) => {
    const opt = document.createElement("option");
    opt.value = city;
    opt.textContent = city;
    citySelect.appendChild(opt);
  });
  if (REGIONS[region].includes(current)) citySelect.value = current;
}
function initFilters() {
  const params = new URLSearchParams(window.location.search);
  const categoryParam = params.get("category") || "";
  const subcategoryParam = params.get("subcategory") || "";
  const regionParam = params.get("region") || "";
  const cityParam = params.get("city") || "";
  const conditionParam = params.get("condition") || "";
  const sortParam = params.get("sort") || "newest";
  console.log("[ADS] URL params:", {
    categoryParam,
    subcategoryParam,
    regionParam,
    cityParam,
    conditionParam,
    sortParam,
  });
  const categorySelect = document.getElementById("categoryFilter");
  if (
    categorySelect &&
    Array.from(categorySelect.options).some((o) => o.value === categoryParam)
  ) {
    categorySelect.value = categoryParam;
    console.log("[ADS] Category set to:", categoryParam);
    if (categoryParam) {
      populateSubcategories(categoryParam);
      if (subcategoryParam) {
        setTimeout(() => {
          const subSelect = document.getElementById("subcategoryFilter");
          if (subSelect) subSelect.value = subcategoryParam;
        }, 50);
      }
    }
  }
  const regionSelect = document.getElementById("regionFilter");
  if (regionSelect && regionParam) {
    regionSelect.value = regionParam;
    console.log("[ADS] Region set to:", regionParam);
    populateCities(regionParam);
  }
  const citySelect = document.getElementById("cityFilter");
  if (citySelect && cityParam && regionParam) {
    setTimeout(() => {
      citySelect.value = cityParam;
    }, 50);
  }
  const minPriceInput = document.getElementById("minPrice");
  const maxPriceInput = document.getElementById("maxPrice");
  const newCheckbox = document.getElementById("newCondition");
  const usedCheckbox = document.getElementById("usedCondition");
  const sortSelect = document.getElementById("sortFilter");
  if (minPriceInput) minPriceInput.value = params.get("min_price") || "";
  if (maxPriceInput) maxPriceInput.value = params.get("max_price") || "";
  if (newCheckbox) newCheckbox.checked = conditionParam === "new";
  if (usedCheckbox) usedCheckbox.checked = conditionParam === "used";
  if (sortSelect) sortSelect.value = sortParam;
  const clearFiltersBtn = document.getElementById("clearFilters");
  const applyFiltersBtn = document.getElementById("applyFilters");
  const sidebar = document.querySelector(".filters-sidebar");
  console.log("[ADS] Buttons found:", {
    clearFiltersBtn: !!clearFiltersBtn,
    applyFiltersBtn: !!applyFiltersBtn,
  });
  if (clearFiltersBtn) {
    clearFiltersBtn.addEventListener("click", (e) => {
      console.log("[ADS] Clear clicked");
      e.preventDefault();
      e.stopPropagation();
      const subEl = document.getElementById("subcategoryFilter");
      document.getElementById("categoryFilter").value = "";
      if (subEl) {
        subEl.value = "";
        subEl.disabled = true;
        subEl.setAttribute("disabled", "");
      }
      document.getElementById("regionFilter").value = "";
      document.getElementById("cityFilter").innerHTML =
        '<option value="">Все города</option>';
      document.getElementById("minPrice").value = "";
      document.getElementById("maxPrice").value = "";
      document.getElementById("newCondition").checked = false;
      document.getElementById("usedCondition").checked = false;
      document.getElementById("sortFilter").value = "newest";
      syncUrlWithFilters();
      window.location.href = window.location.href;
    });
  }
  if (applyFiltersBtn) {
    applyFiltersBtn.addEventListener("click", (e) => {
      console.log("[ADS] Apply clicked");
      e.preventDefault();
      e.stopPropagation();
      syncUrlWithFilters();
      window.location.href = window.location.href;
    });
  }
  const filterInputs = document.querySelectorAll(
    ".filter-select, .price-input, .filter-checkbox input",
  );
  filterInputs.forEach((input) => {
    input.addEventListener("change", (evt) => {
      const el = evt.target;
      if (!el) return;
      if (el.id === "categoryFilter") {
        populateSubcategories(el.value);
        const sub = document.getElementById("subcategoryFilter");
        if (sub) {
          sub.value = "";
          sub.disabled = !el.value;
          if (el.value) sub.removeAttribute("disabled");
          else sub.setAttribute("disabled", "");
        }
      }
      if (el.id === "regionFilter") {
        populateCities(el.value);
        const cityEl = document.getElementById("cityFilter");
        if (cityEl) cityEl.value = "";
      }
      syncUrlWithFilters();
      window.location.href = window.location.href;
    });
  });
  if (sidebar) {
    sidebar.addEventListener("click", (e) => {
      const t = e.target;
      if (t && t.id === "applyFilters") {
        e.preventDefault();
        syncUrlWithFilters();
        window.location.href = window.location.href;
      }
      if (t && t.id === "clearFilters") {
        e.preventDefault();
      }
    });
  }
}
function syncUrlWithFilters() {
  const params = new URLSearchParams(window.location.search);
  const setOrDelete = (key, val) => {
    if (val === undefined || val === null || val === "") params.delete(key);
    else params.set(key, String(val));
  };
  
  
  const searchQuery = params.get("q");
  
  setOrDelete(
    "category",
    document.getElementById("categoryFilter")?.value || "",
  );
  setOrDelete(
    "subcategory",
    document.getElementById("subcategoryFilter")?.value || "",
  );
  const regionValue = document.getElementById("regionFilter")?.value || "";
  setOrDelete("region", regionValue);
  const cityValue = document.getElementById("cityFilter")?.value || "";
  setOrDelete("city", cityValue);
  setOrDelete("location", cityValue);
  let minPrice = document.getElementById("minPrice")?.value || "";
  let maxPrice = document.getElementById("maxPrice")?.value || "";
  if (minPrice && maxPrice && parseFloat(minPrice) > parseFloat(maxPrice)) {
    const temp = minPrice;
    minPrice = maxPrice;
    maxPrice = temp;
    document.getElementById("minPrice").value = minPrice;
    document.getElementById("maxPrice").value = maxPrice;
  }
  setOrDelete("min_price", minPrice);
  setOrDelete("max_price", maxPrice);
  setOrDelete("sort", document.getElementById("sortFilter")?.value || "newest");
  const newChecked = document.getElementById("newCondition")?.checked;
  const usedChecked = document.getElementById("usedCondition")?.checked;
  if (newChecked && !usedChecked) {
    setOrDelete("condition", "new");
  } else if (usedChecked && !newChecked) {
    setOrDelete("condition", "used");
  } else {
    params.delete("condition");
  }
  
  
  if (searchQuery) {
    params.set("q", searchQuery);
  }
  
  history.replaceState(null, "", `${location.pathname}?${params.toString()}`);
}
function initLoadMore() {
  const loadMoreBtn = document.getElementById("loadMoreBtn");
  if (!loadMoreBtn) return;
  let currentPage = 1;
  loadMoreBtn.addEventListener("click", async () => {
    console.log("Загрузка дополнительных объявлений...");
    loadMoreBtn.textContent = "Загрузка...";
    loadMoreBtn.disabled = true;
    try {
      const params = new URLSearchParams(window.location.search);
      currentPage++;
      params.set("page", currentPage);
      const response = await fetch(`/api/ads?${params.toString()}`);
      const data = await response.json();
      if (data.success && data.data && data.data.length > 0) {
        const adsGrid = document.getElementById("adsGrid");
        data.data.forEach((ad) => {
          const adCard = createAdCard(ad);
          adsGrid.appendChild(adCard);
        });
        if (currentPage >= data.totalPages) {
          loadMoreBtn.style.display = "none";
        }
      } else {
        loadMoreBtn.style.display = "none";
      }
    } catch (error) {
      console.error("Ошибка загрузки объявлений:", error);
    } finally {
      loadMoreBtn.textContent = "Показать ещё";
      loadMoreBtn.disabled = false;
    }
  });
}
function createAdCard(ad) {
  const card = document.createElement("a");
  card.href = `/ads/${ad.adID}/`;
  card.className = "ad-card";
  card.style.cssText = "text-decoration: none; color: inherit";
  const createdDate = new Date(ad.createdAt);
  const formattedDate = createdDate.toLocaleDateString("ru-RU");
  const imageUrl = ad.imageID
    ? `/api/images/${ad.imageID}/`
    : "/static/pictures/logo.png";
  card.innerHTML = `
    <div class="ad-image">
      <img src="${imageUrl}" alt="${ad.title}" loading="lazy" />
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
        <span class="ad-date">${formattedDate}</span>
        <span class="ad-views">👁️ ${ad.viewsCount}</span>
      </div>
    </div>
  `;
  return card;
}

function initSearchQuery() {
  const params = new URLSearchParams(window.location.search);
  const searchQuery = params.get("q");
  const searchQueryDisplay = document.getElementById("searchQueryDisplay");
  const searchQueryText = document.getElementById("searchQueryText");
  const clearSearchBtn = document.getElementById("clearSearchBtn");
  
  if (searchQuery && searchQueryDisplay && searchQueryText) {
    searchQueryText.textContent = searchQuery;
    searchQueryDisplay.style.display = "block";
    
    if (clearSearchBtn) {
      clearSearchBtn.addEventListener("click", () => {
        const newParams = new URLSearchParams(window.location.search);
        newParams.delete("q");
        window.location.href = `${window.location.pathname}?${newParams.toString()}`;
      });
    }
  }
}
