// ad.js — только функционал страницы объявления
document.addEventListener("DOMContentLoaded", function () {
  // Инициализируем страницу объявления, если находим её признаки
  if (
    document.getElementById("adTitle") &&
    document.getElementById("adAuthor")
  ) {
    initAdPage();
  }
});

// ----------------------------
// AD-SPECIFIC LOGIC
// ----------------------------

function initAdPage() {
  // Получаем adId из URL: /ads/abc-123/ → adId = "abc-123"
  const pathParts = window.location.pathname.split("/").filter((p) => p);
  const adId = pathParts.length >= 2 ? pathParts[pathParts.length - 1] : null;
  if (!adId) return;

  // === Галерея изображений ===
  // Используем setTimeout чтобы убедиться что DOM полностью загружен
  setTimeout(() => {
    initImageGallery();
  }, 100);

  // === Избранное ===
  const favBtn = document.getElementById("favoriteBtn");
  if (favBtn) {
    favBtn.addEventListener("click", () => {
      alert("Функция избранного пока не реализована");
    });
  }

  // === Написать продавцу ===
  const contactBtn = document.getElementById("contactBtn");
  if (contactBtn) {
    contactBtn.addEventListener("click", (e) => {
      e.preventDefault();
      localStorage.setItem("chat_ad_id", adId);
      window.location.href = "/messages";
    });
  }

  // === Показать телефон ===
  const showBtn = document.getElementById("showPhoneBtn");
  const phoneDisp = document.getElementById("phoneDisplay");
  if (showBtn && phoneDisp) {
    showBtn.addEventListener("click", () => {
      fetch(`/api/ads/${encodeURIComponent(adId)}/phone`)
        .then((r) => r.text())
        .then((phone) => {
          phoneDisp.textContent = phone;
          showBtn.style.display = "none";
        })
        .catch(() => alert("Не удалось загрузить номер телефона"));
    });
  }
}

function initImageGallery() {
  const mainImage = document.getElementById("mainImage");
  const thumbnails = document.querySelectorAll(".ad-thumbnail");
  
  if (!mainImage || thumbnails.length === 0) {
    return;
  }

  // Обработчик клика на миниатюры
  thumbnails.forEach((thumb) => {
    thumb.style.cursor = "pointer";
    
    thumb.addEventListener("click", (e) => {
      e.preventDefault();
      
      // Меняем главное изображение
      mainImage.src = thumb.src;
      mainImage.alt = thumb.alt;
      
      // Добавляем визуальный эффект активной миниатюры
      thumbnails.forEach(t => t.classList.remove("active"));
      thumb.classList.add("active");
    });
  });
}
