"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import type { ServerMessage } from "@/lib/protocol";
import { openTableSocket, type TableSocket } from "@/lib/ws";

export type DraftStatus = "connecting" | "ready" | "error";

/**
 * Used by /create (step 1): opens a socket, creates a table, and keeps
 * the socket open (so the table isn't reaped) while the settings page is
 * mounted. `setName` / `setDescription` push the table's metadata to the
 * server so it's already set by the time we reach the lobby.
 */
export function useTableDraft(): {
  tableId: string | null;
  status: DraftStatus;
  error: string | null;
  setName: (name: string) => void;
  setDescription: (description: string) => void;
} {
  const [tableId, setTableId] = useState<string | null>(null);
  const [status, setStatus] = useState<DraftStatus>("connecting");
  const [error, setError] = useState<string | null>(null);
  const socketRef = useRef<TableSocket | null>(null);

  useEffect(() => {
    const socket = openTableSocket({
      onOpen() {
        socket.send({ type: "create_table" });
      },
      onMessage(msg: ServerMessage) {
        if (msg.type === "table_created") {
          const id = (msg.payload as { tableID: string }).tableID;
          setTableId(id);
          setStatus("ready");
        } else if (msg.type === "error") {
          setError(msg.error ?? "could not create a table");
          setStatus("error");
        }
      },
      onError() {
        setError("connection failed");
        setStatus("error");
      },
    });
    socketRef.current = socket;
    return () => {
      socket.close();
      socketRef.current = null;
    };
  }, []);

  const setName = useCallback((name: string) => {
    socketRef.current?.send({ type: "set_name", name });
  }, []);
  const setDescription = useCallback((description: string) => {
    socketRef.current?.send({ type: "set_description", description });
  }, []);

  return { tableId, status, error, setName, setDescription };
}
