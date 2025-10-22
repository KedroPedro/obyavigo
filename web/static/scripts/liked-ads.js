document.addEventListener('DOMContentLoaded', function() {
    const urlParams = new URLSearchParams(window.location.search);
    const adId = urlParams.get('id');
    if (!adId) {
        // Для тестирования без редиректа — просто показываем демо
        loadMockAd();
        return;
    }
    loadAdData(adId);
    loadSimilarAds(adId);
});

function loadMockAd() {
    const mockAd = {
        title: 'Apple iPhone 15 Pro 256GB',
        price: 3200,
        city: 'Минск',
        district: 'Центральный',
        created_at: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
        views: 127,
        description: 'Продается новый iPhone 15 Pro в оригинальной упаковке. Куплен в официальном магазине. Все аксессуары на месте. Гарантия до декабря 2025 года.\n\n🔹 Цвет: Титан черный\n🔹 Память: 256 ГБ\n🔹 Состояние: Новое\n🔹 Гарантия: До 2025 года',
        images: [
            './pictures/iphone1.jpg',
            './pictures/iphone2.jpg',
            './pictures/iphone3.jpg'
        ],
        seller_name: 'Александр',
        seller_online: true,
        seller_email: 'alexander@example.com',
        seller_avatar: './pictures/profile.png'
    };
    
    renderAd(mockAd);
    loadImages(mockAd.images || []);
    document.getElementById('favoriteBtn').classList.remove('active');
}

function loadAdData(adId) {
    // Здесь можно оставить fetch, но для демо — сразу мок
    loadMockAd();
}

function renderAd(ad) {
    document.title = `${ad.title} — Obyavigo`;
    document.getElementById('adTitle').textContent = ad.title;
    document.getElementById('adPrice').textContent = `${ad.price.toLocaleString('ru-RU')} BYN`;
    
    const metaHtml = `
        <span>📍 ${ad.city}, ${ad.district}</span>
        <span>🕐 ${formatDate(ad.created_at)}</span>
        <span>👁️ ${ad.views} просмотров</span>
    `;
    document.getElementById('adMeta').innerHTML = metaHtml;
    document.getElementById('adDescription').textContent = ad.description;
    
    const authorHtml = `
        <div class="author-info">
            <img src="${ad.seller_avatar || './pictures/profile.png'}" alt="Автор" class="author-avatar" loading="lazy">
            <div>
                <div class="author-name">${ad.seller_name}</div>
                <div class="author-online">${ad.seller_online ? 'online' : 'offline'}</div>
            </div>
        </div>
        <div class="author-contacts">
            <div class="contact-item">
                <span>📞</span>
                <span id="phoneDisplay">Показать телефон</span>
                <button class="show-phone-btn" id="showPhoneBtn">Показать</button>
            </div>
            ${ad.seller_email ? `<div class="contact-item"><span>📧</span> ${ad.seller_email}</div>` : ''}
        </div>
    `;
    document.getElementById('adAuthor').innerHTML = authorHtml;
    
    // Инициализация кнопок
    setTimeout(() => {
        document.getElementById('showPhoneBtn')?.addEventListener('click', function() {
            document.getElementById('phoneDisplay').textContent = '+375 (29) 123-45-67';
            this.style.display = 'none';
        });
        
        document.getElementById('favoriteBtn')?.addEventListener('click', function() {
            this.classList.toggle('active');
            this.innerHTML = this.classList.contains('active') 
                ? '<span>❤️</span><span>В избранном</span>' 
                : '<span>❤️</span><span>В избранное</span>';
        });
        
        document.getElementById('contactBtn')?.addEventListener('click', () => {
            alert('Написать продавцу');
        });
    }, 100);
}

function loadImages(imageUrls) {
    const mainImage = document.getElementById('mainImage');
    const thumbnailContainer = document.getElementById('thumbnailContainer');
    
    if (imageUrls.length === 0) return;
    
    mainImage.src = imageUrls[0];
    mainImage.onload = () => { mainImage.style.opacity = '1'; };
    
    thumbnailContainer.innerHTML = '';
    imageUrls.forEach((url, index) => {
        const img = document.createElement('img');
        img.src = url;
        img.alt = `Изображение ${index + 1}`;
        img.className = 'ad-thumbnail';
        img.loading = 'lazy';
        if (index === 0) img.classList.add('active');
        img.addEventListener('click', () => {
            mainImage.src = url;
            document.querySelectorAll('.ad-thumbnail').forEach(thumb => {
                thumb.classList.remove('active');
            });
            img.classList.add('active');
        });
        thumbnailContainer.appendChild(img);
    });
}

function loadSimilarAds(adId) {
    const similarAds = [
        { id: '1', title: 'Samsung Galaxy S24 Ultra', price: 2800, city: 'Минск', image: './pictures/samsung.jpg' },
        { id: '2', title: 'Google Pixel 8 Pro', price: 2400, city: 'Брест', image: './pictures/pixel.jpg' },
        { id: '3', title: 'iPhone 14 Pro Max', price: 2900, city: 'Гомель', image: './pictures/iphone14.jpg' }
    ];
    
    const grid = document.getElementById('similarAdsGrid');
    grid.innerHTML = similarAds.map(ad => `
        <div class="similar-ad-card">
            <img src="${ad.image}" alt="${ad.title}" class="similar-ad-image" loading="lazy">
            <div class="similar-ad-content">
                <h3 class="similar-ad-title">${ad.title}</h3>
                <div class="similar-ad-price">${ad.price.toLocaleString('ru-RU')} BYN</div>
                <div class="similar-ad-location">${ad.city}</div>
            </div>
        </div>
    `).join('');
}

function formatDate(dateString) {
    const date = new Date(dateString);
    const today = new Date();
    const yesterday = new Date(today);
    yesterday.setDate(yesterday.getDate() - 1);
    
    if (date.toDateString() === today.toDateString()) {
        return `Сегодня, ${date.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' })}`;
    } else if (date.toDateString() === yesterday.toDateString()) {
        return `Вчера, ${date.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' })}`;
    } else {
        return date.toLocaleDateString('ru-RU', { 
            day: 'numeric', 
            month: 'long', 
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    }
}