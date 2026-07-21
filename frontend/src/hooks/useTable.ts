"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import type {
  ChatMessage,
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
// Cap on chat lines kept client-side.
const CHAT_KEEP = 200;

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
  /** Where this client is seated, or null if standing in the lobby. */
  mySeat: Seat | null;
  messages: ChatMessage[];
  bid: (call: string) => void;
  playCard: (card: string) => void;
  sitAt: (dir: Seat) => void;
  stand: () => void;
  takeSeat: (dir: Seat) => void;
  addBot: (dir: Seat) => void;
  removeBot: (dir: Seat) => void;
  moveTo: (dir: Seat) => void;
  setName: (name: string) => void;
  sendChat: (text: string) => void;
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
  const [mySeat, setMySeat] = useState<Seat | null>(null);
  const [messages, setMessages] = useState<ChatMessage[]>([]);

  const socketRef = useRef<TableSocket | null>(null);
  const nextBoardTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const holdTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const pendingView = useRef<PlayerView | null>(null);
  // seat requests the host has already sent this session (keyed seat+role),
  // so a burst of table_state messages doesn't re-send them
  const seatRequests = useRef<Set<string>>(new Set());
  const startRequested = useRef(false);
  const mySeatRef = useRef<Seat | null>(null);
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

    // The host lays out its seat plan on every lobby update: it fills any
    // seat the plan wants that is still empty, once each. Idempotent, so a
    // socket remount self-heals and a late joiner still sees the host +
    // bots. It never auto-starts — the host clicks Start — so a friend
    // always has a lobby to join from, even on an all-bot table.
    function reconcileHostSeats(state: TableState) {
      if (!isHost || state.started) return;
      const plan = seatPlan(configRef.current);
      const want = new Map<Seat, "you" | "bot">();
      if (plan.you) want.set(plan.you, "you");
      for (const b of plan.bots) want.set(b, "bot");

      for (const [seat, role] of want) {
        if (state.seats[seat] !== "") continue;
        const key = `${seat}:${role}`;
        if (seatRequests.current.has(key)) continue;
        seatRequests.current.add(key);
        if (role === "you") {
          socket.send({ type: "sit", direction: SEAT_LETTER[seat] });
        } else {
          socket.send({ type: "sit_bot", direction: SEAT_LETTER[seat], difficulty: 1 });
        }
      }
    }

    // Guards against a torn-down socket (e.g. StrictMode's throwaway first
    // mount) firing a late error/close that poisons the live connection's
    // state.
    let disposed = false;

    const socket = openTableSocket({
      onOpen() {
        if (disposed) return;
        setError(null);
        setStatus("connecting");
        socket.send({ type: "join_table", tableID: tableId });
      },
      onMessage(msg: ServerMessage) {
        if (disposed) return;
        switch (msg.type) {
          case "table_joined":
            break;
          case "seated": {
            // sit_bot also acks with "seated" (bot:"true") — that's not us
            const seat = msg.payload as { direction: Seat; bot?: string };
            if (seat.bot === "true") break;
            mySeatRef.current = seat.direction;
            setMySeat(seat.direction);
            break;
          }
          case "stood":
            mySeatRef.current = null;
            setMySeat(null);
            break;
          case "chat_history":
            setMessages((msg.payload as ChatMessage[]).slice(-CHAT_KEEP));
            break;
          case "chat_message":
            setMessages((m) =>
              [...m, msg.payload as ChatMessage].slice(-CHAT_KEEP),
            );
            break;
          case "table_state": {
            const state = msg.payload as TableState;
            setTableState(state);
            if (!state.started) setStatus("lobby");
            reconcileHostSeats(state);
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
        if (disposed) return;
        setStatus((s) => (s === "error" ? s : "closed"));
      },
      onError() {
        if (disposed) return;
        setError("connection failed");
        setStatus("error");
      },
    });
    socketRef.current = socket;

    return () => {
      disposed = true;
      for (const timer of [nextBoardTimer, holdTimer]) {
        if (timer.current !== null) {
          clearTimeout(timer.current);
          timer.current = null;
        }
      }
      pendingView.current = null;
      seatRequests.current = new Set();
      startRequested.current = false;
      mySeatRef.current = null;
      setMySeat(null);
      setMessages([]);
      setPaused(false);
      setLastCompletedTrick([]);
      setError(null);
      setJoinError(null);
      setView(null);
      setTableState(null);
      setStatus("connecting");
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
    if (mySeatRef.current && mySeatRef.current !== dir) {
      socketRef.current?.send({ type: "stand" });
    }
    socketRef.current?.send({ type: "sit", direction: SEAT_LETTER[dir] });
  }, []);
  const stand = useCallback(() => {
    socketRef.current?.send({ type: "stand" });
  }, []);
  // claim a seat a bot currently holds (lobby only — the backend
  // refuses remove_bot once a game is running)
  const takeSeat = useCallback((dir: Seat) => {
    if (mySeatRef.current && mySeatRef.current !== dir) {
      socketRef.current?.send({ type: "stand" });
    }
    socketRef.current?.send({ type: "remove_bot", direction: SEAT_LETTER[dir] });
    socketRef.current?.send({ type: "sit", direction: SEAT_LETTER[dir] });
  }, []);
  const addBot = useCallback((dir: Seat) => {
    socketRef.current?.send({ type: "sit_bot", direction: SEAT_LETTER[dir], difficulty: 1 });
  }, []);
  const removeBot = useCallback((dir: Seat) => {
    socketRef.current?.send({ type: "remove_bot", direction: SEAT_LETTER[dir] });
  }, []);
  // Move to another seat, leaving a bot in the seat vacated so the table
  // stays full. Lobby-only (backend refuses remove_bot mid-game).
  const moveTo = useCallback((dir: Seat) => {
    const s = socketRef.current;
    const from = mySeatRef.current;
    if (!s || from === dir) return;
    if (from) {
      s.send({ type: "stand" });
      s.send({ type: "sit_bot", direction: SEAT_LETTER[from], difficulty: 1 });
    }
    s.send({ type: "remove_bot", direction: SEAT_LETTER[dir] });
    s.send({ type: "sit", direction: SEAT_LETTER[dir] });
  }, []);
  const setName = useCallback((name: string) => {
    socketRef.current?.send({ type: "set_name", name });
  }, []);
  const sendChat = useCallback((text: string) => {
    const trimmed = text.trim();
    if (trimmed) socketRef.current?.send({ type: "chat", text: trimmed });
  }, []);
  const startGame = useCallback(() => {
    if (startRequested.current) return;
    startRequested.current = true;
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
    mySeat,
    messages,
    bid,
    playCard,
    sitAt,
    stand,
    takeSeat,
    addBot,
    removeBot,
    moveTo,
    setName,
    sendChat,
    startGame,
  };
}
