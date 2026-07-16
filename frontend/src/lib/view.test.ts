import { describe, it, expect } from "vitest";
import type { PlayerView } from "./protocol";
import {
  auctionRows,
  boardResultText,
  canPlay,
  cardFace,
  declarerSeat,
  groupHandBySuit,
  handCount,
  hasCall,
  trumpSuit,
  nextSeat,
  normalizeView,
  partnerSeat,
  trickCardFor,
} from "./view";

function view(overrides: Partial<PlayerView> = {}): PlayerView {
  return {
    phase: "Play",
    boardNumber: 1,
    dealer: "North",
    vulnerability: "None",
    seat: "South",
    hand: [],
    calls: [],
    currentTrick: [],
    tricksNS: 0,
    tricksEW: 0,
    turn: "South",
    ...overrides,
  };
}

describe("seat helpers", () => {
  it("advances clockwise", () => {
    expect(nextSeat("North")).toBe("East");
    expect(nextSeat("West")).toBe("North");
  });
  it("finds partner across the table", () => {
    expect(partnerSeat("North")).toBe("South");
    expect(partnerSeat("East")).toBe("West");
  });
});

describe("cardFace", () => {
  it("maps suit letters to symbols and colours", () => {
    expect(cardFace("SA")).toEqual({ rank: "A", suit: "♠", red: false });
    expect(cardFace("HT")).toEqual({ rank: "10", suit: "♥", red: true });
    expect(cardFace("D2")).toEqual({ rank: "2", suit: "♦", red: true });
    expect(cardFace("CK")).toEqual({ rank: "K", suit: "♣", red: false });
  });
});

describe("handCount", () => {
  it("is 13 for every seat before a card is played", () => {
    const v = view({ phase: "Auction" });
    expect(handCount(v, "North")).toBe(13);
    expect(handCount(v, "South")).toBe(13);
  });
  it("subtracts completed tricks and a card in the live trick", () => {
    const v = view({
      tricksNS: 2,
      tricksEW: 1,
      currentTrick: [
        { seat: "West", card: "SA" },
        { seat: "North", card: "S2" },
      ],
    });
    expect(handCount(v, "West")).toBe(13 - 3 - 1);
    expect(handCount(v, "North")).toBe(13 - 3 - 1);
    expect(handCount(v, "South")).toBe(13 - 3 - 0);
  });
});

describe("declarerSeat", () => {
  it("is null until a dummy exists", () => {
    expect(declarerSeat(view({ dummy: undefined }))).toBeNull();
  });
  it("is the dummy's partner", () => {
    expect(declarerSeat(view({ dummy: "North" }))).toBe("South");
  });
});

describe("canPlay", () => {
  it("is true on the client's own turn", () => {
    expect(canPlay(view({ turn: "South", seat: "South" }))).toBe(true);
  });
  it("is true when declarer and it is dummy's turn", () => {
    expect(canPlay(view({ seat: "South", dummy: "North", turn: "North" }))).toBe(true);
  });
  it("is false in the auction", () => {
    expect(canPlay(view({ phase: "Auction" }))).toBe(false);
  });
  it("is false on an opponent's turn", () => {
    expect(canPlay(view({ seat: "South", turn: "East", dummy: "West" }))).toBe(false);
  });
});

describe("trickCardFor", () => {
  it("returns the card a seat played to the live trick", () => {
    const v = view({ currentTrick: [{ seat: "West", card: "HQ" }] });
    expect(trickCardFor(v, "West")).toBe("HQ");
    expect(trickCardFor(v, "North")).toBeUndefined();
  });
});

describe("auctionRows", () => {
  it("pads leading cells so the dealer sits in the right column", () => {
    // dealer North -> one blank under West
    expect(auctionRows(["1C", "P", "1H", "P"], "North")).toEqual([
      ["", "1C", "P", "1H"],
      ["P", "", "", ""],
    ]);
  });
  it("starts flush when West deals", () => {
    expect(auctionRows(["P", "P"], "West")).toEqual([["P", "P", "", ""]]);
  });
  it("returns nothing for an empty auction", () => {
    expect(auctionRows([], "North")).toEqual([]);
  });
});

describe("boardResultText", () => {
  it("summarises a made contract using the declarer's side tricks", () => {
    expect(
      boardResultText(
        view({
          phase: "Complete",
          result: {
            contract: "1H by South",
            declarer: "South",
            tricksNS: 8,
            tricksEW: 5,
            score: 110,
            passedOut: false,
          },
        }),
      ),
    ).toBe("1H by South · 8 tricks · +110");
  });
  it("uses EW tricks and a signed score for an EW declarer going down", () => {
    expect(
      boardResultText(
        view({
          result: {
            contract: "4S by West",
            declarer: "West",
            tricksNS: 3,
            tricksEW: 10,
            score: -100,
            passedOut: false,
          },
        }),
      ),
    ).toBe("4S by West · 10 tricks · -100");
  });
  it("handles a passed-out board", () => {
    expect(
      boardResultText(
        view({ result: { tricksNS: 0, tricksEW: 0, score: 0, passedOut: true } }),
      ),
    ).toBe("Passed out");
  });
});

describe("trumpSuit", () => {
  it("reads the strain from a contract string", () => {
    expect(trumpSuit("6D by South")).toBe("D");
    expect(trumpSuit("4HX by North")).toBe("H");
    expect(trumpSuit("7SXX by East")).toBe("S");
  });
  it("is null for no-trump and no contract", () => {
    expect(trumpSuit("1NT by South")).toBeNull();
    expect(trumpSuit("3NT by West")).toBeNull();
    expect(trumpSuit(undefined)).toBeNull();
  });
});

describe("groupHandBySuit", () => {
  const hand = ["SA", "SK", "H9", "H2", "DQ", "C7", "C3"];
  it("uses S/H/D/C order with no trump", () => {
    expect(groupHandBySuit(hand, null)).toEqual([
      ["SA", "SK"],
      ["H9", "H2"],
      ["DQ"],
      ["C7", "C3"],
    ]);
  });
  it("moves the trump suit to the front", () => {
    expect(groupHandBySuit(hand, "D")).toEqual([
      ["DQ"],
      ["SA", "SK"],
      ["H9", "H2"],
      ["C7", "C3"],
    ]);
  });
  it("drops empty suits", () => {
    expect(groupHandBySuit(["SA", "SK"], "H")).toEqual([["SA", "SK"]]);
  });
});

describe("normalizeView", () => {
  it("coerces nulled slices to arrays", () => {
    const raw = {
      ...view({ phase: "Auction" }),
      hand: null,
      calls: null,
      currentTrick: null,
    } as unknown as PlayerView;
    const v = normalizeView(raw);
    expect(v.hand).toEqual([]);
    expect(v.calls).toEqual([]);
    expect(v.currentTrick).toEqual([]);
  });
  it("leaves real data untouched", () => {
    const v = normalizeView(view({ calls: ["P", "1C"] }));
    expect(v.calls).toEqual(["P", "1C"]);
  });
});

describe("hasCall", () => {
  it("checks membership safely", () => {
    expect(hasCall(["P", "1C"], "1C")).toBe(true);
    expect(hasCall(["P", "1C"], "X")).toBe(false);
    expect(hasCall(undefined, "P")).toBe(false);
  });
});
