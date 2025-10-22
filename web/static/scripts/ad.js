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

  // === Избранное ===
  const favBtn = document.getElementById("favoriteBtn");
  if (favBtn) {
    // Проверка при загрузке
    fetch(`/api/loved/check/${encodeURIComponent(adId)}`)
      .then((r) => r.json())
      .then((data) => {
        if (data.is_loved) {
          favBtn.classList.add("active");
          favBtn.innerHTML = "<span>❤️</span><span>В избранном</span>";
        }
      })
      .catch(() => {
        // Опционально: обработка ошибки
      });

    // Переключение избранного
    favBtn.addEventListener("click", () => {
      const isLoved = favBtn.classList.contains("active");
      fetch("/api/loved/toggle", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ad_id: adId }),
      })
        .then((r) => r.json())
        .then((data) => {
          if (data.success) {
            favBtn.classList.toggle("active", !isLoved);
            favBtn.innerHTML = isLoved
              ? "<span>❤️</span><span>В избранное</span>"
              : "<span>❤️</span><span>В избранном</span>";
          }
        })
        .catch(() => alert("Ошибка при работе с избранным"));
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

  // === Похожие объявления ===
  const grid = document.getElementById("similarAdsGrid");
  if (grid) {
    fetch(`/api/ads/${encodeURIComponent(adId)}/similar`)
      .then((r) => r.json())
      .then((ads) => {
        if (!Array.isArray(ads) || ads.length === 0) {
          grid.innerHTML =
            '<p style="text-align:center; grid-column:1/-1;">Похожих объявлений не найдено</p>';
          return;
        }
        grid.innerHTML = ads
          .map(
            (ad) => `
              <div class="similar-ad-card">
                <img
                  src="${ad.image || "/static/pictures/placeholder.jpg"}"
                  alt="${ad.title}"
                  class="similar-ad-image"
                  loading="lazy"
                >
                <div class="similar-ad-content">
                  <h3>${ad.title}</h3>
                  <div class="similar-ad-price">${ad.price ? ad.price + " BYN" : "Договорная"}</div>
                  <div class="similar-ad-location">${ad.city || "—"}</div>
                </div>
              </div>
            `,
          )
          .join("");
      })
      .catch(() => {
        grid.innerHTML =
          '<p style="text-align:center; grid-column:1/-1;">Ошибка загрузки похожих объявлений</p>';
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
