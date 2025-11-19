const MessagesApp = (() => {
    const state = {
        chats: [],
        selectedChatId: null,
        currentUser: null,
        ws: null,
        reconnectTimer: null,
        messagesCache: new Map(),
        messagesMeta: new Map(),
    };

    function init() {
        if (!document.getElementById('chatsList')) {
            return;
        }
        bindUI();
        loadInitialData();
    }

    function bindUI() {
        const messageInput = document.getElementById('chatMessageInput');
        const sendBtn = document.getElementById('chatSendBtn');
        const searchInput = document.getElementById('chatsSearch');
        const newChatBtn = document.getElementById('newChatBtn');

        messageInput?.addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                sendMessage();
            }
        });

        sendBtn?.addEventListener('click', sendMessage);
        searchInput?.addEventListener('input', (e) => filterChats(e.target.value));
        newChatBtn?.addEventListener('click', () => {
            window.location.href = '/ads/';
        });
    }

    async function loadInitialData() {
        await fetchCurrentUser();
        await loadChats();
        await initWebSocket();
    }

    async function fetchCurrentUser() {
        try {
            const response = await fetch('/api/user/profile/');
            if (!response.ok) return;
            state.currentUser = await response.json();
        } catch (error) {
            console.error('Не удалось получить данные пользователя', error);
        }
    }

    async function loadChats() {
        const chatsList = document.getElementById('chatsList');
        showChatsSkeletons(chatsList);
        try {
            const response = await fetch('/api/chats/');
            if (!response.ok) {
                throw new Error('Ошибка загрузки чатов');
            }
            const data = await response.json();
            state.chats = data.chats || [];
            renderChats();
        } catch (error) {
            console.error(error);
            chatsList.innerHTML = `
                <div class="no-chats">
                    <p>Не удалось загрузить чаты. Попробуйте позже.</p>
                </div>
            `;
        }
    }

    function renderChats() {
        const chatsList = document.getElementById('chatsList');
        if (!state.chats.length) {
            chatsList.innerHTML = `
                <div class="no-chats">
                    <div style="font-size: 2rem; margin-bottom: 15px;">💬</div>
                    <p>У вас пока нет сообщений</p>
                </div>
            `;
            return;
        }

        chatsList.innerHTML = state.chats
            .map((chat) => {
                const isActive = state.selectedChatId === chat.chat_id;
                const lastMessage = chat.last_message || 'Начните диалог';
                const lastMessageTime = chat.last_message_time
                    ? formatRelativeTime(chat.last_message_time)
                    : '';
                const avatarSrc = chat.companion_avatar_id
                    ? `/api/avatars/${chat.companion_avatar_id}/`
                    : '/static/pictures/profile.png';

                return `
                    <div class="chat-item ${isActive ? 'active' : ''}" data-chat-id="${chat.chat_id}" data-chat-name="${chat.companion_name}">
                        <div class="chat-avatar">
                            <img src="${avatarSrc}" alt="${chat.companion_name}" onerror="this.src='/static/pictures/profile.png'">
                        </div>
                        <div class="chat-info">
                            <div class="chat-header">
                                <h4 class="chat-name">${chat.companion_name}</h4>
                                <span class="chat-time">${lastMessageTime}</span>
                            </div>
                            <div class="chat-preview">
                                <p class="chat-message">${escapeHTML(lastMessage)}</p>
                            </div>
                        </div>
                    </div>
                `;
            })
            .join('');

        chatsList.querySelectorAll('.chat-item').forEach((item) => {
            item.addEventListener('click', () => {
                const chatId = item.getAttribute('data-chat-id');
                openChat(chatId);
            });
        });
    }

    function showChatsSkeletons(container) {
        container.innerHTML = `
            <div class="chat-item skeleton"></div>
            <div class="chat-item skeleton"></div>
            <div class="chat-item skeleton"></div>
        `;
    }

    function openChat(chatId) {
        if (!chatId) return;
        state.selectedChatId = chatId;
        document.querySelectorAll('.chat-item').forEach((item) => {
            item.classList.toggle('active', item.getAttribute('data-chat-id') === chatId);
        });

        const chatData = state.chats.find((chat) => chat.chat_id === chatId);
        if (!chatData) return;

        document.getElementById('chatPlaceholder').style.display = 'none';
        const chatWindow = document.getElementById('chatWindow');
        chatWindow.style.display = 'flex';
        document.getElementById('chatUserAvatar').src = chatData.companion_avatar_id
            ? `/api/avatars/${chatData.companion_avatar_id}/`
            : '/static/pictures/profile.png';
        document.getElementById('chatUserName').textContent = chatData.companion_name;
        document.getElementById('chatUserStatus').textContent = 'в сети';

        if (state.messagesCache.has(chatId)) {
            renderMessages(chatId);
        } else {
            loadMessages(chatId);
        }
    }

    async function loadMessages(chatId) {
        showMessagesSkeletons();
        try {
            const response = await fetch(`/api/chats/${chatId}/messages/?offset=0&limit=50`);
            if (!response.ok) {
                throw new Error('Ошибка загрузки сообщений');
            }
            const data = await response.json();
            state.messagesCache.set(chatId, data.messages || []);
            state.messagesMeta.set(chatId, {
                nextOffset: data.next_offset || 0,
                hasMore: data.has_more,
            });
            renderMessages(chatId);
        } catch (error) {
            console.error(error);
            const chatMessages = document.getElementById('chatMessages');
            chatMessages.innerHTML = `
                <div class="no-messages">
                    <p>Не удалось загрузить сообщения</p>
                </div>
            `;
        }
    }

    function renderMessages(chatId) {
        const chatMessages = document.getElementById('chatMessages');
        const messages = state.messagesCache.get(chatId) || [];
        const currentUserId = state.currentUser?.user_id;

        if (!messages.length) {
            chatMessages.innerHTML = `
                <div class="no-messages">
                    <p>Сообщений пока нет</p>
                </div>
            `;
            return;
        }

        chatMessages.innerHTML = messages
            .map((message) => {
                const isOwn = currentUserId && message.sender_id === currentUserId;
                const time = message.created_at ? formatMessageTime(message.created_at) : '';
                return `
                    <div class="message ${isOwn ? 'own' : ''}">
                        ${
                            !isOwn
                                ? `<img src="/static/pictures/profile.png" alt="Аватар" class="message-avatar">`
                                : ''
                        }
                        <div class="message-content">
                            <p class="message-text">${escapeHTML(message.text)}</p>
                            <div class="message-time">${time}</div>
                        </div>
                    </div>
                `;
            })
            .join('');

        chatMessages.scrollTop = chatMessages.scrollHeight;
    }

    function showMessagesSkeletons() {
        const chatMessages = document.getElementById('chatMessages');
        chatMessages.innerHTML = `
            <div class="message skeleton"></div>
            <div class="message skeleton"></div>
            <div class="message skeleton"></div>
        `;
    }

    function filterChats(query) {
        const normalized = query.toLowerCase();
        document.querySelectorAll('.chat-item').forEach((item) => {
            const name = item.getAttribute('data-chat-name')?.toLowerCase() || '';
            const lastMessage = item.querySelector('.chat-message')?.textContent.toLowerCase() || '';
            item.style.display = name.includes(normalized) || lastMessage.includes(normalized) ? 'flex' : 'none';
        });
    }

    async function initWebSocket() {
        try {
            const tokenResponse = await fetch('/api/messenger/ws-token/');
            if (!tokenResponse.ok) {
                throw new Error('Не удалось получить токен для WebSocket');
            }
            const { token } = await tokenResponse.json();
            const protocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
            state.ws = new WebSocket(`${protocol}${window.location.host}/ws/chat/?token=${token}`);

            state.ws.onmessage = (event) => {
                try {
                    const message = JSON.parse(event.data);
                    handleIncomingMessage(message);
                } catch (error) {
                    console.error('Ошибка обработки сообщения', error);
                }
            };

            state.ws.onclose = () => {
                if (state.reconnectTimer) return;
                state.reconnectTimer = setTimeout(() => {
                    state.reconnectTimer = null;
                    initWebSocket();
                }, 2000);
            };

            state.ws.onerror = (error) => {
                console.error('WebSocket ошибка', error);
                state.ws.close();
            };
        } catch (error) {
            console.error(error);
        }
    }

    function handleIncomingMessage(message) {
        if (!message?.chat_id) return;

        const existingMessages = state.messagesCache.get(message.chat_id) || [];
        existingMessages.push(message);
        state.messagesCache.set(message.chat_id, existingMessages);

        const chat = state.chats.find((c) => c.chat_id === message.chat_id);
        if (chat) {
            chat.last_message = message.text;
            chat.last_message_time = message.created_at;
            state.chats = [
                chat,
                ...state.chats.filter((c) => c.chat_id !== message.chat_id),
            ];
            renderChats();
        }

        if (state.selectedChatId === message.chat_id) {
            renderMessages(message.chat_id);
        }
    }

    function sendMessage() {
        const input = document.getElementById('chatMessageInput');
        const text = input.value.trim();
        if (!text || !state.selectedChatId) {
            return;
        }
        if (!state.ws || state.ws.readyState !== WebSocket.OPEN) {
            alert('Соединение с сервером сообщений отсутствует. Попробуйте позже.');
            return;
        }
        state.ws.send(JSON.stringify({
            chat_id: state.selectedChatId,
            text: text,
        }));
        input.value = '';
    }

    function formatRelativeTime(value) {
        const date = new Date(value);
        if (Number.isNaN(date.getTime())) return '';
        const diffMs = Date.now() - date.getTime();
        const diffMinutes = Math.floor(diffMs / (1000 * 60));
        if (diffMinutes < 1) return 'только что';
        if (diffMinutes < 60) return `${diffMinutes} мин назад`;
        const diffHours = Math.floor(diffMinutes / 60);
        if (diffHours < 24) return `${diffHours} ч назад`;
        return date.toLocaleDateString('ru-RU');
    }

    function formatMessageTime(value) {
        const date = new Date(value);
        if (Number.isNaN(date.getTime())) return '';
        return date.toLocaleTimeString('ru-RU', {
            hour: '2-digit',
            minute: '2-digit',
        });
    }

    function escapeHTML(unsafe) {
        if (!unsafe) return '';
        const div = document.createElement('div');
        div.textContent = unsafe;
        return div.innerHTML;
    }

    return { init };
})();

document.addEventListener('DOMContentLoaded', MessagesApp.init);
