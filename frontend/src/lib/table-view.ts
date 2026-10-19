import type { Person, PlayerView, Seat } from "./protocol";
import { SEATS } from "./protocol";
import { nextSeat, seatVulnerable } from "./view";
import type { TableStatus } from "@/hooks/useTable";

// Screen slot for a compass seat, rotated so the viewer is at the bottom
// and the table turns clockwise from there.
export const SCREEN_SLOTS = ["s", "w", "n", "e"] as const;

export function screenSlot(
  seat: Seat,
  mySeat: Seat,
): (typeof SCREEN_SLOTS)[number] {
  return SCREEN_SLOTS[(SEATS.indexOf(seat) - SEATS.indexOf(mySeat) + 4) % 4];
}

/** Seats in plate order: the viewer first, then clockwise round the table. */
export function plateOrder(mySeat: Seat): Seat[] {
  const i = SEATS.indexOf(mySeat);
  return [...SEATS.slice(i), ...SEATS.slice(0, i)];
}

/** Seated people keyed by seat; watchers are left out. */
export function peopleBySeat(
  people: Person[] | undefined,
): Partial<Record<Seat, Person>> {
  const out: Partial<Record<Seat, Person>> = {};
  for (const p of people ?? []) {
    if (p.seat) out[p.seat] = p;
  }
  return out;
}

const VUL_ON = "var(--red)";
const VUL_OFF = "var(--vul-safe)";

/** Edge colours from the viewer's seat: `axis` is the viewer's own
 *  partnership (top/bottom of the table), `cross` is the opponents
 *  (left/right). */
export function vulEdges(view: PlayerView | null): {
  axis: string;
  cross: string;
} {
  if (!view) return { axis: VUL_OFF, cross: VUL_OFF };
  return {
    axis: seatVulnerable(view.seat, view.vulnerability) ? VUL_ON : VUL_OFF,
    cross: seatVulnerable(nextSeat(view.seat), view.vulnerability)
      ? VUL_ON
      : VUL_OFF,
  };
}

export function statusLabel(status: TableStatus): string {
  return status === "closed" ? "disconnected" : "connecting…";
}
