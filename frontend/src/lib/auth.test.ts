import { describe, it, expect, vi, afterEach } from "vitest";
import { getMe, login, signup } from "./auth";

function stubFetch(res: { ok: boolean; status?: number; json?: unknown; text?: string }) {
  vi.stubGlobal(
    "fetch",
    vi.fn(async () => ({
      ok: res.ok,
      status: res.status ?? (res.ok ? 200 : 400),
      json: async () => res.json,
      text: async () => res.text ?? "",
    })),
  );
}

function stubFetchThrows() {
  vi.stubGlobal(
    "fetch",
    vi.fn(async () => {
      throw new Error("network down");
    }),
  );
}

afterEach(() => vi.unstubAllGlobals());

describe("getMe", () => {
  it("returns the user on 200", async () => {
    stubFetch({ ok: true, json: { id: "1", username: "ann" } });
    expect(await getMe()).toMatchObject({ username: "ann" });
  });
  it("returns null on 401", async () => {
    stubFetch({ ok: false, status: 401 });
    expect(await getMe()).toBeNull();
  });
  it("returns null when the request fails", async () => {
    stubFetchThrows();
    expect(await getMe()).toBeNull();
  });
});

describe("login / signup", () => {
  it("returns ok + user on success", async () => {
    stubFetch({ ok: true, json: { id: "1", username: "ann" } });
    expect(await login("ann", "pw")).toEqual({ ok: true, user: { id: "1", username: "ann" } });
  });
  it("passes the server's error text through, trimmed", async () => {
    stubFetch({ ok: false, status: 409, text: "username already taken\n" });
    expect(await signup("ann", "pw")).toEqual({ ok: false, error: "username already taken" });
  });
  it("falls back to a status message when the body is empty", async () => {
    stubFetch({ ok: false, status: 500, text: "" });
    expect(await login("ann", "pw")).toEqual({ ok: false, error: "Request failed (500)" });
  });
  it("reports an unreachable server", async () => {
    stubFetchThrows();
    expect(await login("ann", "pw")).toEqual({ ok: false, error: "Can't reach the server." });
  });
});
