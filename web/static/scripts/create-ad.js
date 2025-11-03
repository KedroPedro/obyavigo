const REGIONS = {
  "Минская область": ["Минск", "Борисов", "Солигорск", "Молодечно", "Слуцк"],
  "Гомельская область": ["Гомель", "Мозырь", "Жлобин", "Речица", "Петриков"],
  "Могилёвская область": ["Могилёв", "Кричев", "Славгород", "Шклов"],
  "Витебская область": ["Витебск", "Новополоцк", "Полоцк", "Браслав", "Орша"],
  "Гродненская область": ["Гродно", "Лида", "Слоним", "Щучин", "Волковыск"],
  "Брестская область": ["Брест", "Пинск", "Барановичи", "Кобрин"],
};
const categories = [
  {
    id: "auto",
    name: "Авто",
    icon: "🚗",
    subcategories: [
      { id: "cars", name: "Легковые автомобили" },
      { id: "trucks", name: "Грузовики" },
      { id: "motorcycles", name: "Мотоциклы" },
      { id: "special", name: "Спецтехника" },
      { id: "parts", name: "Запчасти и аксессуары" },
    ],
  },
  {
    id: "realty",
    name: "Недвижимость",
    icon: "🏠",
    subcategories: [
      { id: "apartments", name: "Квартиры" },
      { id: "houses", name: "Дома и коттеджи" },
      { id: "commercial", name: "Коммерческая недвижимость" },
      { id: "land", name: "Земельные участки" },
      { id: "garages", name: "Гаражи и стоянки" },
    ],
  },
  {
    id: "electronics",
    name: "Электроника",
    icon: "📱",
    subcategories: [
      { id: "phones", name: "Смартфоны и планшеты" },
      { id: "computers", name: "Ноутбуки и компьютеры" },
      { id: "audio", name: "Телевизоры и аудио" },
      { id: "photo", name: "Фото и видео" },
      { id: "appliances", name: "Бытовая техника" },
    ],
  },
  {
    id: "fashion",
    name: "Мода",
    icon: "👗",
    subcategories: [
      { id: "clothing", name: "Одежда" },
      { id: "shoes", name: "Обувь" },
      { id: "accessories", name: "Аксессуары" },
      { id: "jewelry", name: "Часы и украшения" },
      { id: "cosmetics", name: "Косметика" },
    ],
  },
  {
    id: "home",
    name: "Для дома",
    icon: "🛌️",
    subcategories: [
      { id: "furniture", name: "Мебель" },
      { id: "interior", name: "Интерьер" },
      { id: "tableware", name: "Посуда" },
      { id: "textile", name: "Текстиль" },
      { id: "household", name: "Хозяйственные товары" },
    ],
  },
  {
    id: "hobby",
    name: "Хобби",
    icon: "⚽",
    subcategories: [
      { id: "sports", name: "Спорт и отдых" },
      { id: "books", name: "Книги и журналы" },
      { id: "collectibles", name: "Коллекционирование" },
      { id: "music", name: "Музыкальные инструменты" },
      { id: "tourism", name: "Туризм и рыбалка" },
    ],
  },
  {
    id: "services",
    name: "Услуги",
    icon: "🔧",
    subcategories: [
      { id: "repair", name: "Ремонт техники" },
      { id: "beauty", name: "Красота и здоровье" },
      { id: "education", name: "Образование" },
      { id: "transport", name: "Транспортные услуги" },
      { id: "construction", name: "Ремонт и строительство" },
    ],
  },
  {
    id: "work",
    name: "Работа",
    icon: "💼",
    subcategories: [
      { id: "vacancies", name: "Вакансии" },
      { id: "resume", name: "Резюме" },
      { id: "freelance", name: "Фриланс" },
      { id: "remote", name: "Удаленная работа" },
      { id: "parttime", name: "Подработка" },
    ],
  },
];
document.addEventListener("DOMContentLoaded", function () {
  initStepNavigation();
  initCategorySelection();
  initImageUpload();
  initRegionCitySelection();
  initDraftSaving();
  initPreview();
});
function initStepNavigation() {
  const nextButtons = document.querySelectorAll(".next-btn");
  const prevButtons = document.querySelectorAll(".prev-btn");
  const steps = document.querySelectorAll(".step");
  const formSteps = document.querySelectorAll(".form-step");
  nextButtons.forEach((button) => {
    button.addEventListener("click", () => {
      const currentStep = parseInt(button.closest(".form-step").dataset.step);
      if (validateStep(currentStep)) {
        showStep(currentStep + 1);
      }
    });
  });
  prevButtons.forEach((button) => {
    button.addEventListener("click", () => {
      const currentStep = parseInt(button.closest(".form-step").dataset.step);
      showStep(currentStep - 1);
    });
  });
  function showStep(stepNumber) {
    formSteps.forEach((step) => step.classList.remove("active"));
    steps.forEach((step) => step.classList.remove("active"));
    document
      .querySelector(`.form-step[data-step="${stepNumber}"]`)
      .classList.add("active");
    document
      .querySelector(`.step[data-step="${stepNumber}"]`)
      .classList.add("active");
    document.querySelector(".prev-btn").disabled = stepNumber === 1;
    document.querySelector(".next-btn").disabled = stepNumber === 3;
    if (stepNumber === 3) {
      if (typeof updatePreview === "function") updatePreview();
    }
  }
}
function initCategorySelection() {
  const grid = document.getElementById("categoryGrid");
  grid.innerHTML = categories
    .map(
      (cat) => `
        <div class="category-card" data-category="${cat.id}">
            <div class="category-icon">${cat.icon}</div>
            <h3>${cat.name}</h3>
        </div>
    `,
    )
    .join("");
  const cards = document.querySelectorAll(".category-card");
  const subcategoriesSection = document.getElementById("subcategoriesSection");
  const subcategoriesGrid = document.getElementById("subcategoriesGrid");
  cards.forEach((card) => {
    card.addEventListener("click", () => {
      cards.forEach((c) => c.classList.remove("selected"));
      card.classList.add("selected");
      const categoryId = card.dataset.category;
      const category = categories.find((cat) => cat.id === categoryId);
      if (category && category.subcategories) {
        subcategoriesGrid.innerHTML = category.subcategories
          .map(
            (sub) => `
                    <div class="subcategory-card" data-subcategory="${sub.id}">
                        <h4>${sub.name}</h4>
                    </div>
                `,
          )
          .join("");
        subcategoriesSection.style.display = "block";
        const subcategoryCards = document.querySelectorAll(".subcategory-card");
        subcategoryCards.forEach((subCard) => {
          subCard.addEventListener("click", () => {
            subcategoryCards.forEach((sc) => sc.classList.remove("selected"));
            subCard.classList.add("selected");
            document.getElementById("createAdForm").dataset.category =
              categoryId;
            document.getElementById("createAdForm").dataset.subcategory =
              subCard.dataset.subcategory;
            document.querySelector(".next-btn").disabled = false;
          });
        });
      } else {
        subcategoriesSection.style.display = "none";
        document.getElementById("createAdForm").dataset.category = categoryId;
        document.querySelector(".next-btn").disabled = false;
      }
    });
  });
}
function validateStep(stepNumber) {
  clearAllErrors();
  if (stepNumber === 1) {
    const category = document.getElementById("createAdForm").dataset.category;
    if (!category) {
      alert("Пожалуйста, выберите категорию");
      return false;
    }
    return true;
  }
  if (stepNumber === 2) {
    let isValid = true;
    const title = document.getElementById("adTitle").value.trim();
    if (!title) {
      showError("titleError", "Введите заголовок");
      isValid = false;
    } else if (title.length < 5) {
      showError("titleError", "Заголовок должен быть не менее 5 символов");
      isValid = false;
    }
    const price = document.getElementById("adPrice").value;
    if (!price || price <= 0) {
      showError("priceError", "Укажите корректную цену");
      isValid = false;
    }
    const region = document.getElementById("adRegion").value;
    if (!region) {
      showError("regionError", "Выберите регион");
      isValid = false;
    }
    const city = document.getElementById("adCity").value;
    if (!city) {
      showError("cityError", "Выберите город");
      isValid = false;
    }
    const desc = document.getElementById("adDescription").value.trim();
    if (!desc) {
      showError("descriptionError", "Введите описание");
      isValid = false;
    } else if (desc.length < 20) {
      showError(
        "descriptionError",
        "Описание должно быть не менее 20 символов",
      );
      isValid = false;
    }
    const condition = document.getElementById("adCondition").value;
    if (!condition) {
      showError("conditionError", "Выберите состояние товара");
      isValid = false;
    }
    return isValid;
  }
  if (stepNumber === 3) {
    let isValid = true;
    const rules = document.getElementById("agreeRules");
    if (!rules.checked) {
      showError("rulesError", "Необходимо согласие с правилами");
      isValid = false;
    }
    return isValid;
  }
  return true;
}
function initImageUpload() {
  const dropZone = document.getElementById("dropZone");
  const fileInput = document.getElementById("adImages");
  const preview = document.getElementById("imagePreview");
  dropZone.addEventListener("click", () => {
    fileInput.click();
  });
  ["dragenter", "dragover", "dragleave", "drop"].forEach((eventName) => {
    dropZone.addEventListener(eventName, preventDefaults, false);
  });
  function preventDefaults(e) {
    e.preventDefault();
    e.stopPropagation();
  }
  ["dragenter", "dragover"].forEach((eventName) => {
    dropZone.addEventListener(eventName, highlight, false);
  });
  ["dragleave", "drop"].forEach((eventName) => {
    dropZone.addEventListener(eventName, unhighlight, false);
  });
  function highlight() {
    dropZone.classList.add("dragover");
  }
  function unhighlight() {
    dropZone.classList.remove("dragover");
  }
  dropZone.addEventListener("drop", handleDrop, false);
  function handleDrop(e) {
    const dt = e.dataTransfer;
    const files = dt.files;
    handleFiles(files);
  }
  fileInput.addEventListener("change", () => {
    handleFiles(fileInput.files);
  });
  let uploadedImages = [];
  function handleFiles(files) {
    preview.innerHTML = "";
    uploadedImages = [];
    const validFiles = Array.from(files).slice(0, 10);
    validFiles.forEach((file) => {
      if (!file.type.match("image.*")) return;
      const reader = new FileReader();
      reader.onload = (e) => {
        uploadedImages.push(e.target.result);
        const img = document.createElement("img");
        img.src = e.target.result;
        preview.appendChild(img);
        updatePreview();
      };
      reader.readAsDataURL(file);
    });
  }
}
function initRegionCitySelection() {
  const regionSelect = document.getElementById("adRegion");
  const citySelect = document.getElementById("adCity");
  if (!regionSelect || !citySelect) {
    console.error("Region or city select not found!");
    return;
  }
  console.log("Region-city selection initialized");
  console.log("Available regions:", Object.keys(REGIONS));
  regionSelect.addEventListener("change", () => {
    const selectedRegion = regionSelect.value;
    console.log("Region changed to:", selectedRegion);
    citySelect.innerHTML = '<option value="">Выберите город</option>';
    if (selectedRegion && REGIONS[selectedRegion]) {
      console.log("Adding cities:", REGIONS[selectedRegion]);
      REGIONS[selectedRegion].forEach((city) => {
        const option = document.createElement("option");
        option.value = city;
        option.textContent = city;
        citySelect.appendChild(option);
      });
      console.log("Cities added, total options:", citySelect.options.length);
    }
    updatePreview();
  });
  citySelect.addEventListener("change", updatePreview);
}
function initPreview() {
  const fields = ["adTitle", "adPrice", "adCity", "adDescription"];
  fields.forEach((id) => {
    const el = document.getElementById(id);
    if (el) {
      el.addEventListener("input", updatePreview);
    }
  });
  updatePreview();
}
function updatePreview() {
  document.getElementById("previewTitle").textContent =
    document.getElementById("adTitle").value || "Заголовок объявления";
  const price = document.getElementById("adPrice").value;
  document.getElementById("previewPrice").textContent = price
    ? `${price} руб.`
    : "Цена не указана";
  const city = document.getElementById("adCity").value;
  document.getElementById("previewCity").textContent = city || "—";
  document.getElementById("previewDescription").textContent =
    document.getElementById("adDescription").value || "Описание...";
  const previewImagePlaceholder = document.querySelector(".preview-image-placeholder");
  if (previewImagePlaceholder) {
    const imagePreview = document.getElementById("imagePreview");
    const firstImage = imagePreview?.querySelector("img");
    if (firstImage && firstImage.src) {
      previewImagePlaceholder.style.backgroundImage = `url(${firstImage.src})`;
      previewImagePlaceholder.style.backgroundSize = "cover";
      previewImagePlaceholder.style.backgroundPosition = "center";
      previewImagePlaceholder.textContent = "";
    } else {
      previewImagePlaceholder.style.backgroundImage = "none";
      previewImagePlaceholder.textContent = "📷";
    }
  }
}
function initDraftSaving() {
  const fields = [
    "adTitle",
    "adPrice",
    "adRegion",
    "adCity",
    "adDescription",
  ];
  const draft = localStorage.getItem("obyavigo_draft");
  if (draft) {
    const data = JSON.parse(draft);
    if (data.adRegion) {
      const regionSelect = document.getElementById("adRegion");
      if (regionSelect) {
        regionSelect.value = data.adRegion;
        const citySelect = document.getElementById("adCity");
        if (citySelect && REGIONS[data.adRegion]) {
          citySelect.innerHTML = '<option value="">Выберите город</option>';
          REGIONS[data.adRegion].forEach((city) => {
            const option = document.createElement("option");
            option.value = city;
            option.textContent = city;
            citySelect.appendChild(option);
          });
        }
      }
    }
    fields.forEach((field) => {
      const el = document.getElementById(field);
      if (el && data[field] !== undefined) {
        el.value = data[field];
      }
    });
    updatePreview();
  }
  setInterval(() => {
    const draftData = {};
    fields.forEach((field) => {
      const el = document.getElementById(field);
      if (el) draftData[field] = el.value;
    });
    localStorage.setItem("obyavigo_draft", JSON.stringify(draftData));
  }, 5000);
}
function showError(elementId, message) {
  document.getElementById(elementId).textContent = message;
}
function clearAllErrors() {
  const errors = document.querySelectorAll(".error-message");
  errors.forEach((el) => (el.textContent = ""));
}
function isValidBelarusPhone(phone) {
  const cleaned = phone.replace(/\D/g, "");
  if (cleaned.length === 12 && cleaned.startsWith("375")) {
    return /^375(25|29|33|44)\d{7}$/.test(cleaned);
  }
  if (cleaned.length === 9) {
    return /^(25|29|33|44)\d{7}$/.test(cleaned);
  }
  return false;
}
document
  .getElementById("createAdForm")
  .addEventListener("submit", function (e) {
    e.preventDefault();
    if (!validateStep(3)) return;
    const formData = new FormData();
    const categorySlug = this.dataset.category;
    const categoryNameMap = {
      auto: "Авто",
      realty: "Недвижимость",
      electronics: "Электроника",
      fashion: "Мода",
      home: "Для дома",
      hobby: "Хобби",
      services: "Услуги",
      work: "Работа",
    };
    const categoryName = categoryNameMap[categorySlug] || "Другое";
    formData.append("categoryName", categoryName);
    const subcategorySlug = this.dataset.subcategory;
    if (subcategorySlug) {
      const category = categories.find((cat) => cat.id === categorySlug);
      const subcategory = category?.subcategories?.find(
        (sub) => sub.id === subcategorySlug,
      );
      if (subcategory) {
        formData.append("subcategoryName", subcategory.name);
      }
    }
    const cityValue = document.getElementById("adCity").value;
    const locationName = cityValue || "Не указан";
    formData.append("locationName", locationName);
    formData.append("title", document.getElementById("adTitle").value.trim());
    formData.append(
      "desc",
      document.getElementById("adDescription").value.trim(),
    );
    formData.append(
      "price",
      Number(document.getElementById("adPrice").value.trim()) * 100,
    );
    formData.append("condition", document.getElementById("adCondition").value);
    const files = document.getElementById("adImages").files;
    for (let i = 0; i < files.length; i++) {
      formData.append("images", files[i]);
    }
    fetch("/api/create-ad/", {
      method: "POST",
      body: formData,
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.success) {
          localStorage.removeItem("obyavigo_draft");
          alert("Объявление успешно опубликовано");
          window.location.href = data.url || "/ads/" + data.adId;
        } else {
          alert("Ошибка: " + (data.message || "Не удалось создать объявление"));
        }
      })
      .catch((error) => {
        console.error("Ошибка при создании объявления:", error);
        alert("Ошибка подключения к серверу");
      });
  });
