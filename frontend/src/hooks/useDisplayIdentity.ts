"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { useAuth } from "@/hooks/useAuth";
import type { User } from "@/lib/auth";
import type { JoinIdentity } from "@/lib/protocol";
import {
  loadGuestIdentity,
  saveGuestName,
  type GuestIdentity,
} from "@/lib/guestName";

export interface DisplayIdentity {
  /** False until auth and the stored guest identity have both loaded. */
  ready: boolean;
  user: User | null;
  guest: GuestIdentity | null;
  /** A guest who has never confirmed a name gets the "Play as" prompt. */
  needsPrompt: boolean;
  /** Fields to send with create_table / join_table. */
  join: JoinIdentity;
  /** The name this client plays under (before any server-side suffix). */
  name: string;
  chooseName: (name: string) => void;
}

/** Logged-in users play as their username; guests as a stored nickname. */
export function useDisplayIdentity(): DisplayIdentity {
  const { user, loading } = useAuth();
  const [guest, setGuest] = useState<GuestIdentity | null>(null);

  // localStorage only exists in the browser — load after mount
  useEffect(() => {
    setGuest(loadGuestIdentity());
  }, []);

  const chooseName = useCallback((name: string) => {
    setGuest(saveGuestName(name));
  }, []);

  const join = useMemo<JoinIdentity>(
    () =>
      user || !guest ? {} : { displayName: guest.name, guestID: guest.guestID },
    [user, guest],
  );

  const ready = !loading && (user != null || guest != null);
  return {
    ready,
    user,
    guest,
    needsPrompt: ready && user == null && guest != null && !guest.chosen,
    join,
    name: user?.username ?? guest?.name ?? "",
    chooseName,
  };
}
