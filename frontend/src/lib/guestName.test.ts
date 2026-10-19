import { describe, it, expect } from "vitest";
import {
  isValidName,
  loadGuestIdentity,
  normalizeName,
  randomGuestName,
  saveGuestName,
} from "./guestName";

/** Minimal in-memory Storage. */
function memoryStorage(): Storage {
  const m = new Map<string, string>();
  return {
    get length() {
      return m.size;
    },
    clear: () => m.clear(),
    getItem: (k) => m.get(k) ?? null,
    key: (i) => [...m.keys()][i] ?? null,
    removeItem: (k) => void m.delete(k),
    setItem: (k, v) => void m.set(k, String(v)),
  };
}

function throwingStorage(): Storage {
  const fail = () => {
    throw new Error("blocked");
  };
  return {
    length: 0,
    clear: fail,
    getItem: fail,
    key: fail,
    removeItem: fail,
    setItem: fail,
  };
}

describe("normalizeName", () => {
  it("trims and collapses whitespace", () => {
    expect(normalizeName("  a   b  ")).toBe("a b");
  });
});

describe("isValidName", () => {
  it("accepts 2–20 letters, digits, spaces, _ and -", () => {
    expect(isValidName("Martin")).toBe(true);
    expect(isValidName("José Ñ")).toBe(true);
    expect(isValidName("under_score-1")).toBe(true);
    expect(isValidName("x".repeat(20))).toBe(true);
  });
  it("rejects too short, too long, and other characters", () => {
    expect(isValidName("a")).toBe(false);
    expect(isValidName("   ")).toBe(false);
    expect(isValidName("x".repeat(21))).toBe(false);
    expect(isValidName("bad<name>")).toBe(false);
    expect(isValidName("emoji 🃏")).toBe(false);
  });
  it("judges the normalized form", () => {
    expect(isValidName("  a  b ")).toBe(true);
  });
});

describe("randomGuestName", () => {
  it("is Guest + four digits", () => {
    expect(randomGuestName(() => 0)).toBe("Guest 1000");
    expect(randomGuestName(() => 0.99999)).toBe("Guest 9999");
    expect(randomGuestName()).toMatch(/^Guest \d{4}$/);
  });
});

describe("loadGuestIdentity", () => {
  it("creates and persists an unchosen default the first time", () => {
    const s = memoryStorage();
    const first = loadGuestIdentity(s);
    expect(first.name).toMatch(/^Guest \d{4}$/);
    expect(first.guestID).not.toBe("");
    expect(first.chosen).toBe(false);
    expect(loadGuestIdentity(s)).toEqual(first);
  });

  it("replaces a corrupt or invalid stored value", () => {
    const s = memoryStorage();
    s.setItem("guest-identity", "{not json");
    expect(loadGuestIdentity(s).name).toMatch(/^Guest \d{4}$/);
    s.setItem("guest-identity", JSON.stringify({ name: "x", guestID: "g", chosen: true }));
    expect(loadGuestIdentity(s).chosen).toBe(false);
  });

  it("still returns an identity when storage throws", () => {
    const id = loadGuestIdentity(throwingStorage());
    expect(id.name).toMatch(/^Guest \d{4}$/);
    expect(id.guestID).not.toBe("");
  });
});

describe("saveGuestName", () => {
  it("normalizes, marks chosen, and keeps the guest id", () => {
    const s = memoryStorage();
    const before = loadGuestIdentity(s);
    const after = saveGuestName("  Martin   Z ", s);
    expect(after).toEqual({ name: "Martin Z", guestID: before.guestID, chosen: true });
    expect(loadGuestIdentity(s)).toEqual(after);
  });
});
