document.addEventListener('DOMContentLoaded', function() {
    initChatsList();
    initChatWindow();
    initMessageInput();
    initSearch();
    initNewChat();
});
function initChatsList() {
    loadChats();
}
function loadChats() {
    const chatsList = document.getElementById('chatsList');
    showChatsSkeletons();
    setTimeout(() => {
        const chats = generateMockChats();
        renderChats(chats);
    }, 1000);
}
function showChatsSkeletons() {
    const chatsList = document.getElementById('chatsList');
    chatsList.innerHTML = `
        <div class="chat-item skeleton">
            <div class="chat-avatar" style="background: var(--skeleton-bg);"></div>
            <div class="chat-info">
                <div class="chat-user-name" style="background: var(--skeleton-bg); height: 16px; border-radius: 4px; margin-bottom: 8px;"></div>
                <div class="chat-last-message" style="background: var(--skeleton-bg); height: 12px; border-radius: 4px;"></div>
            </div>
        </div>
        <div class="chat-item skeleton">
            <div class="chat-avatar" style="background: var(--skeleton-bg);"></div>
            <div class="chat-info">
                <div class="chat-user-name" style="background: var(--skeleton-bg); height: 16px; border-radius: 4px; margin-bottom: 8px;"></div>
                <div class="chat-last-message" style="background: var(--skeleton-bg); height: 12px; border-radius: 4px;"></div>
            </div>
        </div>
        <div class="chat-item skeleton">
            <div class="chat-avatar" style="background: var(--skeleton-bg);"></div>
            <div class="chat-info">
                <div class="chat-user-name" style="background: var(--skeleton-bg); height: 16px; border-radius: 4px; margin-bottom: 8px;"></div>
                <div class="chat-last-message" style="background: var(--skeleton-bg); height: 12px; border-radius: 4px;"></div>
            </div>
        </div>
    `;
}
function generateMockChats() {
    return [];
}
function renderChats(chats) {
    const chatsList = document.getElementById('chatsList');
    if (chats.length === 0) {
        chatsList.innerHTML = `
            <div style="text-align: center; padding: 40px; color: var(--text-secondary);">
                <div style="font-size: 2rem; margin-bottom: 15px;">💬</div>
                <h3>Нет сообщений</h3>
                <p>Начните переписку с продавцами</p>
            </div>
        `;
        return;
    }
    chatsList.innerHTML = chats.map(chat => `
        <div class="chat-item" data-chat-id="${chat.id}">
            <img src="${chat.avatar}" alt="${chat.name}" class="chat-avatar" onerror="this.src='./pictures/profile.png'">
            <div class="chat-info">
                <div class="chat-user-name">${chat.name}</div>
                <div class="chat-last-message">${chat.lastMessage}</div>
                <div class="chat-meta">
                    <span class="chat-time">${formatTime(chat.time)}</span>
                    ${chat.unread > 0 ? `<span class="chat-unread-badge">${chat.unread}</span>` : ''}
                </div>
            </div>
        </div>
    `).join('');
    document.querySelectorAll('.chat-item').forEach(item => {
        item.addEventListener('click', () => {
            const chatId = item.dataset.chatId;
            openChat(chatId, chats.find(chat => chat.id == chatId));
        });
    });
}
function formatTime(date) {
    const now = new Date();
    const diffTime = now - date;
    const diffHours = Math.floor(diffTime / (1000 * 60 * 60));
    const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));
    if (diffHours < 1) return 'только что';
    if (diffHours < 24) return `${diffHours}ч назад`;
    if (diffDays === 1) return 'вчера';
    if (diffDays < 7) return `${diffDays} дн. назад`;
    return date.toLocaleDateString('ru-RU');
}
function initChatWindow() {
}
function openChat(chatId, chatData) {
    document.querySelectorAll('.chat-item').forEach(item => {
        item.classList.remove('active');
    });
    document.querySelector(`[data-chat-id="${chatId}"]`).classList.add('active');
    document.getElementById('chatPlaceholder').style.display = 'none';
    document.getElementById('chatWindow').style.display = 'flex';
    document.getElementById('chatUserAvatar').src = chatData.avatar;
    document.getElementById('chatUserName').textContent = chatData.name;
    document.getElementById('chatUserStatus').textContent = chatData.online ? 'онлайн' : 'был(а) недавно';
    loadMessages(chatId);
}
function loadMessages(chatId) {
    const chatMessages = document.getElementById('chatMessages');
    showMessagesSkeletons();
    setTimeout(() => {
        const messages = generateMockMessages(chatId);
        renderMessages(messages);
    }, 500);
}
function showMessagesSkeletons() {
    const chatMessages = document.getElementById('chatMessages');
    chatMessages.innerHTML = `
        <div class="message skeleton">
            <div class="message-avatar" style="background: var(--skeleton-bg); width: 32px; height: 32px; border-radius: 50%;"></div>
            <div class="message-content" style="background: var(--skeleton-bg); height: 40px; border-radius: 18px;"></div>
        </div>
        <div class="message own skeleton">
            <div class="message-content" style="background: var(--skeleton-bg); height: 40px; border-radius: 18px; margin-left: auto; max-width: 200px;"></div>
        </div>
        <div class="message skeleton">
            <div class="message-avatar" style="background: var(--skeleton-bg); width: 32px; height: 32px; border-radius: 50%;"></div>
            <div class="message-content" style="background: var(--skeleton-bg); height: 40px; border-radius: 18px;"></div>
        </div>
    `;
}
function generateMockMessages(chatId) {
    return [];
}
function renderMessages(messages) {
    const chatMessages = document.getElementById('chatMessages');
    chatMessages.innerHTML = messages.map(message => `
        <div class="message ${message.own ? 'own' : ''}">
            ${!message.own ? `<img src="./pictures/profile.png" alt="Аватар" class="message-avatar" onerror="this.src='./pictures/profile.png'">` : ''}
            <div class="message-content">
                <p class="message-text">${message.text}</p>
                <div class="message-time">${formatMessageTime(message.time)}</div>
                ${message.status ? `<div class="message-status">${getStatusIcon(message.status)}</div>` : ''}
            </div>
        </div>
    `).join('');
    chatMessages.scrollTop = chatMessages.scrollHeight;
}
function formatMessageTime(date) {
    return date.toLocaleTimeString('ru-RU', { 
        hour: '2-digit', 
        minute: '2-digit' 
    });
}
function getStatusIcon(status) {
    const icons = {
        sent: '✓',
        delivered: '✓✓',
        read: '✓✓'
    };
    return icons[status] || '';
}
function initMessageInput() {
    const messageInput = document.getElementById('chatMessageInput');
    const sendBtn = document.getElementById('chatSendBtn');
    messageInput?.addEventListener('keypress', (e) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            sendMessage();
        }
    });
    sendBtn?.addEventListener('click', sendMessage);
}
function sendMessage() {
    const messageInput = document.getElementById('chatMessageInput');
    const text = messageInput.value.trim();
    if (!text) return;
    addMessageToChat(text, true);
    messageInput.value = '';
}
function addMessageToChat(text, isOwn) {
    const chatMessages = document.getElementById('chatMessages');
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${isOwn ? 'own' : ''}`;
    const now = new Date();
    const timeString = now.toLocaleTimeString('ru-RU', { 
        hour: '2-digit', 
        minute: '2-digit' 
    });
    messageDiv.innerHTML = `
        ${!isOwn ? `<img src="./pictures/profile.png" alt="Аватар" class="message-avatar" onerror="this.src='./pictures/profile.png'">` : ''}
        <div class="message-content">
            <p class="message-text">${text}</p>
            <div class="message-time">${timeString}</div>
            ${isOwn ? '<div class="message-status">✓</div>' : ''}
        </div>
    `;
    chatMessages.appendChild(messageDiv);
    chatMessages.scrollTop = chatMessages.scrollHeight;
}
function initSearch() {
    const searchInput = document.getElementById('chatsSearch');
    searchInput?.addEventListener('input', (e) => {
        const query = e.target.value.toLowerCase();
        filterChats(query);
    });
}
function filterChats(query) {
    const chatItems = document.querySelectorAll('.chat-item');
    chatItems.forEach(item => {
        const name = item.querySelector('.chat-user-name').textContent.toLowerCase();
        const message = item.querySelector('.chat-last-message').textContent.toLowerCase();
        if (name.includes(query) || message.includes(query)) {
            item.style.display = 'flex';
        } else {
            item.style.display = 'none';
        }
    });
}
function initNewChat() {
    const newChatBtn = document.getElementById('newChatBtn');
    newChatBtn?.addEventListener('click', () => {
        alert('potom');
    });
}
