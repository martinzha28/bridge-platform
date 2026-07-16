// Pure derivations from a PlayerView. No React, no network — unit-tested
// in view.test.ts. The backend view is authoritative; nothing here
// invents game data, it only reshapes what the view already contains.

import type { PlayerView, Seat } from "./protocol";
import { SEATS } from "./protocol";

/** The backend omits or nulls empty slices (`currentTrick`, `hand`,
 *  `calls`); coerce them to arrays so the rest of the code needn't. */
export function normalizeView(raw: PlayerView): PlayerView {
  return {
    ...raw,
    hand: raw.hand ?? [],
    calls: raw.calls ?? [],
    currentTrick: raw.currentTrick ?? [],
  };
}

export const SEAT_LETTER: Record<Seat, "N" | "E" | "S" | "W"> = {
  North: "N",
  East: "E",
  South: "S",
  West: "W",
};

export function nextSeat(seat: Seat): Seat {
  return SEATS[(SEATS.indexOf(seat) + 1) % 4];
}

export function partnerSeat(seat: Seat): Seat {
  return SEATS[(SEATS.indexOf(seat) + 2) % 4];
}

const SUIT_SYMBOL: Record<string, string> = {
  S: "♠",
  H: "♥",
  D: "♦",
  C: "♣",
};

const RANK_LABEL: Record<string, string> = { T: "10" };

export interface CardFace {
  rank: string;
  suit: string;
  red: boolean;
}

/** cardFace("HT") -> { rank: "10", suit: "♥", red: true } */
export function cardFace(code: string): CardFace {
  const suit = code[0];
  const rank = code[1];
  return {
    rank: RANK_LABEL[rank] ?? rank,
    suit: SUIT_SYMBOL[suit] ?? suit,
    red: suit === "H" || suit === "D",
  };
}

const DEFAULT_SUIT_ORDER = ["S", "H", "D", "C"];

/** The contract's trump suit letter, or null for no-trump / no contract.
 *  Contract strings look like "6D by South" or "4HX by North". */
export function trumpSuit(contract: string | undefined): string | null {
  const s = contract?.[1];
  return s === "S" || s === "H" || s === "D" || s === "C" ? s : null;
}

/** Split a hand into suit groups (high→low within a suit is already the
 *  backend's order). Trumps, if any, come first — leftmost for both your
 *  hand and dummy. */
export function groupHandBySuit(cards: string[], trump: string | null): string[][] {
  const order = trump
    ? [trump, ...DEFAULT_SUIT_ORDER.filter((s) => s !== trump)]
    : DEFAULT_SUIT_ORDER;
  return order.map((s) => cards.filter((c) => c[0] === s)).filter((g) => g.length > 0);
}

/** Cards still in a seat's hand, derived from tricks won + the live trick. */
export function handCount(view: PlayerView, seat: Seat): number {
  const completed = view.tricksNS + view.tricksEW;
  const playedToLiveTrick = view.currentTrick.some((c) => c.seat === seat) ? 1 : 0;
  return 13 - completed - playedToLiveTrick;
}

/** The declarer is the dummy's partner; unknown until the auction ends. */
export function declarerSeat(view: PlayerView): Seat | null {
  return view.dummy ? partnerSeat(view.dummy) : null;
}

export function canBid(view: PlayerView): boolean {
  return view.phase === "Auction" && view.turn === view.seat;
}

/** True when this client is on turn to play — its own turn, or dummy's
 *  turn while this client is declarer. */
export function canPlay(view: PlayerView): boolean {
  if (view.phase !== "Play") return false;
  if (view.turn === view.seat) return true;
  return view.turn === view.dummy && view.seat === declarerSeat(view);
}

export function trickCardFor(view: PlayerView, seat: Seat): string | undefined {
  return view.currentTrick.find((c) => c.seat === seat)?.card;
}

/** One-line summary of a finished board, shown during the pause before
 *  the next deal. */
export function boardResultText(view: PlayerView): string {
  const r = view.result;
  if (!r) return "Board complete";
  if (r.passedOut) return "Passed out";
  const declarerTricks =
    r.declarer === "North" || r.declarer === "South" ? r.tricksNS : r.tricksEW;
  const score = `${r.score >= 0 ? "+" : ""}${r.score}`;
  return `${r.contract} · ${declarerTricks} tricks · ${score}`;
}

// --- Auction ladder (right-side panel) -------------------------------------

export const AUCTION_COLUMNS: Seat[] = ["West", "North", "East", "South"];

/** Lay calls out under W/N/E/S columns, padding leading cells so the
 *  dealer's first call sits in the dealer's column. */
export function auctionRows(calls: string[], dealer: Seat): string[][] {
  if (calls.length === 0) return [];
  const lead = AUCTION_COLUMNS.indexOf(dealer);
  const cells: string[] = Array(lead).fill("").concat(calls);
  while (cells.length % 4 !== 0) cells.push("");
  const rows: string[][] = [];
  for (let i = 0; i < cells.length; i += 4) rows.push(cells.slice(i, i + 4));
  return rows;
}

// --- Bidding box ----------------------------------------------------------

export const BID_LEVELS = [1, 2, 3, 4, 5, 6, 7] as const;
export const BID_STRAINS = ["C", "D", "H", "S", "NT"] as const;
export type BidStrain = (typeof BID_STRAINS)[number];

export function strainLabel(strain: BidStrain): string {
  return strain === "NT" ? "NT" : SUIT_SYMBOL[strain];
}

export function strainIsRed(strain: BidStrain): boolean {
  return strain === "H" || strain === "D";
}

export function hasCall(legalCalls: string[] | undefined, call: string): boolean {
  return !!legalCalls && legalCalls.includes(call);
}
