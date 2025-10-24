document.addEventListener("DOMContentLoaded", function () {
  initStepNavigation();
  initCategorySelection();
  initFormValidation();
  initImageUpload();
  initAutoDetectCity();
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
  const categories = [
    {
      id: "auto",
      name: "Авто",
      icon: "🚗",
      subcategories: [
        { id: "cars", name: "Легковые автомобили" },
        { id: "trucks", name: "Грузовики" },
        { id: "motorcycles", name: "Мотоциклы" },
        { id: "parts", name: "Запчасти" },
        { id: "accessories", name: "Аксессуары" },
      ],
    },
    {
      id: "realty",
      name: "Недвижимость",
      icon: "🏠",
      subcategories: [
        { id: "apartments", name: "Квартиры" },
        { id: "houses", name: "Дома" },
        { id: "commercial", name: "Коммерческая" },
        { id: "land", name: "Земельные участки" },
        { id: "rent", name: "Аренда" },
      ],
    },
    {
      id: "electronics",
      name: "Электроника",
      icon: "📱",
      subcategories: [
        { id: "phones", name: "Телефоны" },
        { id: "computers", name: "Компьютеры" },
        { id: "audio", name: "Аудио" },
        { id: "photo", name: "Фото/Видео" },
        { id: "gaming", name: "Игры" },
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
        { id: "jewelry", name: "Украшения" },
        { id: "bags", name: "Сумки" },
      ],
    },
    {
      id: "home",
      name: "Для дома",
      icon: "🛋️",
      subcategories: [
        { id: "furniture", name: "Мебель" },
        { id: "appliances", name: "Бытовая техника" },
        { id: "decor", name: "Декор" },
        { id: "kitchen", name: "Кухня" },
        { id: "garden", name: "Сад/Огород" },
      ],
    },
    {
      id: "hobby",
      name: "Хобби",
      icon: "⚽",
      subcategories: [
        { id: "sports", name: "Спорт" },
        { id: "music", name: "Музыка" },
        { id: "books", name: "Книги" },
        { id: "collectibles", name: "Коллекционирование" },
        { id: "art", name: "Искусство" },
      ],
    },
    {
      id: "services",
      name: "Услуги",
      icon: "🔧",
      subcategories: [
        { id: "repair", name: "Ремонт" },
        { id: "cleaning", name: "Уборка" },
        { id: "transport", name: "Транспорт" },
        { id: "education", name: "Образование" },
        { id: "beauty", name: "Красота" },
      ],
    },
    {
      id: "work",
      name: "Работа",
      icon: "💼",
      subcategories: [
        { id: "fulltime", name: "Полная занятость" },
        { id: "parttime", name: "Частичная занятость" },
        { id: "freelance", name: "Фриланс" },
        { id: "internship", name: "Стажировка" },
        { id: "remote", name: "Удаленная работа" },
      ],
    },
  ];

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

    return isValid;
  }

  if (stepNumber === 3) {
    let isValid = true;

    const name = document.getElementById("adName").value.trim();
    if (!name) {
      showError("nameError", "Введите имя");
      isValid = false;
    }

    const phone = document.getElementById("adPhone").value.trim();
    if (!phone) {
      showError("phoneError", "Введите телефон");
      isValid = false;
    } else if (!isValidBelarusPhone(phone)) {
      showError("phoneError", "Неверный формат белорусского номера");
      isValid = false;
    }

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

  function handleFiles(files) {
    preview.innerHTML = "";
    const validFiles = Array.from(files).slice(0, 10);

    validFiles.forEach((file) => {
      if (!file.type.match("image.*")) return;

      const reader = new FileReader();
      reader.onload = (e) => {
        const img = document.createElement("img");
        img.src = e.target.result;
        preview.appendChild(img);
        updatePreview();
      };
      reader.readAsDataURL(file);
    });
  }
}
function initAutoDetectCity() {
  document.getElementById("autoDetectCity").addEventListener("click", () => {
    // Для демо — просто ставим Минск
    document.getElementById("adCity").value = "minsk";
    updatePreview();
  });
}
function initPreview() {
  const fields = ["adTitle", "adPrice", "adCity", "adDescription", "adName"];
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
  const cityNames = {
    minsk: "Минск",
    brest: "Брест",
    vitebsk: "Витебск",
    gomel: "Гомель",
    grodno: "Гродно",
    mogilev: "Могилёв",
    other: "Другой город",
  };
  document.getElementById("previewCity").textContent = city
    ? cityNames[city] || city
    : "—";

  document.getElementById("previewDescription").textContent =
    document.getElementById("adDescription").value || "Описание...";
}
function initDraftSaving() {
  const fields = [
    "adTitle",
    "adPrice",
    "adCity",
    "adDescription",
    "adPhone",
    "adName",
  ];
  const draft = localStorage.getItem("obyavigo_draft");
  if (draft) {
    const data = JSON.parse(draft);
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

    // === Категория: slug → имя ===
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

    // === Город: значение → название ===
    const cityValue = document.getElementById("adCity").value;
    const cityNames = {
      minsk: "Минск",
      brest: "Брест",
      vitebsk: "Витебск",
      gomel: "Гомель",
      grodno: "Гродно",
      mogilev: "Могилёв",
      other: "Другой город",
    };
    const locationName = cityNames[cityValue] || cityValue || "Не указан";
    formData.append("locationName", locationName);

    // === Остальные текстовые поля ===
    formData.append("title", document.getElementById("adTitle").value.trim());
    formData.append(
      "desc",
      document.getElementById("adDescription").value.trim(),
    );
    formData.append(
      "price",
      Number(document.getElementById("adPrice").value.trim()) * 100,
    );
    formData.append("phone", document.getElementById("adPhone").value.trim());
    // Если нужно condition (с опечаткой "codition" — исправьте на сервере!)
    formData.append("codition", "used"); // или "new"

    // === Изображения ===
    const files = document.getElementById("adImages").files;
    for (let i = 0; i < files.length; i++) {
      formData.append("images", files[i]); // имя поля "images" — как ожидает бэкенд
    }

    // === Отправка как multipart/form-data (автоматически) ===
    fetch("/api/create-ad/", {
      method: "POST",
      body: formData,
      // ⚠️ НЕ указываем Content-Type! Браузер сам установит boundary
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.success) {
          localStorage.removeItem("obyavigo_draft");
          alert("Объявление успешно опубликовано!");
          window.location.href = "/my-ads";
        } else {
          alert("Ошибка: " + (data.message || "Не удалось создать объявление"));
        }
      })
      .catch((error) => {
        console.error("Ошибка при создании объявления:", error);
        alert("Ошибка подключения к серверу");
      });
  });
