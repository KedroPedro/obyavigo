document.addEventListener('DOMContentLoaded', function() {
    initFilters();
    loadAds();
    initLoadMore();
});
function initFilters() {
    const clearFiltersBtn = document.getElementById('clearFilters');
    const applyFiltersBtn = document.getElementById('applyFilters');
    clearFiltersBtn?.addEventListener('click', () => {
        document.getElementById('categoryFilter').value = '';
        document.getElementById('cityFilter').value = '';
        document.getElementById('minPrice').value = '';
        document.getElementById('maxPrice').value = '';
        document.getElementById('newCondition').checked = false;
        document.getElementById('usedCondition').checked = false;
        document.getElementById('sortFilter').value = 'newest';
        loadAds();
    });
    applyFiltersBtn?.addEventListener('click', () => {
        loadAds();
    });
    const filterInputs = document.querySelectorAll('.filter-select, .price-input, .filter-checkbox input');
    filterInputs.forEach(input => {
        input.addEventListener('change', () => {
            loadAds();
        });
    });
}
function loadAds() {
    const adsGrid = document.getElementById('adsGrid');
    const adsCount = document.getElementById('adsCount');
    showSkeletons();
    setTimeout(() => {
        const ads = generateMockAds();
        const filteredAds = applyFilters(ads);
        
        renderAds(filteredAds);
        updateAdsCount(filteredAds.length);
    }, 1000);
}

function showSkeletons() {
    const adsGrid = document.getElementById('adsGrid');
    adsGrid.innerHTML = `
        <div class="skeleton-ad"></div>
        <div class="skeleton-ad"></div>
        <div class="skeleton-ad"></div>
        <div class="skeleton-ad"></div>
        <div class="skeleton-ad"></div>
        <div class="skeleton-ad"></div>
    `;
}

function generateMockAds() {
    return [];
}

function applyFilters(ads) {
    const category = document.getElementById('categoryFilter').value;
    const city = document.getElementById('cityFilter').value;
    const minPrice = parseInt(document.getElementById('minPrice').value) || 0;
    const maxPrice = parseInt(document.getElementById('maxPrice').value) || Infinity;
    const newCondition = document.getElementById('newCondition').checked;
    const usedCondition = document.getElementById('usedCondition').checked;
    const sortBy = document.getElementById('sortFilter').value;
    
    let filteredAds = ads.filter(ad => {
        if (category && ad.category !== category) return false;
        if (city && ad.city !== city) return false;
        if (ad.price < minPrice || ad.price > maxPrice) return false;
        if (newCondition && usedCondition) return true;
        if (newCondition && ad.condition !== 'new') return false;
        if (usedCondition && ad.condition !== 'used') return false;
        return true;
    });
    switch (sortBy) {
        case 'newest':
            filteredAds.sort((a, b) => b.date - a.date);
            break;
        case 'oldest':
            filteredAds.sort((a, b) => a.date - b.date);
            break;
        case 'price-low':
            filteredAds.sort((a, b) => a.price - b.price);
            break;
        case 'price-high':
            filteredAds.sort((a, b) => b.price - a.price);
            break;
        case 'popular':
            filteredAds.sort((a, b) => b.views - a.views);
            break;
    }
    
    return filteredAds;
}

function renderAds(ads) {
    const adsGrid = document.getElementById('adsGrid');
    
    if (ads.length === 0) {
        adsGrid.innerHTML = `
            <div style="grid-column: 1 / -1; text-align: center; padding: 40px; color: var(--text-secondary);">
                <h3>Объявления не найдены</h3>
                <p>Попробуйте изменить параметры поиска</p>
            </div>
        `;
        return;
    }
    
    adsGrid.innerHTML = ads.map(ad => `
        <div class="ad-card" data-ad-id="${ad.id}">
            <img src="${ad.image}" alt="${ad.title}" class="ad-image" onerror="this.src='./pictures/placeholder.jpg'">
            <div class="ad-content">
                <h3 class="ad-title">${ad.title}</h3>
                <div class="ad-price">${ad.price.toLocaleString()} руб.</div>
                <div class="ad-location">
                    <span>📍</span>
                    <span>${getCityName(ad.city)}</span>
                </div>
                <div class="ad-meta">
                    <div class="ad-date">
                        <span>📅</span>
                        <span>${formatDate(ad.date)}</span>
                    </div>
                    <div class="ad-views">
                        <span>👁</span>
                        <span>${ad.views}</span>
                    </div>
                </div>
            </div>
        </div>
    `).join('');
    document.querySelectorAll('.ad-card').forEach(card => {
        card.addEventListener('click', () => {
            const adId = card.dataset.adId;
            window.location.href = `./ad.html?id=${adId}`;
        });
    });
}

function getCityName(cityCode) {
    const cities = {
        minsk: 'Минск',
        brest: 'Брест',
        vitebsk: 'Витебск',
        gomel: 'Гомель',
        grodno: 'Гродно',
        mogilev: 'Могилёв'
    };
    return cities[cityCode] || cityCode;
}

function formatDate(date) {
    const now = new Date();
    const diffTime = Math.abs(now - date);
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays === 1) return 'Вчера';
    if (diffDays < 7) return `${diffDays} дн. назад`;
    return date.toLocaleDateString('ru-RU');
}

function updateAdsCount(count) {
    const adsCount = document.getElementById('adsCount');
    adsCount.textContent = `Найдено: ${count} объявлений`;
}
function initLoadMore() {
    const loadMoreBtn = document.getElementById('loadMoreBtn');
    
    loadMoreBtn?.addEventListener('click', () => {
        console.log('Загрузка дополнительных объявлений...');
        loadMoreBtn.textContent = 'Загрузка...';
        loadMoreBtn.disabled = true;
        
        setTimeout(() => {
            loadMoreBtn.textContent = 'Показать ещё';
            loadMoreBtn.disabled = false;
        }, 1000);
    });
}
