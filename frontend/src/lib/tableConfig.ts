import type { Seat } from "./protocol";
import { SEATS } from "./protocol";

export type SeatRole = "you" | "bot" | "open";
export type Mode = "casual" | "competitive";
export type Visibility = "public" | "private";
export type Scoring = "imps" | "matchpoints";

export interface TableConfig {
  name: string;
  description: string;
  mode: Mode;
  visibility: Visibility;
  scoring: Scoring;
  boards: number | null; // null = unlimited
  seats: Record<Seat, SeatRole>;
}

export const DEFAULT_CONFIG: TableConfig = {
  name: "Practice table",
  description: "",
  mode: "casual",
  visibility: "private",
  scoring: "imps",
  boards: 24,
  seats: { South: "you", North: "bot", East: "bot", West: "bot" },
};

export const MAX_BOARDS = 64;

/** Parse the boards field: "" → unlimited (null); otherwise 1..MAX_BOARDS. */
export function parseBoards(raw: string): number | null {
  const t = raw.trim();
  if (t === "") return null;
  const n = Math.floor(Number(t));
  if (!Number.isFinite(n) || n < 1) return null;
  return Math.min(n, MAX_BOARDS);
}

export function boardsInput(boards: number | null): string {
  return boards == null ? "" : String(boards);
}

export interface SeatPlan {
  you: Seat | null;
  bots: Seat[];
  open: Seat[];
}

export function seatPlan(config: TableConfig): SeatPlan {
  return {
    you: SEATS.find((s) => config.seats[s] === "you") ?? null,
    bots: SEATS.filter((s) => config.seats[s] === "bot"),
    open: SEATS.filter((s) => config.seats[s] === "open"),
  };
}

export function inviteUrl(origin: string, tableId: string): string {
  return `${origin}/table/${tableId}`;
}

// --- per-table config passed from /create to /table/[id] via sessionStorage ---

const key = (id: string) => `table-config:${id}`;

export function storeConfig(id: string, config: TableConfig): void {
  try {
    sessionStorage.setItem(key(id), JSON.stringify(config));
  } catch {
    // sessionStorage unavailable — the table page falls back to defaults
  }
}

export function loadConfig(id: string): TableConfig | null {
  try {
    const raw = sessionStorage.getItem(key(id));
    return raw ? (JSON.parse(raw) as TableConfig) : null;
  } catch {
    return null;
  }
}
