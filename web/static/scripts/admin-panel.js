let currentUserRole = null;
document.addEventListener('DOMContentLoaded', function() {
    initNavigation();
    const usersRoleFilter = document.getElementById('usersRoleFilter');
    const usersSearch = document.getElementById('usersSearch');
    if (usersRoleFilter) {
        usersRoleFilter.addEventListener('change', () => loadUsers());
    }
    if (usersSearch) {
        let debounceTimer;
        usersSearch.addEventListener('input', () => {
            clearTimeout(debounceTimer);
            debounceTimer = setTimeout(() => loadUsers(), 500);
        });
    }
    const detailsCloseBtn = document.getElementById('detailsCloseBtn');
    if (detailsCloseBtn) {
        detailsCloseBtn.onclick = () => {
            document.getElementById('detailsModal').style.display = 'none';
        };
    }
    getUserRole().then(() => {
        loadDashboardData();
    });
});
function initNavigation() {
    const navItems = document.querySelectorAll('.nav-item');
    const tabContents = document.querySelectorAll('.tab-content');
    const pageTitle = document.getElementById('pageTitle');
    const sidebar = document.getElementById('sidebar');
    const sidebarOverlay = document.getElementById('sidebarOverlay');
    const sidebarToggle = document.getElementById('sidebarToggle');
    const mobileTabSelect = document.getElementById('mobileTabSelect');
    const titles = {
        users: 'Пользователи',
        moderation: 'Модерация',
        reports: 'Жалобы',
        stats: 'Статистика'
    };

    const closeSidebar = () => {
        sidebar?.classList.remove('active');
        sidebarOverlay?.classList.remove('active');
        document.body.classList.remove('sidebar-open');
    };

    const openSidebar = () => {
        sidebar?.classList.add('active');
        sidebarOverlay?.classList.add('active');
        document.body.classList.add('sidebar-open');
    };

    const activateTab = (tabKey) => {
        if (!tabKey) return;
        navItems.forEach(i => i.classList.toggle('active', i.dataset.tab === tabKey));
        tabContents.forEach(t => t.classList.toggle('active', t.id === `${tabKey}-tab`));
        pageTitle.textContent = titles[tabKey] || 'Админ-панель';
        if (mobileTabSelect && mobileTabSelect.value !== tabKey) {
            mobileTabSelect.value = tabKey;
        }
        loadTabData(tabKey);
        if (window.innerWidth <= 768) {
            closeSidebar();
        }
    };

    navItems.forEach(item => {
        item.addEventListener('click', (e) => {
            e.preventDefault();
            activateTab(item.dataset.tab);
        });
    });

    mobileTabSelect?.addEventListener('change', (e) => activateTab(e.target.value));
    sidebarToggle?.addEventListener('click', openSidebar);
    sidebarOverlay?.addEventListener('click', closeSidebar);
    window.addEventListener('resize', () => {
        if (window.innerWidth > 768) {
            closeSidebar();
        }
    });

    if (mobileTabSelect) {
        const activeItem = document.querySelector('.nav-item.active');
        if (activeItem) {
            mobileTabSelect.value = activeItem.dataset.tab;
        }
    }
}
function loadTabData(tab) {
    switch(tab) {
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
            loadStats();
            break;
    }
}
async function getUserRole() {
    try {
        const res = await fetch('/api/user/profile/');
        if (!res.ok) throw new Error('Failed to get user profile');
        const data = await res.json();
        currentUserRole = data.role;
    } catch (err) {
        console.error('Error getting user role:', err);
        currentUserRole = 'user';
    }
}
function showDetailsModal(title, content) {
    const modal = document.getElementById('detailsModal');
    document.getElementById('detailsModalTitle').textContent = title;
    document.getElementById('detailsModalContent').textContent = content;
    modal.style.display = 'flex';
    modal.onclick = (e) => {
        if (e.target === modal) {
            modal.style.display = 'none';
        }
    };
}
function createTruncatedText(text, maxLength, title) {
    if (!text) return '-';
    if (text.length <= maxLength) return text;
    const truncated = text.substring(0, maxLength) + '...';
    return `<span class="truncate-text" onclick="showDetailsModal('${title}', ${JSON.stringify(text).replace(/"/g, '&quot;')})">${truncated}</span>`;
}
function loadStats() {
    fetch('/api/admin/stats/')
        .then(res => res.json())
        .then(data => {
            document.getElementById('totalAds').textContent = data.totalAds || '0';
            document.getElementById('totalUsers').textContent = data.totalUsers || '0';
            document.getElementById('pendingReports').textContent = data.pendingReports || '0';
            document.getElementById('pendingModeration').textContent = data.pendingModeration || '0';
        })
        .catch(err => console.error('Error loading stats:', err));
}

