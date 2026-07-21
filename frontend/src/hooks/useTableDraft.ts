"use client";

import { useEffect, useRef, useState } from "react";
import type { ServerMessage } from "@/lib/protocol";
import { openTableSocket, type TableSocket } from "@/lib/ws";

export type DraftStatus = "connecting" | "ready" | "error";

/**
 * Used by /create: opens a socket, creates a table, and keeps the socket
 * open (so the table isn't reaped) while the page is mounted. The
 * returned `tableId` powers the invite link before anything is
 * configured.
 */
export function useTableDraft(): {
  tableId: string | null;
  status: DraftStatus;
  error: string | null;
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

  return { tableId, status, error };
}
