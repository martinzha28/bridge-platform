// Thin client over the backend's cookie auth. All requests go same-origin
// to /api/* (Next proxies to the backend), so the HttpOnly token cookie
// is set and sent automatically.

export interface User {
  id: string;
  username: string;
  rating: number;
  gamesPlayed: number;
  karma: number;
  email?: string;
  about?: string;
  nationality?: string;
  systems?: string[];
  createdAt: string;
}

export type AuthResult = { ok: true; user: User } | { ok: false; error: string };

const AUTH = "/api/v1/auth";

/** The current user, or null if not logged in. */
export async function getMe(): Promise<User | null> {
  try {
    const res = await fetch(`${AUTH}/me`, { credentials: "include" });
    return res.ok ? ((await res.json()) as User) : null;
  } catch {
    return null;
  }
}

async function authRequest(
  path: string,
  username: string,
  password: string,
): Promise<AuthResult> {
  let res: Response;
  try {
    res = await fetch(`${AUTH}${path}`, {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });
  } catch {
    return { ok: false, error: "Can't reach the server." };
  }
  if (res.ok) return { ok: true, user: (await res.json()) as User };
  // the backend sends plain-text errors via http.Error
  const text = (await res.text()).trim();
  return { ok: false, error: text || `Request failed (${res.status})` };
}

export function login(username: string, password: string): Promise<AuthResult> {
  return authRequest("/login", username, password);
}

export function signup(username: string, password: string): Promise<AuthResult> {
  return authRequest("/register", username, password);
}

export function logout(): Promise<Response> {
  return fetch(`${AUTH}/logout`, { method: "POST", credentials: "include" });
}
