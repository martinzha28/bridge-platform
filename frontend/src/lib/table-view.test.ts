import { describe, it, expect } from "vitest";
import { peopleBySeat } from "./table-view";
import type { Person } from "./protocol";

describe("peopleBySeat", () => {
  it("indexes seated people by seat and skips watchers", () => {
    const people: Person[] = [
      { id: 1, name: "alice", guest: false, seat: "South", host: true },
      { id: 2, name: "Guest 4821", guest: true, seat: "", host: false },
      { id: 3, name: "bob", guest: true, seat: "North", host: false },
    ];
    expect(peopleBySeat(people)).toEqual({ South: people[0], North: people[2] });
  });
  it("tolerates a missing list", () => {
    expect(peopleBySeat(undefined)).toEqual({});
  });
});
