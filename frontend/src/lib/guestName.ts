// A guest's display identity, kept in localStorage so it survives
// refreshes and carries across tables. Mirrors the backend's name rules
// (ws/identity.go): 2–20 letters, digits, spaces, '_' or '-'.

export interface GuestIdentity {
  name: string;
  /** Random per-browser id; lets the server recognise a reconnect. */
  guestID: string;
  /** False until the guest confirms a name in the "Play as" prompt. */
  chosen: boolean;
}

const KEY = "guest-identity";
export const NAME_MIN = 2;
export const NAME_MAX = 20;
const NAME_CHARS = /^[\p{L}\p{Nd} _-]+$/u;

// used when localStorage is unavailable, so the identity is at least
// stable for this page session
let memoryIdentity: GuestIdentity | null = null;

export function normalizeName(raw: string): string {
  return raw.trim().split(/\s+/).filter(Boolean).join(" ");
}

export function isValidName(raw: string): boolean {
  const name = normalizeName(raw);
  const len = [...name].length;
  return len >= NAME_MIN && len <= NAME_MAX && NAME_CHARS.test(name);
}

export function randomGuestName(rand: () => number = Math.random): string {
  return `Guest ${1000 + Math.floor(rand() * 9000)}`;
}

function randomID(): string {
  if (typeof crypto !== "undefined" && "randomUUID" in crypto) {
    return crypto.randomUUID();
  }
  return Math.random().toString(36).slice(2) + Date.now().toString(36);
}

function defaultStorage(): Storage | null {
  try {
    return typeof localStorage === "undefined" ? null : localStorage;
  } catch {
    return null;
  }
}

function read(storage: Storage | null): GuestIdentity | null {
  try {
    const raw = storage?.getItem(KEY);
    if (!raw) return null;
    const v = JSON.parse(raw) as Partial<GuestIdentity>;
    if (typeof v.name !== "string" || !isValidName(v.name)) return null;
    if (typeof v.guestID !== "string" || v.guestID === "") return null;
    return { name: v.name, guestID: v.guestID, chosen: v.chosen === true };
  } catch {
    return null;
  }
}

function persists(storage: Storage | null): boolean {
  try {
    storage?.setItem(`${KEY}:probe`, "1");
    storage?.removeItem(`${KEY}:probe`);
    return storage != null;
  } catch {
    return false;
  }
}

function write(storage: Storage | null, id: GuestIdentity): void {
  memoryIdentity = id;
  try {
    storage?.setItem(KEY, JSON.stringify(id));
  } catch {
    // storage blocked — memoryIdentity still holds it for this session
  }
}

/** The stored guest identity, creating (and saving) a default if none. */
export function loadGuestIdentity(
  storage: Storage | null = defaultStorage(),
): GuestIdentity {
  const stored = read(storage);
  if (stored) return stored;
  if (memoryIdentity && !persists(storage)) return memoryIdentity;
  const fresh = { name: randomGuestName(), guestID: randomID(), chosen: false };
  write(storage, fresh);
  return fresh;
}

/** Saves a confirmed nickname, keeping the guest id. */
export function saveGuestName(
  name: string,
  storage: Storage | null = defaultStorage(),
): GuestIdentity {
  const current = loadGuestIdentity(storage);
  const next = { ...current, name: normalizeName(name), chosen: true };
  write(storage, next);
  return next;
}
