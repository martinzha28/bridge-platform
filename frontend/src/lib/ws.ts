import type { ClientMessage, ServerMessage } from "./protocol";

const WS_URL = process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:8080/ws";

export interface TableSocket {
  send(msg: ClientMessage): void;
  close(): void;
}

interface Handlers {
  onOpen?: () => void;
  onMessage: (msg: ServerMessage) => void;
  onClose?: () => void;
  onError?: () => void;
}

/** Opens the game WebSocket and adapts it to typed JSON messages. */
export function openTableSocket(handlers: Handlers): TableSocket {
  const ws = new WebSocket(WS_URL);

  ws.addEventListener("open", () => handlers.onOpen?.());
  ws.addEventListener("close", () => handlers.onClose?.());
  ws.addEventListener("error", () => handlers.onError?.());
  ws.addEventListener("message", (ev) => {
    try {
      handlers.onMessage(JSON.parse(ev.data) as ServerMessage);
    } catch {
      // ignore frames that aren't JSON
    }
  });

  return {
    send(msg) {
      if (ws.readyState === WebSocket.OPEN) ws.send(JSON.stringify(msg));
    },
    close() {
      ws.close();
    },
  };
}
