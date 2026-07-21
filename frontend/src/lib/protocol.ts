// Wire types mirroring the backend's ws + game packages.
// Card codes are suit letter + rank letter, e.g. "SA", "HT", "C2".
// Call codes are "P" | "X" | "XX" | "1C".."7NT".

export type Seat = "North" | "East" | "South" | "West";
export const SEATS: Seat[] = ["North", "East", "South", "West"];

export type Phase = "Deal" | "Auction" | "Play" | "Complete";

export interface PlayedCard {
  seat: Seat;
  card: string;
}

export interface ResultView {
  contract?: string;
  declarer?: Seat;
  tricksNS: number;
  tricksEW: number;
  score: number;
  passedOut: boolean;
}

export interface PlayerView {
  phase: Phase;
  boardNumber: number;
  dealer: Seat;
  vulnerability: "None" | "NS" | "EW" | "Both";

  seat: Seat;
  hand: string[];

  calls: string[];
  contract?: string;

  dummy?: Seat;
  dummyHand?: string[];

  currentTrick: PlayedCard[];
  lastTrick?: PlayedCard[]; // the just-finished trick, still on the table
  tricksNS: number;
  tricksEW: number;

  turn: Seat;

  result?: ResultView;

  legalCalls?: string[];
  legalCards?: string[];
}

// Client -> server
export type ClientMessage =
  | { type: "create_table" }
  | { type: "join_table"; tableID: string }
  | { type: "sit"; direction: "N" | "E" | "S" | "W" }
  | { type: "sit_bot"; direction: "N" | "E" | "S" | "W"; difficulty: number }
  | { type: "start"; seed?: number }
  | { type: "bid"; call: string }
  | { type: "play_card"; card: string };

export type SeatOccupant = "human" | "bot" | "";

/** Lobby view of a table before/around the game starting. */
export interface TableState {
  tableID: string;
  seats: Record<Seat, SeatOccupant>;
  started: boolean;
}

// Server -> client
export interface ServerMessage {
  type:
    | "table_created"
    | "table_joined"
    | "seated"
    | "game_state"
    | "table_state"
    | "error";
  payload?: unknown;
  error?: string;
}
