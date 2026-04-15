import { appConfig } from '@/config';

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
    if (
      this.socket &&
      (this.socket.readyState === WebSocket.CONNECTING ||
        this.socket.readyState === WebSocket.OPEN)
    ) {
      return;
    }

    const socket = new WebSocket(url);

    this.socket = socket;
    this.notifyStatus(socket);

    socket.addEventListener('open', () => {
      if (this.socket !== socket) {
        return;
      }

      this.notifyStatus(socket);
    });

    socket.addEventListener('close', () => {
      if (this.socket !== socket) {
        return;
      }

      this.socket = null;
      this.notifyStatus(socket);
    });

    socket.addEventListener('error', () => {
      if (this.socket !== socket) {
        return;
      }

      this.statusListeners.forEach((listener) => listener('error'));
    });

    socket.addEventListener('message', (event) => {
      if (this.socket !== socket) {
        return;
      }

      const payload = this.parseMessage(event.data);

      this.messageListeners.forEach((listener) => listener(payload));
    });
  }

  public disconnect() {
    if (!this.socket) {
      return;
    }

    const socket = this.socket;

    this.socket = null;

    if (socket.readyState !== WebSocket.CLOSED) {
      socket.close();
    }

    this.notifyStatus(socket);
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

  private notifyStatus(socket: WebSocket | null = this.socket) {
    const status = getStatusFromSocket(socket);

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
