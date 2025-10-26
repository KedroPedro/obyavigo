// Функция переключения избранного
function toggleFavorite(button, adId) {
    const card = button.closest('.ad-card');
    const isActive = button.classList.contains('active');
    
    if (isActive) {
        // Удаляем из избранного
        fetch('/api/liked-ads/remove', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ ad_id: adId }),
            credentials: 'include'
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                // Плавно удаляем карточку
                card.style.transition = 'opacity 0.3s ease, transform 0.3s ease';
                card.style.opacity = '0';
                card.style.transform = 'scale(0.9)';
                
                setTimeout(() => {
                    card.remove();
                    updateLikedCount();
                    
                    // Проверяем, остались ли объявления
                    const grid = document.getElementById('likedAdsGrid');
                    const remainingCards = grid.querySelectorAll('.ad-card');
                    
                    if (remainingCards.length === 0) {
                        // Показываем пустое состояние
                        grid.innerHTML = `
                            <div class="no-ads">
                                <div class="no-ads-content">
                                    <div class="no-ads-icon">❤️</div>
                                    <h3>Пока нет избранных объявлений</h3>
                                    <p>Добавьте объявления в избранное, чтобы они отображались здесь</p>
                                    <a href="/ads/" class="btn-primary">Посмотреть объявления</a>
                                </div>
                            </div>
                        `;
                    }
                }, 300);
            }
        })
        .catch(error => {
            console.error('Error removing from favorites:', error);
        });
    }
}

// Обновление счетчика избранных
function updateLikedCount() {
    const grid = document.getElementById('likedAdsGrid');
    const count = grid.querySelectorAll('.ad-card').length;
    const countElement = document.getElementById('likedAdsCount');
    
    let word = 'объявлений';
    if (count === 1) {
        word = 'объявление';
    } else if (count >= 2 && count <= 4) {
        word = 'объявления';
    }
    
    countElement.textContent = `В избранном: ${count} ${word}`;
}

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', function() {
    console.log('Liked ads page loaded');
});
