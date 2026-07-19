"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import type { PlayedCard, PlayerView, ServerMessage } from "@/lib/protocol";
import { openTableSocket, type TableSocket } from "@/lib/ws";
import {
  boardSummary,
  normalizeView,
  trickJustFinished,
  type BoardResult,
} from "@/lib/view";

export type TableStatus = "connecting" | "setup" | "live" | "closed" | "error";

// How long the finished board stays on screen before the next deal.
const NEXT_BOARD_DELAY_MS = 4000;
// How long a completed trick sits on the table before play catches up.
const TRICK_PAUSE_MS = 1500;

export interface UseTable {
  view: PlayerView | null;
  status: TableStatus;
  error: string | null;
  /** Finished boards at this table, oldest first. */
  history: BoardResult[];
  /** True while a finished trick is held on the table (auto-pause). */
  paused: boolean;
  /** The most recently completed trick this board, or empty. */
  lastCompletedTrick: PlayedCard[];
  bid: (call: string) => void;
  playCard: (card: string) => void;
}

/**
 * Connects to the game socket and drives the fixed local setup:
 * create a table, sit South, fill N/E/W with bots, start board one.
 * From then on the returned `view` is whatever the backend last sent,
 * except that a finished trick is briefly held on screen before play
 * catches up.
 */
export function useTable(): UseTable {
  const [view, setView] = useState<PlayerView | null>(null);
  const [status, setStatus] = useState<TableStatus>("connecting");
  const [error, setError] = useState<string | null>(null);
  const [history, setHistory] = useState<BoardResult[]>([]);
  const [paused, setPaused] = useState(false);
  const [lastCompletedTrick, setLastCompletedTrick] = useState<PlayedCard[]>([]);

  const socketRef = useRef<TableSocket | null>(null);
  const nextBoardTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const holdTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const pendingView = useRef<PlayerView | null>(null);

  useEffect(() => {
    // endHold releases the frozen screen and catches up to whatever
    // state arrived while it was held.
    function endHold() {
      holdTimer.current = null;
      setPaused(false);
      const buffered = pendingView.current;
      pendingView.current = null;
      if (buffered) ingest(buffered);
    }

    function ingest(next: PlayerView) {
      setView(next);
      setStatus("live");

      if (next.phase === "Auction") {
        setLastCompletedTrick([]);
      } else if (next.lastTrick && next.lastTrick.length > 0) {
        setLastCompletedTrick(next.lastTrick);
      }

      if (next.phase === "Complete" && nextBoardTimer.current === null) {
        setHistory((h) => [...h, boardSummary(next)]);
        nextBoardTimer.current = setTimeout(() => {
          nextBoardTimer.current = null;
          socket.send({ type: "start" });
        }, NEXT_BOARD_DELAY_MS);
        return;
      }

      if (trickJustFinished(next)) {
        setPaused(true);
        holdTimer.current = setTimeout(endHold, TRICK_PAUSE_MS);
      }
    }

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
            if (holdTimer.current !== null) {
              pendingView.current = next; // caught up when the hold ends
            } else {
              ingest(next);
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
      for (const timer of [nextBoardTimer, holdTimer]) {
        if (timer.current !== null) {
          clearTimeout(timer.current);
          timer.current = null;
        }
      }
      pendingView.current = null;
      setPaused(false);
      setLastCompletedTrick([]);
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

  return { view, status, error, history, paused, lastCompletedTrick, bid, playCard };
}
