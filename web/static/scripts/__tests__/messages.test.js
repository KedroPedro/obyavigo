const fs = require('fs');
const path = require('path');
const vm = require('vm');
const { JSDOM } = require('jsdom');

const scriptPath = path.resolve(__dirname, '../messages.js');
const scriptContent = fs.readFileSync(scriptPath, 'utf-8');

const CHAT_ID = '123e4567-e89b-12d3-a456-426614174000';
const SAMPLE_CHAT = {
    chat_id: CHAT_ID,
    companion_name: 'Алиса',
    companion_avatar_id: null,
    last_message: 'Привет!',
    last_message_time: '2025-01-01T12:00:00Z',
};

const SAMPLE_MESSAGE = {
    message_id: 'msg-1',
    chat_id: CHAT_ID,
    sender_id: 99,
    text: 'Первое сообщение',
    created_at: '2025-01-01T12:01:00Z',
};

const flushAsync = () => new Promise((resolve) => setTimeout(resolve, 0));

function createFetchMock() {
    return jest.fn((url, options = {}) => {
        if (url === '/api/user/profile/') {
            return Promise.resolve({
                ok: true,
                json: () => Promise.resolve({ user_id: 99, name: 'User' }),
            });
        }
        if (url === '/api/chats/' && (!options.method || options.method === 'GET')) {
            return Promise.resolve({
                ok: true,
                json: () => Promise.resolve({ chats: [SAMPLE_CHAT] }),
            });
        }
        if (url.startsWith(`/api/chats/${CHAT_ID}/messages/`)) {
            return Promise.resolve({
                ok: true,
                json: () => Promise.resolve({
                    messages: [SAMPLE_MESSAGE],
                    next_offset: 0,
                    has_more: false,
                }),
            });
        }
        if (url === '/api/messenger/ws-token/') {
            return Promise.resolve({
                ok: true,
                json: () => Promise.resolve({ token: 'ws-token' }),
            });
        }
        return Promise.reject(new Error(`Unhandled fetch url: ${url}`));
    });
}

function setupTestContext() {
    const dom = new JSDOM(
        `<!DOCTYPE html>
        <body>
            <div id="chatsList"></div>
            <textarea id="chatMessageInput"></textarea>
            <button id="chatSendBtn">➤</button>
            <input id="chatsSearch" />
            <button id="newChatBtn"></button>
            <button id="openChatsDrawer"></button>
            <button id="closeChatsDrawer"></button>
            <div id="mobileChatsOverlay"></div>
            <div id="chatPlaceholder" style="display:block;"></div>
            <div id="chatWindow" style="display:none;"></div>
            <img id="chatUserAvatar" />
            <div id="chatUserName"></div>
            <div id="chatUserStatus"></div>
            <div id="chatMessages"></div>
            <div class="chats-sidebar"></div>
        </body>`,
        { url: 'https://example.com/messages' }
    );

    const { window } = dom;
    const { document } = window;
    Object.defineProperty(window, 'innerWidth', { value: 1024, writable: true });

    const fetchMock = createFetchMock();
    const alertMock = jest.fn();

    class MockWebSocket {
        static CONNECTING = 0;
        static OPEN = 1;
        static CLOSING = 2;
        static CLOSED = 3;
        static instances = [];

        constructor(url) {
            this.url = url;
            this.readyState = MockWebSocket.CONNECTING;
            this.send = jest.fn();
            this.close = jest.fn(() => {
                this.readyState = MockWebSocket.CLOSED;
            });
            MockWebSocket.instances.push(this);
            setTimeout(() => {
                this.readyState = MockWebSocket.OPEN;
                this.onopen?.();
            }, 0);
        }
    }

    const consoleStub = {
        log: jest.fn(),
        error: jest.fn(),
        warn: jest.fn(),
    };

    const baseContext = {
        window,
        document,
        console: consoleStub,
        localStorage: window.localStorage,
        navigator: window.navigator,
        location: window.location,
        fetch: fetchMock,
        alert: alertMock,
        WebSocket: MockWebSocket,
        setTimeout,
        clearTimeout,
        setInterval,
        clearInterval,
        performance,
    };

    baseContext.global = baseContext;
    Object.assign(window, {
        fetch: fetchMock,
        alert: alertMock,
        WebSocket: MockWebSocket,
        console: consoleStub,
    });

    const context = vm.createContext(baseContext);
    vm.runInContext(scriptContent, context);
    const app = vm.runInContext('MessagesApp', context);

    return { context, window, document, fetchMock, alertMock, MockWebSocket, app };
}

describe('MessagesApp', () => {
    afterEach(() => {
        jest.clearAllMocks();
    });

    test('renders chats and loads messages after selecting a chat', async () => {
        const { app, window, document } = setupTestContext();
        document.dispatchEvent(new window.Event('DOMContentLoaded'));

        await flushAsync();
        await flushAsync();

        const chatItems = document.querySelectorAll('.chat-item');
        expect(chatItems).toHaveLength(1);
        expect(chatItems[0].querySelector('.chat-name').textContent).toBe('Алиса');

        chatItems[0].click();
        await flushAsync();
        await flushAsync();

        const chatMessages = document.getElementById('chatMessages');
        expect(chatMessages.textContent).toContain('Первое сообщение');
        expect(chatMessages.querySelectorAll('.message').length).toBe(1);
        expect(document.getElementById('chatWindow').style.display).toBe('flex');

        app.destroy();
        window.close();
    });

    test('sends message via WebSocket when chat is selected', async () => {
        const { app, window, document, alertMock, MockWebSocket } = setupTestContext();
        document.dispatchEvent(new window.Event('DOMContentLoaded'));

        await flushAsync();
        await flushAsync();

        document.querySelector('.chat-item').click();
        await flushAsync();
        await flushAsync();

        const wsInstance = MockWebSocket.instances.at(-1);
        expect(wsInstance).toBeDefined();

        await flushAsync();

        const input = document.getElementById('chatMessageInput');
        input.value = '   Новое сообщение   ';
        document.getElementById('chatSendBtn').click();

        expect(wsInstance.send).toHaveBeenCalledTimes(1);
        const payload = JSON.parse(wsInstance.send.mock.calls[0][0]);
        expect(payload).toEqual({
            chat_id: CHAT_ID,
            text: 'Новое сообщение',
        });
        expect(input.value).toBe('');
        expect(alertMock).not.toHaveBeenCalled();

        app.destroy();
        window.close();
    });
});

