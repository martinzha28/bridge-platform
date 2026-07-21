import { describe, it, expect } from "vitest";
import {
  DEFAULT_CONFIG,
  boardsInput,
  configValid,
  cycleSeat,
  inviteUrl,
  parseBoards,
  seatPlan,
} from "./tableConfig";

describe("parseBoards", () => {
  it("empty → unlimited", () => {
    expect(parseBoards("")).toBeNull();
    expect(parseBoards("   ")).toBeNull();
  });
  it("parses and clamps", () => {
    expect(parseBoards("24")).toBe(24);
    expect(parseBoards("1")).toBe(1);
    expect(parseBoards("999")).toBe(64);
    expect(parseBoards("3.9")).toBe(3);
  });
  it("junk / non-positive → unlimited", () => {
    expect(parseBoards("0")).toBeNull();
    expect(parseBoards("-5")).toBeNull();
    expect(parseBoards("abc")).toBeNull();
  });
});

describe("boardsInput", () => {
  it("round-trips with parseBoards", () => {
    expect(boardsInput(24)).toBe("24");
    expect(boardsInput(null)).toBe("");
  });
});

describe("cycleSeat", () => {
  const seats = DEFAULT_CONFIG.seats; // S=you, rest=bot

  it("cycles a seat you → bot → open → you", () => {
    let s = seats;
    s = cycleSeat(s, "North"); // bot → open
    expect(s.North).toBe("open");
    s = cycleSeat(s, "North"); // open → you
    expect(s.North).toBe("you");
    s = cycleSeat(s, "North"); // you → bot
    expect(s.North).toBe("bot");
  });
  it("landing on 'you' demotes the previous 'you' to bot", () => {
    let s = seats;
    s = cycleSeat(s, "North"); // bot → open
    s = cycleSeat(s, "North"); // open → you
    expect(s.North).toBe("you");
    expect(s.South).toBe("bot"); // was "you"
  });
});

describe("seatPlan", () => {
  it("splits seats by role", () => {
    const plan = seatPlan({
      ...DEFAULT_CONFIG,
      seats: { South: "you", North: "bot", East: "open", West: "bot" },
    });
    expect(plan).toEqual({ you: "South", bots: ["North", "West"], open: ["East"] });
  });
});

describe("configValid", () => {
  it("requires exactly one 'you'", () => {
    expect(configValid(DEFAULT_CONFIG)).toBe(true);
    expect(
      configValid({ ...DEFAULT_CONFIG, seats: { South: "bot", North: "bot", East: "bot", West: "bot" } }),
    ).toBe(false);
    expect(
      configValid({ ...DEFAULT_CONFIG, seats: { South: "you", North: "you", East: "bot", West: "bot" } }),
    ).toBe(false);
  });
});

describe("inviteUrl", () => {
  it("builds the table link", () => {
    expect(inviteUrl("https://x.app", "a3f9")).toBe("https://x.app/table/a3f9");
  });
});
