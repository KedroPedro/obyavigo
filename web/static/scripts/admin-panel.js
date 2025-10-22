document.addEventListener('DOMContentLoaded', function() {

    
    // Инициализация навигации
    initNavigation();
    
    // Загрузка данных
    loadDashboardData();
});



// Навигация по вкладкам
function initNavigation() {
    const navItems = document.querySelectorAll('.nav-item');
    const tabContents = document.querySelectorAll('.tab-content');
    const pageTitle = document.getElementById('pageTitle');
    
    navItems.forEach(item => {
        item.addEventListener('click', (e) => {
            e.preventDefault();
            
            // Убираем активный класс у всех
            navItems.forEach(i => i.classList.remove('active'));
            tabContents.forEach(t => t.classList.remove('active'));
            
            // Добавляем активный класс текущему
            item.classList.add('active');
            const tabId = item.dataset.tab + '-tab';
            document.getElementById(tabId).classList.add('active');
            
            // Обновляем заголовок
            const titles = {
                ads: 'Объявления',
                users: 'Пользователи',
                moderation: 'Модерация',
                reports: 'Жалобы',
                stats: 'Статистика'
            };
            pageTitle.textContent = titles[item.dataset.tab];
            
            // Загружаем данные для вкладки
            loadTabData(item.dataset.tab);
        });
    });
}

// Загрузка данных для вкладки
function loadTabData(tab) {
    switch(tab) {
        case 'ads':
            loadAds();
            break;
        case 'users':
            loadUsers();
            break;
        case 'moderation':
            loadModeration();
            break;
        case 'reports':
            loadReports();
            break;
        case 'stats':
            // Статистика уже загружена
            break;
    }
}

// Загрузка дашборда
function loadDashboardData() {
    // Статистика
    document.getElementById('totalAds').textContent = '24,567';
    document.getElementById('totalUsers').textContent = '18,932';
    document.getElementById('pendingReports').textContent = '12';
    document.getElementById('pendingModeration').textContent = '8';
    
    // Загружаем данные для первой вкладки (объявления)
    loadAds();
}

// Загрузка объявлений
function loadAds() {
    const tbody = document.getElementById('adsTableBody');
    tbody.innerHTML = `
        <tr>
            <td>12345</td>
            <td>iPhone 15 Pro Max 256GB</td>
            <td>Алексей</td>
            <td>3 200 BYN</td>
            <td>Минск</td>
            <td><span class="status active">Активное</span></td>
            <td>
                <button class="action-btn view">👁️</button>
                <button class="action-btn reject">❌</button>
            </td>
        </tr>
        <tr>
            <td>12344</td>
            <td>Квартира 2к, центр</td>
            <td>Мария</td>
            <td>320 000 BYN</td>
            <td>Брест</td>
            <td><span class="status pending">На модерации</span></td>
            <td>
                <button class="action-btn approve">✅</button>
                <button class="action-btn reject">❌</button>
            </td>
        </tr>
        <tr>
            <td>12343</td>
            <td>BMW X5 2020</td>
            <td>Иван</td>
            <td>89 000 BYN</td>
            <td>Гомель</td>
            <td><span class="status rejected">Отклонено</span></td>
            <td>
                <button class="action-btn view">👁️</button>
                <button class="action-btn ban">🚫</button>
            </td>
        </tr>
    `;
    
    // Добавляем обработчики кнопок
    addActionButtonListeners();
}

// Загрузка пользователей
function loadUsers() {
    const tbody = document.getElementById('usersTableBody');
    tbody.innerHTML = `
        <tr>
            <td>1001</td>
            <td>Алексей</td>
            <td>alexey@mail.ru</td>
            <td>+375291234567</td>
            <td>15.05.2024</td>
            <td><span class="status active">Активен</span></td>
            <td>
                <button class="action-btn ban">🚫</button>
            </td>
        </tr>
        <tr>
            <td>1002</td>
            <td>Мария</td>
            <td>maria@gmail.com</td>
            <td>+375297654321</td>
            <td>10.05.2024</td>
            <td><span class="status banned">Забанен</span></td>
            <td>
                <button class="action-btn view">👁️</button>
            </td>
        </tr>
    `;
    addActionButtonListeners();
}

// Загрузка модерации
function loadModeration() {
    const tbody = document.getElementById('moderationTableBody');
    tbody.innerHTML = `
        <tr>
            <td>12344</td>
            <td>Квартира 2к, центр</td>
            <td>Мария</td>
            <td>Сегодня, 14:30</td>
            <td>
                <button class="action-btn approve">✅ Одобрить</button>
                <button class="action-btn reject">❌ Отклонить</button>
            </td>
        </tr>
        <tr>
            <td>12346</td>
            <td>MacBook Pro M3</td>
            <td>Дмитрий</td>
            <td>Сегодня, 12:15</td>
            <td>
                <button class="action-btn approve">✅ Одобрить</button>
                <button class="action-btn reject">❌ Отклонить</button>
            </td>
        </tr>
    `;
    addActionButtonListeners();
}

// Загрузка жалоб
function loadReports() {
    const tbody = document.getElementById('reportsTableBody');
    tbody.innerHTML = `
        <tr>
            <td>501</td>
            <td>iPhone 15 Pro Max</td>
            <td>Подозрение в мошенничестве</td>
            <td>Сегодня, 15:20</td>
            <td><span class="status pending">Новая</span></td>
            <td>
                <button class="action-btn view">👁️</button>
                <button class="action-btn ban">🚫</button>
            </td>
        </tr>
        <tr>
            <td>500</td>
            <td>Квартира в центре</td>
            <td>Неверная информация</td>
            <td>Вчера, 18:45</td>
            <td><span class="status active">Рассмотрена</span></td>
            <td>
                <button class="action-btn view">👁️</button>
            </td>
        </tr>
    `;
    addActionButtonListeners();
}

// Обработчики кнопок действий
function addActionButtonListeners() {
    document.querySelectorAll('.action-btn').forEach(btn => {
        btn.addEventListener('click', function() {
            const action = this.className.includes('approve') ? 'approve' :
                          this.className.includes('reject') ? 'reject' :
                          this.className.includes('ban') ? 'ban' : 'view';
            
            if (action === 'view') {
                alert('Просмотр объявления/пользователя');
                return;
            }
            
            showConfirmModal(action, () => {
                // Здесь будет запрос к бэкенду
                console.log(`Выполнено действие: ${action}`);
                alert(`Действие "${action}" выполнено успешно!`);
            });
        });
    });
}

// Модальное окно подтверждения
function showConfirmModal(action, callback) {
    const modal = document.getElementById('confirmModal');
    const title = document.getElementById('modalTitle');
    const message = document.getElementById('modalMessage');
    const confirmBtn = document.getElementById('confirmBtn');
    
    const texts = {
        approve: { title: 'Одобрение объявления', message: 'Вы уверены, что хотите одобрить это объявление?' },
        reject: { title: 'Отклонение объявления', message: 'Вы уверены, что хотите отклонить это объявление?' },
        ban: { title: 'Блокировка', message: 'Вы уверены, что хотите заблокировать пользователя?' }
    };
    
    title.textContent = texts[action].title;
    message.textContent = texts[action].message;
    modal.style.display = 'flex';
    
    document.getElementById('cancelBtn').onclick = () => {
        modal.style.display = 'none';
    };
    
    confirmBtn.onclick = () => {
        modal.style.display = 'none';
        callback();
    };
    
    modal.onclick = (e) => {
        if (e.target === modal) {
            modal.style.display = 'none';
        }
    };
}