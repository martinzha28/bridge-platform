"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import type {
  PlayedCard,
  PlayerView,
  Seat,
  ServerMessage,
  TableState,
} from "@/lib/protocol";
import { openTableSocket, type TableSocket } from "@/lib/ws";
import {
  boardSummary,
  normalizeView,
  trickJustFinished,
  type BoardResult,
} from "@/lib/view";
import { seatPlan, type TableConfig } from "@/lib/tableConfig";

export type TableStatus =
  | "connecting"
  | "lobby"
  | "live"
  | "closed"
  | "error";

// How long the finished board stays on screen before the next deal.
const NEXT_BOARD_DELAY_MS = 4000;
// How long a completed trick sits on the table before play catches up.
const TRICK_PAUSE_MS = 1500;

const SEAT_LETTER: Record<Seat, "N" | "E" | "S" | "W"> = {
  North: "N",
  East: "E",
  South: "S",
  West: "W",
};

export interface UseTableOptions {
  tableId: string;
  config: TableConfig;
  /** The creator — applies the seat plan and drives board advancement. */
  isHost: boolean;
}

export interface UseTable {
  view: PlayerView | null;
  tableState: TableState | null;
  status: TableStatus;
  error: string | null;
  /** Set if the table id doesn't exist. */
  joinError: string | null;
  history: BoardResult[];
  paused: boolean;
  lastCompletedTrick: PlayedCard[];
  bid: (call: string) => void;
  playCard: (card: string) => void;
  sitAt: (dir: Seat) => void;
  startGame: () => void;
}

export function useTable({ tableId, config, isHost }: UseTableOptions): UseTable {
  const [view, setView] = useState<PlayerView | null>(null);
  const [tableState, setTableState] = useState<TableState | null>(null);
  const [status, setStatus] = useState<TableStatus>("connecting");
  const [error, setError] = useState<string | null>(null);
  const [joinError, setJoinError] = useState<string | null>(null);
  const [history, setHistory] = useState<BoardResult[]>([]);
  const [paused, setPaused] = useState(false);
  const [lastCompletedTrick, setLastCompletedTrick] = useState<PlayedCard[]>([]);

  const socketRef = useRef<TableSocket | null>(null);
  const nextBoardTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const holdTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const pendingView = useRef<PlayerView | null>(null);
  const planApplied = useRef(false);
  const configRef = useRef(config);
  configRef.current = config;

  useEffect(() => {
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
        setHistory((h) => {
          const boards = h.length + 1;
          const limit = configRef.current.boards;
          if (isHost && (limit == null || boards < limit)) {
            nextBoardTimer.current = setTimeout(() => {
              nextBoardTimer.current = null;
              socket.send({ type: "start" });
            }, NEXT_BOARD_DELAY_MS);
          }
          return [...h, boardSummary(next)];
        });
        return;
      }

      if (trickJustFinished(next)) {
        setPaused(true);
        holdTimer.current = setTimeout(endHold, TRICK_PAUSE_MS);
      }
    }

    // The host seats itself + its bots the first time it sees the lobby.
    function applyPlanOnce(state: TableState) {
      if (planApplied.current || !isHost || state.started) return;
      planApplied.current = true;
      const plan = seatPlan(configRef.current);
      if (plan.you) socket.send({ type: "sit", direction: SEAT_LETTER[plan.you] });
      for (const b of plan.bots) {
        socket.send({ type: "sit_bot", direction: SEAT_LETTER[b], difficulty: 1 });
      }
      if (plan.open.length === 0) socket.send({ type: "start" });
    }

    const socket = openTableSocket({
      onOpen() {
        setStatus("connecting");
        socket.send({ type: "join_table", tableID: tableId });
      },
      onMessage(msg: ServerMessage) {
        switch (msg.type) {
          case "table_joined":
            break;
          case "table_state": {
            const state = msg.payload as TableState;
            setTableState(state);
            if (!state.started) setStatus("lobby");
            applyPlanOnce(state);
            break;
          }
          case "game_state": {
            const next = normalizeView(msg.payload as PlayerView);
            if (holdTimer.current !== null) {
              pendingView.current = next;
            } else {
              ingest(next);
            }
            break;
          }
          case "error":
            if (msg.error?.startsWith("table not found")) {
              setJoinError("This table no longer exists.");
              setStatus("error");
            } else {
              setError(msg.error ?? "unknown error");
            }
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
      planApplied.current = false;
      setPaused(false);
      setLastCompletedTrick([]);
      socket.close();
      socketRef.current = null;
    };
  }, [tableId, isHost]);

  const bid = useCallback((call: string) => {
    socketRef.current?.send({ type: "bid", call });
  }, []);
  const playCard = useCallback((card: string) => {
    socketRef.current?.send({ type: "play_card", card });
  }, []);
  const sitAt = useCallback((dir: Seat) => {
    socketRef.current?.send({ type: "sit", direction: SEAT_LETTER[dir] });
  }, []);
  const startGame = useCallback(() => {
    socketRef.current?.send({ type: "start" });
  }, []);

  return {
    view,
    tableState,
    status,
    error,
    joinError,
    history,
    paused,
    lastCompletedTrick,
    bid,
    playCard,
    sitAt,
    startGame,
  };
}
