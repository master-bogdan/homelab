import { appConfig } from '@/shared/config/env';

export type WebSocketMessageHandler = (payload: unknown) => void;
export type WebSocketStatus =
  | 'closed'
  | 'closing'
  | 'connecting'
  | 'error'
  | 'open'
  | 'uninitialized';
export type WebSocketStatusHandler = (status: WebSocketStatus) => void;

const getStatusFromSocket = (socket: WebSocket | null): WebSocketStatus => {
  if (!socket) {
    return 'uninitialized';
  }

  switch (socket.readyState) {
    case WebSocket.CONNECTING:
      return 'connecting';
    case WebSocket.OPEN:
      return 'open';
    case WebSocket.CLOSING:
      return 'closing';
    case WebSocket.CLOSED:
      return 'closed';
    default:
      return 'closed';
  }
};

export class WebSocketClient {
  private messageListeners = new Set<WebSocketMessageHandler>();
  private socket: WebSocket | null = null;
  private statusListeners = new Set<WebSocketStatusHandler>();

  public connect(url = appConfig.wsBaseUrl) {
    if (this.socket && this.socket.readyState <= WebSocket.OPEN) {
      return;
    }

    this.socket = new WebSocket(url);
    this.notifyStatus();

    this.socket.addEventListener('open', () => {
      this.notifyStatus();
    });

    this.socket.addEventListener('close', () => {
      this.notifyStatus();
    });

    this.socket.addEventListener('error', () => {
      this.statusListeners.forEach((listener) => listener('error'));
    });

    this.socket.addEventListener('message', (event) => {
      const payload = this.parseMessage(event.data);

      this.messageListeners.forEach((listener) => listener(payload));
    });
  }

  public disconnect() {
    this.socket?.close();
  }

  public getStatus(): WebSocketStatus {
    return getStatusFromSocket(this.socket);
  }

  public send(payload: Record<string, unknown>) {
    if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
      throw new Error('WebSocket connection is not open.');
    }

    this.socket.send(JSON.stringify(payload));
  }

  public subscribeToMessages(listener: WebSocketMessageHandler) {
    this.messageListeners.add(listener);

    return () => {
      this.messageListeners.delete(listener);
    };
  }

  public subscribeToStatus(listener: WebSocketStatusHandler) {
    this.statusListeners.add(listener);
    listener(this.getStatus());

    return () => {
      this.statusListeners.delete(listener);
    };
  }

  private notifyStatus() {
    const status = this.getStatus();

    this.statusListeners.forEach((listener) => listener(status));
  }

  private parseMessage(payload: string) {
    try {
      return JSON.parse(payload) as unknown;
    } catch {
      return payload;
    }
  }
}

export const appWebSocketClient = new WebSocketClient();
