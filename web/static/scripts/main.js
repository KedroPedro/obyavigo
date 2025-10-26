class ObyavigoApp {
  constructor() {
    this.isInitialized = false;
    this.init();
  }

  init() {
    if (this.isInitialized) return;

    document.addEventListener("DOMContentLoaded", () => {
      this.initializeComponents();
      this.isInitialized = true;
    });
  }

  initializeComponents() {
    this.loadDynamicContent();
  }
  loadDynamicContent() {
    this.showSkeletons();
    setTimeout(() => {
      this.fetchStats();
      this.fetchCategories();
      this.fetchFeaturedAds();
      this.fetchLovedItems();
    }, 1000);
  }

  showSkeletons() {
    if (window.location.pathname.startsWith('/ads')) return;
    
    const categoriesGrid = document.getElementById("categoriesGrid");
    const adsGrid = document.getElementById("adsGrid");

    if (categoriesGrid) {
      categoriesGrid.innerHTML = `
                <div class="skeleton-card"></div>
                <div class="skeleton-card"></div>
                <div class="skeleton-card"></div>
                <div class="skeleton-card"></div>
            `;
    }

    if (adsGrid) {
      adsGrid.innerHTML = `
                <div class="skeleton-ad"></div>
                <div class="skeleton-ad"></div>
                <div class="skeleton-ad"></div>
            `;
    }
  }

  fetchStats() {
    const stats = {
      yearsOnMarket: "15 лет",
      citiesCovered: "85 городов",
      appRating: "4.9/5",
      dailyDeals: "3 200+",
    };

    Object.entries(stats).forEach(([key, value]) => {
      const element = document.getElementById(
        key.replace(/([A-Z])/g, "-$1").toLowerCase(),
      );
      if (element) element.textContent = value;
    });
  }

  fetchCategories() {
    const categoriesGrid = document.getElementById("categoriesGrid");
    if (!categoriesGrid) return;
    categoriesGrid.innerHTML = `
            <div class="skeleton-card"></div>
            <div class="skeleton-card"></div>
            <div class="skeleton-card"></div>
            <div class="skeleton-card"></div>
        `;
  }

  fetchFeaturedAds() {
    const adsGrid = document.getElementById("adsGrid");
    if (!adsGrid) return;
    if (window.location.pathname.startsWith('/ads')) return;

    adsGrid.innerHTML = `
            <div class="skeleton-ad"></div>
            <div class="skeleton-ad"></div>
            <div class="skeleton-ad"></div>
        `;
  }

  fetchLovedItems() {
    const lovedList = document.getElementById("lovedList");
    if (!lovedList) return;

    lovedList.innerHTML = `
            <div class="no-loved">Вы ещё ничего не добавили в избранное</div>
        `;
  }

}
const app = new ObyavigoApp();
