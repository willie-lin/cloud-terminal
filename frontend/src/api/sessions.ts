/**
 * WebSocket session management for terminal connections.
 * Connects to the backend WebSocket endpoint at /api/terminal.
 */

const WS_BASE = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}`;

export interface SessionEvents {
  onData: (data: ArrayBuffer) => void;
  onError: (error: Event) => void;
  onClose: (code: number, reason: string) => void;
}

export interface SessionController {
  send: (data: ArrayBuffer | string) => void;
  resize: (cols: number, rows: number) => void;
  close: () => void;
}

/**
 * Open a WebSocket terminal session.
 *
 * @param urn - Resource URN to connect to
 * @param token - STS token for authentication
 * @param events - Event callbacks
 * @returns SessionController for sending data and managing the session
 */
export function connectTerminalSession(
  urn: string,
  token: string,
  events: SessionEvents,
): SessionController {
  const url = `${WS_BASE}/api/terminal?urn=${encodeURIComponent(urn)}&token=${encodeURIComponent(token)}`;
  const ws = new WebSocket(url);

  ws.binaryType = 'arraybuffer';

  ws.onmessage = (event) => {
    if (event.data instanceof ArrayBuffer) {
      events.onData(event.data);
    }
  };

  ws.onerror = (event) => {
    events.onError(event);
  };

  ws.onclose = (event) => {
    events.onClose(event.code, event.reason);
  };

  return {
    send: (data) => {
      if (ws.readyState === WebSocket.OPEN) {
        if (typeof data === 'string') {
          ws.send(data);
        } else {
          ws.send(data);
        }
      }
    },
    resize: (cols: number, rows: number) => {
      if (ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ type: 'resize', data: { cols, rows } }));
      }
    },
    close: () => {
      ws.close();
    },
  };
}

export default connectTerminalSession;
