import { describe, it, expect } from "vitest";
import {
  DEFAULT_CONFIG,
  boardsInput,
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

describe("seatPlan", () => {
  it("splits seats by role", () => {
    const plan = seatPlan({
      ...DEFAULT_CONFIG,
      seats: { South: "you", North: "bot", East: "open", West: "bot" },
    });
    expect(plan).toEqual({ you: "South", bots: ["North", "West"], open: ["East"] });
  });
});

describe("inviteUrl", () => {
  it("builds the table link", () => {
    expect(inviteUrl("https://x.app", "a3f9")).toBe("https://x.app/table/a3f9");
  });
});
