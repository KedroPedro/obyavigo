
document.addEventListener("DOMContentLoaded", function () {

  if (
    document.getElementById("adTitle") &&
    document.getElementById("adAuthor")
  ) {
    initAdPage();
  }
});





function initAdPage() {
  console.log("initAdPage запущена");
  const pathParts = window.location.pathname.split("/").filter((p) => p);
  const adId = pathParts.length >= 2 ? pathParts[pathParts.length - 1] : null;


  if (!adId) {
    return; 
  }


  setTimeout(() => {
    initImageGallery();
  }, 100);


  const favBtn = document.getElementById("favoriteBtn");
  if (favBtn) {
    checkFavoriteStatus(adId, favBtn);
    
    favBtn.addEventListener("click", () => {
      toggleFavorite(adId, favBtn);
    });
  }


  const contactBtn = document.getElementById("contactBtn");
  if (contactBtn) {
    contactBtn.addEventListener("click", (e) => {
      e.preventDefault();
      localStorage.setItem("chat_ad_id", adId);
      window.location.href = "/messages";
    });
  }

  initReportButton(adId);


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


  thumbnails.forEach((thumb) => {
    thumb.style.cursor = "pointer";
    
    thumb.addEventListener("click", (e) => {
      e.preventDefault();
      

      mainImage.src = thumb.src;
      mainImage.alt = thumb.alt;
      

      thumbnails.forEach(t => t.classList.remove("active"));
      thumb.classList.add("active");
    });
  });
}

function checkFavoriteStatus(adId, favBtn) {
  fetch(`/api/favorites/check/${encodeURIComponent(adId)}/`)
    .then((r) => {
      if (!r.ok) throw new Error("Failed to check favorite status");
      return r.json();
    })
    .then((data) => {
      if (data.isFavorite) {
        updateFavoriteButton(favBtn, true);
      } else {
        updateFavoriteButton(favBtn, false);
      }
    })
    .catch((err) => {
      console.error("Ошибка при проверке статуса избранного:", err);
    });
}

function toggleFavorite(adId, favBtn) {
  const isFavorite = favBtn.dataset.isFavorite === "true";
  
  const method = isFavorite ? "DELETE" : "POST";
  const url = `/api/favorites/${encodeURIComponent(adId)}/`;
  
  fetch(url, { method })
    .then((r) => {
      if (!r.ok) throw new Error("Failed to toggle favorite");
      return r.json();
    })
    .then((data) => {
      if (data.success) {
        updateFavoriteButton(favBtn, !isFavorite);
      }
    })
    .catch((err) => {
      console.error("Ошибка при изменении статуса избранного:", err);
      alert("Не удалось изменить статус избранного");
    });
}

function updateFavoriteButton(favBtn, isFavorite) {
  favBtn.dataset.isFavorite = isFavorite;
  const span = favBtn.querySelector("span:last-child");
  
  if (isFavorite) {
    favBtn.classList.add("favorite-active");
    if (span) span.textContent = "Удалить из избранного";
  } else {
    favBtn.classList.remove("favorite-active");
    if (span) span.textContent = "В избранное";
  }
}

function initReportButton(adId) {
  console.log("initReportButton запущена с adId:", adId);
  const reportBtn = document.getElementById("reportBtn");
  const reportModal = document.getElementById("reportModal");
  const closeReportModal = document.getElementById("closeReportModal");
  const cancelReportBtn = document.getElementById("cancelReportBtn");
  const reportForm = document.getElementById("reportForm");
  const reportDescription = document.getElementById("reportDescription");
  const charCount = document.querySelector(".char-count");

  if (!reportBtn || !reportModal) return;

  reportBtn.addEventListener("click", (e) => {
    e.preventDefault();
    reportModal.style.display = "flex";
    document.body.style.overflow = "hidden";
  });

  const closeModal = () => {
    reportModal.style.display = "none";
    document.body.style.overflow = "";
    reportForm.reset();
    if (charCount) charCount.textContent = "0 / 500";
  };

  closeReportModal.addEventListener("click", closeModal);
  cancelReportBtn.addEventListener("click", closeModal);
  
  reportModal.addEventListener("click", (e) => {
    if (e.target === reportModal) {
      closeModal();
    }
  });

  if (reportDescription && charCount) {
    reportDescription.addEventListener("input", () => {
      charCount.textContent = `${reportDescription.value.length} / 500`;
    });
  }

  reportForm.addEventListener("submit", (e) => {
    e.preventDefault();
    
    const reportType = document.getElementById("reportType").value;
    const description = reportDescription.value;

    if (!reportType) {
      alert("Пожалуйста, выберите причину жалобы");
      return;
    }

    fetch("/api/reports/", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        ad_id: adId,
        report_type: reportType,
        description: description
      })
    })
    .then((r) => {
      if (!r.ok) {
        if (r.status === 401) {
          throw new Error("Необходимо войти в систему");
        }
        throw new Error("Не удалось отправить жалобу");
      }
      return r.json();
    })
    .then(() => {
      alert("Жалоба успешно отправлена. Спасибо за обращение!");
      closeModal();
    })
    .catch((err) => {
      console.error("Ошибка при отправке жалобы:", err);
      alert(err.message || "Произошла ошибка при отправке жалобы");
    });
  });
}
