"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import type { PlayerView, ServerMessage } from "@/lib/protocol";
import { openTableSocket, type TableSocket } from "@/lib/ws";
import { normalizeView } from "@/lib/view";

export type TableStatus = "connecting" | "setup" | "live" | "closed" | "error";

// How long the finished board stays on screen before the next deal.
const NEXT_BOARD_DELAY_MS = 4000;

export interface UseTable {
  view: PlayerView | null;
  status: TableStatus;
  error: string | null;
  bid: (call: string) => void;
  playCard: (card: string) => void;
}

/**
 * Connects to the game socket and drives the fixed local setup:
 * create a table, sit South, fill N/E/W with bots, start board one.
 * From then on the returned `view` is whatever the backend last sent.
 */
export function useTable(): UseTable {
  const [view, setView] = useState<PlayerView | null>(null);
  const [status, setStatus] = useState<TableStatus>("connecting");
  const [error, setError] = useState<string | null>(null);
  const socketRef = useRef<TableSocket | null>(null);
  const nextBoardTimer = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    const socket = openTableSocket({
      onOpen() {
        setStatus("setup");
        socket.send({ type: "create_table" });
      },
      onMessage(msg: ServerMessage) {
        switch (msg.type) {
          case "table_created":
            socket.send({ type: "sit", direction: "S" });
            socket.send({ type: "sit_bot", direction: "N", difficulty: 1 });
            socket.send({ type: "sit_bot", direction: "E", difficulty: 1 });
            socket.send({ type: "sit_bot", direction: "W", difficulty: 1 });
            socket.send({ type: "start" });
            break;
          case "game_state": {
            const next = normalizeView(msg.payload as PlayerView);
            setView(next);
            setStatus("live");
            // When a board finishes, deal the next one after a short
            // pause so the result is readable.
            if (next.phase === "Complete" && nextBoardTimer.current === null) {
              nextBoardTimer.current = setTimeout(() => {
                nextBoardTimer.current = null;
                socket.send({ type: "start" });
              }, NEXT_BOARD_DELAY_MS);
            }
            break;
          }
          case "error":
            setError(msg.error ?? "unknown error");
            break;
        }
      },
      onClose() {
        setStatus((s) => (s === "error" ? s : "closed"));
      },
      onError() {
        setError("connection failed");
        setStatus("error");
      },
    });
    socketRef.current = socket;

    return () => {
      if (nextBoardTimer.current !== null) {
        clearTimeout(nextBoardTimer.current);
        nextBoardTimer.current = null;
      }
      socket.close();
      socketRef.current = null;
    };
  }, []);

  const bid = useCallback((call: string) => {
    socketRef.current?.send({ type: "bid", call });
  }, []);

  const playCard = useCallback((card: string) => {
    socketRef.current?.send({ type: "play_card", card });
  }, []);

  return { view, status, error, bid, playCard };
}
