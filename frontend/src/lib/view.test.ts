import { describe, it, expect } from "vitest";
import type { PlayerView } from "./protocol";
import {
  auctionColumns,
  auctionRows,
  seatVulnerable,
  boardResultText,
  boardSummary,
  canPlay,
  cardFace,
  contractSymbol,
  historyRow,
  declarerSeat,
  groupHandBySuit,
  handCount,
  hasCall,
  trumpSuit,
  nextSeat,
  normalizeView,
  partnerSeat,
  sideTricks,
  trickCardFor,
  trickJustFinished,
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

describe("trickJustFinished", () => {
  const done = (over: Partial<PlayerView>) =>
    trickJustFinished(view({ phase: "Play", currentTrick: [], lastTrick: [{ seat: "West", card: "SA" }], ...over }));

  it("is true between tricks in the play phase", () => {
    expect(done({})).toBe(true);
  });
  it("is false once the next trick is led", () => {
    expect(done({ currentTrick: [{ seat: "North", card: "S2" }] })).toBe(false);
  });
  it("is false with no completed trick", () => {
    expect(done({ lastTrick: [] })).toBe(false);
  });
  it("is false outside the play phase", () => {
    expect(done({ phase: "Complete" })).toBe(false);
  });
});

describe("auctionColumns", () => {
  it("puts the viewer's own seat last, clockwise order ahead of it", () => {
    expect(auctionColumns("South")).toEqual(["West", "North", "East", "South"]);
    expect(auctionColumns("North")).toEqual(["East", "South", "West", "North"]);
    expect(auctionColumns("East")).toEqual(["South", "West", "North", "East"]);
    expect(auctionColumns("West")).toEqual(["North", "East", "South", "West"]);
  });
});

describe("sideTricks", () => {
  it("splits trick counts from the viewer's side", () => {
    expect(sideTricks(view({ seat: "South", tricksNS: 5, tricksEW: 3 }))).toEqual({
      ours: 5,
      theirs: 3,
    });
    expect(sideTricks(view({ seat: "East", tricksNS: 5, tricksEW: 3 }))).toEqual({
      ours: 3,
      theirs: 5,
    });
  });
});

describe("seatVulnerable", () => {
  it("reads the board vulnerability per partnership", () => {
    expect(seatVulnerable("South", "EW")).toBe(false);
    expect(seatVulnerable("West", "EW")).toBe(true);
    expect(seatVulnerable("North", "NS")).toBe(true);
    expect(seatVulnerable("East", "Both")).toBe(true);
    expect(seatVulnerable("South", "None")).toBe(false);
  });
});

describe("auctionRows", () => {
  const south = auctionColumns("South"); // W N E S

  it("pads leading cells so the dealer sits in its column", () => {
    // dealer North -> one blank under West
    expect(auctionRows(["1C", "P", "1H", "P"], "North", south)).toEqual([
      ["", "1C", "P", "1H"],
      ["P", "", "", ""],
    ]);
  });
  it("starts flush when the first column dealt", () => {
    expect(auctionRows(["P", "P"], "West", south)).toEqual([["P", "P", "", ""]]);
  });
  it("respects a rotated column order", () => {
    // viewer North: columns E S W N, dealer North -> three blanks
    expect(auctionRows(["1C", "P"], "North", auctionColumns("North"))).toEqual([
      ["", "", "", "1C"],
      ["P", "", "", ""],
    ]);
  });
  it("returns nothing for an empty auction", () => {
    expect(auctionRows([], "North", south)).toEqual([]);
  });
});

describe("boardSummary", () => {
  it("strips the declarer from the contract and takes that side's tricks", () => {
    expect(
      boardSummary(
        view({
          boardNumber: 3,
          vulnerability: "EW",
          contract: "3NT by North",
          result: {
            contract: "3NT by North",
            declarer: "North",
            tricksNS: 10,
            tricksEW: 3,
            score: 430,
            passedOut: false,
          },
        }),
      ),
    ).toEqual({
      board: 3,
      vulnerability: "EW",
      contract: "3NT",
      declarer: "North",
      tricks: 10,
      score: 430,
      passedOut: false,
    });
  });
  it("uses EW tricks for an EW declarer", () => {
    expect(
      boardSummary(
        view({
          result: {
            contract: "4S by West",
            declarer: "West",
            tricksNS: 4,
            tricksEW: 9,
            score: -50,
            passedOut: false,
          },
        }),
      ),
    ).toMatchObject({ contract: "4S", declarer: "West", tricks: 9, score: -50 });
  });
  it("marks a passed-out board", () => {
    expect(
      boardSummary(
        view({ boardNumber: 2, result: { tricksNS: 0, tricksEW: 0, score: 0, passedOut: true } }),
      ),
    ).toMatchObject({ board: 2, contract: null, declarer: null, passedOut: true });
  });
});

describe("contractSymbol", () => {
  it("swaps the strain letter for a suit symbol", () => {
    expect(contractSymbol("1S")).toBe("1♠");
    expect(contractSymbol("4H")).toBe("4♥");
    expect(contractSymbol("6D")).toBe("6♦");
    expect(contractSymbol("2C")).toBe("2♣");
  });
  it("keeps NT and any doubling suffix", () => {
    expect(contractSymbol("3NT")).toBe("3NT");
    expect(contractSymbol("4HX")).toBe("4♥X");
    expect(contractSymbol("6DXX")).toBe("6♦XX");
  });
  it("passes odd input through", () => {
    expect(contractSymbol("passed out")).toBe("passed out");
  });
});

describe("historyRow", () => {
  const base = { board: 1, vulnerability: "None", passedOut: false };

  it("formats a made N/S contract with the score under NS", () => {
    expect(
      historyRow({ ...base, contract: "3NT", declarer: "North", tricks: 10, score: 430 }),
    ).toEqual({ board: 1, result: "3NT N +1", red: false, scoreNS: "430", scoreEW: "–" });
  });
  it("shows '=' for an exactly-made contract", () => {
    expect(
      historyRow({ ...base, contract: "4S", declarer: "North", tricks: 10, score: 420 }),
    ).toMatchObject({ result: "4♠ N =", scoreNS: "420", scoreEW: "–" });
  });
  it("puts a defeated N/S contract's penalty under EW", () => {
    expect(
      historyRow({ ...base, contract: "1NT", declarer: "South", tricks: 0, score: -700 }),
    ).toMatchObject({ result: "1NT S -7", scoreNS: "–", scoreEW: "700" });
  });
  it("handles an E/W declarer both ways", () => {
    expect(
      historyRow({ ...base, contract: "4H", declarer: "West", tricks: 11, score: 450 }),
    ).toMatchObject({ result: "4♥ W +1", red: true, scoreNS: "–", scoreEW: "450" });
    expect(
      historyRow({ ...base, contract: "2S", declarer: "East", tricks: 6, score: -100 }),
    ).toMatchObject({ result: "2♠ E -2", scoreNS: "100", scoreEW: "–" });
  });
  it("marks a passed-out board", () => {
    expect(
      historyRow({ ...base, board: 2, contract: null, declarer: null, tricks: 0, score: 0, passedOut: true }),
    ).toEqual({ board: 2, result: "passed out", red: false, scoreNS: "–", scoreEW: "–" });
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