function loadDashboardData() {
    loadStats();
    loadModeration();
}
function loadAds() {
    const tbody = document.getElementById('adsTableBody');
    const statusFilter = document.getElementById('adsStatusFilter')?.value || 'all';
    const searchQuery = document.getElementById('adsSearch')?.value || '';
    tbody.innerHTML = '<tr><td colspan="7" class="loading">Загрузка...</td></tr>';
    fetch(`/api/admin/ads/?status=${statusFilter}&search=${encodeURIComponent(searchQuery)}`)
        .then(res => res.json())
        .then(data => {
            if (!data.ads || data.ads.length === 0) {
                tbody.innerHTML = '<tr><td colspan="7" class="no-data">Нет данных</td></tr>';
                return;
            }
            tbody.innerHTML = data.ads.map(ad => `
                <tr>
                    <td>${ad.id.substring(0, 8)}...</td>
                    <td>${ad.title}</td>
                    <td>${ad.user_id ? ad.user_id.substring(0, 8) + '...' : '-'}</td>
                    <td>${ad.price} BYN</td>
                    <td>${ad.location_name || '-'}</td>
                    <td><span class="status ${ad.ad_status}">${getStatusText(ad.ad_status)}</span></td>
                    <td>
                        <button class="action-btn view" onclick="viewAd('${ad.id}')">👁️</button>
                        ${ad.ad_status !== 'public' ? `<button class="action-btn approve" onclick="updateAdStatus('${ad.id}', 'public')">✅</button>` : ''}
                        ${ad.ad_status !== 'rejected' ? `<button class="action-btn reject" onclick="updateAdStatus('${ad.id}', 'rejected')">❌</button>` : ''}
                    </td>
                </tr>
            `).join('');
        })
        .catch(err => {
            console.error('Error loading ads:', err);
            tbody.innerHTML = '<tr><td colspan="7" class="no-data">Ошибка загрузки данных</td></tr>';
        });
}
function getStatusText(status) {
    const statusMap = {
        'public': 'Активное',
        'draft': 'Черновик',
        'rejected': 'Отклонено',
        'pending': 'На модерации'
    };
    return statusMap[status] || status;
}
function viewAd(adId) {
    window.open(`/ads/${adId}/`, '_blank');
}
function updateAdStatus(adId, status) {
    showConfirmModal(status === 'public' ? 'approve' : 'reject', () => {
        fetch(`/api/admin/ads/${adId}/status/`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ status })
        })
        .then(res => {
            if (!res.ok) throw new Error('Failed to update status');
            return res.json();
        })
        .then(() => {
            loadModeration();
            loadStats();
        })
        .catch(err => {
            console.error('Error updating ad status:', err);
        });
    });
}
function loadUsers() {
    const tbody = document.getElementById('usersTableBody');
    const roleFilter = document.getElementById('usersRoleFilter')?.value || 'all';
    const searchQuery = document.getElementById('usersSearch')?.value || '';
    tbody.innerHTML = '<tr><td colspan="7" class="loading">Загрузка...</td></tr>';
    fetch(`/api/admin/users/?role=${roleFilter}&search=${encodeURIComponent(searchQuery)}`)
        .then(res => res.json())
        .then(data => {
            if (!data.users || data.users.length === 0) {
                tbody.innerHTML = '<tr><td colspan="7" class="no-data">Нет данных</td></tr>';
                return;
            }
            tbody.innerHTML = data.users.map(user => `
                <tr>
                    <td>${createTruncatedText(user.id, 12, 'ID пользователя')}</td>
                    <td>${createTruncatedText(user.username, 20, 'Имя пользователя')}</td>
                    <td>${createTruncatedText(user.email, 30, 'Email')}</td>
                    <td>${user.role}</td>
                    <td><span class="status ${user.status}">${user.status === 'active' ? 'Активен' : 'Забанен'}</span></td>
                    <td>${new Date(user.registrationDate).toLocaleDateString('ru-RU')}</td>
                    <td class="action-btns">
                        ${user.status === 'active' ? 
                            `<button class="action-btn ban" onclick="updateUserStatus('${user.id}', 'banned')">🚫 Забанить</button>` :
                            `<button class="action-btn approve" onclick="updateUserStatus('${user.id}', 'active')">✅ Разбанить</button>`
                        }
                        ${currentUserRole === 'admin' && user.role === 'user' ? 
                            `<button class=\"btn-primary\" onclick=\"grantModeratorRole('${user.id}')\">⭐ Сделать модератором</button>` : ''
                        }
                        ${currentUserRole === 'admin' && user.role === 'moderator' ? 
                            `<button class=\"btn-primary\" onclick=\"grantAdminRole('${user.id}')\">👑 Сделать администратором</button>` : ''
                        }
                        ${currentUserRole === 'admin' && user.role === 'moderator' ? 
                            `<button class=\"btn-warning\" onclick=\"demoteToUser('${user.id}')\">⬇️ Понизить до пользователя</button>` : ''
                        }
                        ${currentUserRole === 'admin' && user.role === 'admin' ? 
                            `<button class=\"btn-warning\" onclick=\"demoteToUser('${user.id}')\">⬇️ Понизить до пользователя</button>` : ''
                        }
                    </td>
                </tr>
            `).join('');
        })
        .catch(err => {
            console.error('Error loading users:', err);
            tbody.innerHTML = '<tr><td colspan="7" class="no-data">Ошибка загрузки данных</td></tr>';
        });
}
function updateUserStatus(userId, status) {
    showConfirmModal('ban', () => {
        fetch(`/api/admin/users/${userId}/status/`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ status })
        })
        .then(res => {
            if (!res.ok) throw new Error('Failed to update status');
            return res.json();
        })
        .then(() => {
            loadUsers();
        })
        .catch(err => {
            console.error('Error updating user status:', err);
        });
    });
}
function loadModeration() {
    const tbody = document.getElementById('moderationTableBody');
    tbody.innerHTML = '<tr><td colspan="6" class="loading">Загрузка...</td></tr>';
    fetch('/api/admin/ads/?status=moderation')
        .then(res => res.json())
        .then(data => {
            if (!data.ads || data.ads.length === 0) {
                tbody.innerHTML = '<tr><td colspan="6" class="no-data">Нет объявлений на модерации</td></tr>';
                return;
            }
            tbody.innerHTML = data.ads.map(ad => `
                <tr>
                    <td>${createTruncatedText(ad.id || '-', 12, 'ID объявления')}</td>
                    <td>${createTruncatedText(ad.title, 40, 'Заголовок')}</td>
                    <td>${createTruncatedText(ad.user_id || '-', 12, 'ID пользователя')}</td>
                    <td>${createTruncatedText(ad.user_email || '-', 30, 'Email пользователя')}</td>
                    <td>${new Date(ad.created_at).toLocaleDateString('ru-RU')}</td>
                    <td class="action-btns">
                        <button class="action-btn view" onclick="viewAd('${ad.id}')" title="Перейти к объявлению">👁️</button>
                        <button class="action-btn approve" onclick="updateAdStatus('${ad.id}', 'public')">✅ Одобрить</button>
                        <button class="action-btn reject" onclick="updateAdStatus('${ad.id}', 'rejected')">❌ Отклонить</button>
                    </td>
                </tr>
            `).join('');
        })
        .catch(err => {
            console.error('Error loading moderation:', err);
            tbody.innerHTML = '<tr><td colspan="6" class="no-data">Ошибка загрузки данных</td></tr>';
        });
}
function loadReports() {
    const tbody = document.getElementById('reportsTableBody');
    const statusFilter = document.getElementById('reportsStatusFilter')?.value || 'all';
    tbody.innerHTML = '<tr><td colspan="7" class="loading">Загрузка...</td></tr>';
    fetch(`/api/admin/reports/?status=${statusFilter}`)
        .then(res => res.json())
        .then(data => {
            if (!data.reports || data.reports.length === 0) {
                tbody.innerHTML = '<tr><td colspan="7" class="no-data">Нет жалоб</td></tr>';
                return;
            }
            tbody.innerHTML = data.reports.map(report => `
                <tr>
                    <td>${createTruncatedText(report.id, 12, 'ID жалобы')}</td>
                    <td>${getReportTypeText(report.complaint_type)}</td>
                    <td>${report.listing_id ? `<a href="/ads/${report.listing_id}/" target="_blank" style="color: var(--primary-color);">${report.listing_id.substring(0, 8)}...</a>` : '-'}</td>
                    <td>${createTruncatedText(report.complainant_email || '-', 25, 'Email отправителя')}</td>
                    <td>${createTruncatedText(report.description, 40, 'Причина')}</td>
                    <td>${new Date(report.created_at).toLocaleDateString('ru-RU')}</td>
                    <td class="action-btns">
                        ${report.listing_id ? `<button class="action-btn view" onclick="viewAd('${report.listing_id}')" title="Перейти к объявлению">👁️</button>` : ''}
                        ${report.status === 'pending' ? `
                            <button class="action-btn reject" onclick="rejectReport('${report.id}')">❌ Отклонить</button>
                            <button class="action-btn ban" onclick="blockAdFromReport('${report.id}', '${report.listing_id}')">🚫 Заблокировать объявление</button>
                            <button class="btn-danger" onclick="blockAdAndUserFromReport('${report.id}', '${report.listing_id}', '${report.ad_owner_id}')">🚨 Заблокировать объявление и пользователя</button>
                        ` : '<span class="status resolved">Решено</span>'}
                    </td>
                </tr>
            `).join('');
        })
        .catch(err => {
            console.error('Error loading reports:', err);
            tbody.innerHTML = '<tr><td colspan="7" class="no-data">Ошибка загрузки данных</td></tr>';
        });
}
function getReportTypeText(type) {
    const types = {
        'spam': 'Спам',
        'fraud': 'Мошенничество',
        'fake': 'Поддельный товар',
        'wrong_category': 'Неверная категория',
        'offensive': 'Оскорбительный контент',
        'sold': 'Товар уже продан',
        'other': 'Другое'
    };
    return types[type] || type;
}
function rejectReport(reportId) {
    showConfirmModal('reject', () => {
        fetch(`/api/admin/reports/${reportId}/status/`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ status: 'rejected' })
        })
        .then(res => {
            if (!res.ok) throw new Error('Failed to reject report');
            return res.json();
        })
        .then(() => {
            loadReports();
            loadDashboardData(); 
        })
        .catch(err => {
            console.error('Error rejecting report:', err);
        });
    });
}
function blockAdFromReport(reportId, adId) {
    showConfirmModal('ban', () => {
        Promise.all([
            fetch(`/api/admin/ads/${adId}/status/`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ status: 'rejected' })
            }),
            fetch(`/api/admin/reports/${reportId}/status/`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ status: 'resolved', resolution_comment: 'Ad blocked' })
            })
        ])
        .then(responses => {
            if (!responses[0].ok || !responses[1].ok) {
                throw new Error('Failed to block ad');
            }
            return Promise.all(responses.map(r => r.json()));
        })
        .then(() => {
            loadReports();
            loadDashboardData();
        })
        .catch(err => {
            console.error('Error blocking ad:', err);
        });
    });
}
function blockAdAndUserFromReport(reportId, adId, userId) {
    showConfirmModal('ban', () => {
        Promise.all([
            fetch(`/api/admin/ads/${adId}/status/`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ status: 'rejected' })
            }),
            fetch(`/api/admin/users/${userId}/status/`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ status: 'banned' })
            }),
            fetch(`/api/admin/reports/${reportId}/status/`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ status: 'resolved', resolution_comment: 'Ad and user blocked' })
            })
        ])
        .then(responses => {
            if (!responses[0].ok || !responses[1].ok || !responses[2].ok) {
                throw new Error('Failed to block ad and user');
            }
            return Promise.all(responses.map(r => r.json()));
        })
        .then(() => {
            loadReports();
            loadDashboardData();
        })
        .catch(err => {
            console.error('Error blocking ad and user:', err);
        });
    });
}
function grantModeratorRole(userId) {
    showConfirmModal('approve', () => {
        fetch(`/api/admin/users/${userId}/role/`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ role: 'moderator' })
        })
        .then(res => {
            if (!res.ok) throw new Error('Failed to grant moderator role');
            return res.json();
        })
        .then(() => {
            loadUsers();
        })
        .catch(err => {
            console.error('Error granting moderator role:', err);
        });
    });
}
function grantAdminRole(userId) {
    showConfirmModal('approve', () => {
        fetch(`/api/admin/users/${userId}/role/`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ role: 'admin' })
        })
        .then(res => {
            if (!res.ok) throw new Error('Failed to grant admin role');
            return res.json();
        })
        .then(() => {
            loadUsers();
        })
        .catch(err => {
            console.error('Error granting admin role:', err);
        });
    });
}
function demoteToUser(userId) {
    showConfirmModal('ban', () => {
        fetch(`/api/admin/users/${userId}/role/`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ role: 'user' })
        })
        .then(res => {
            if (!res.ok) throw new Error('Failed to demote user');
            return res.json();
        })
        .then(() => {
            loadUsers();
        })
        .catch(err => {
            console.error('Error demoting user:', err);
        });
    });
}
function showConfirmModal(action, callback) {
    const modal = document.getElementById('confirmModal');
    const title = document.getElementById('modalTitle');
    const message = document.getElementById('modalMessage');
    const confirmBtn = document.getElementById('confirmBtn');
    const texts = {
        approve: { title: 'Подтверждение', message: 'Вы уверены, что хотите выполнить это действие?' },
        reject: { title: 'Отклонение', message: 'Вы уверены, что хотите отклонить?' },
        ban: { title: 'Блокировка', message: 'Вы уверены, что хотите заблокировать?' }
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
