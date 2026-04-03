import { WebSocketClient, type WebSocketStatus } from '../webSocketClient';

type WebSocketEventName = 'close' | 'error' | 'message' | 'open';
type WebSocketEventHandler = (event: Event | MessageEvent<string>) => void;

class MockWebSocket {
  public static readonly CONNECTING = 0;
  public static readonly OPEN = 1;
  public static readonly CLOSING = 2;
  public static readonly CLOSED = 3;

  public static instances: MockWebSocket[] = [];

  public readonly url: string;
  public readyState = MockWebSocket.CONNECTING;
  public send = vi.fn();

  private listeners: Record<WebSocketEventName, WebSocketEventHandler[]> = {
    close: [],
    error: [],
    message: [],
    open: []
  };

  public constructor(url: string) {
    this.url = url;
    MockWebSocket.instances.push(this);
  }

  public addEventListener(type: WebSocketEventName, listener: WebSocketEventHandler) {
    this.listeners[type].push(listener);
  }

  public close() {
    this.readyState = MockWebSocket.CLOSING;
  }

  public emitClose() {
    this.readyState = MockWebSocket.CLOSED;
    this.listeners.close.forEach((listener) => listener(new Event('close')));
  }

  public emitMessage(data: string) {
    this.listeners.message.forEach((listener) =>
      listener(new MessageEvent('message', { data }))
    );
  }

  public emitOpen() {
    this.readyState = MockWebSocket.OPEN;
    this.listeners.open.forEach((listener) => listener(new Event('open')));
  }
}

describe('WebSocketClient', () => {
  beforeEach(() => {
    MockWebSocket.instances = [];
    vi.stubGlobal('WebSocket', MockWebSocket as unknown as typeof WebSocket);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('ignores status and message events from a stale socket after reconnect', () => {
    const client = new WebSocketClient();
    const messages: unknown[] = [];
    const statuses: WebSocketStatus[] = [];

    client.subscribeToMessages((payload) => {
      messages.push(payload);
    });
    client.subscribeToStatus((status) => {
      statuses.push(status);
    });

    client.connect('ws://first');

    const firstSocket = MockWebSocket.instances[0];

    firstSocket.emitOpen();
    client.disconnect();
    client.connect('ws://second');

    const secondSocket = MockWebSocket.instances[1];

    secondSocket.emitOpen();

    const statusCountBeforeStaleClose = statuses.length;

    firstSocket.emitMessage(JSON.stringify({ socket: 'stale' }));
    firstSocket.emitClose();
    secondSocket.emitMessage(JSON.stringify({ socket: 'current' }));

    expect(statuses).toHaveLength(statusCountBeforeStaleClose);
    expect(statuses.at(-1)).toBe('open');
    expect(messages).toEqual([{ socket: 'current' }]);
  });
});
