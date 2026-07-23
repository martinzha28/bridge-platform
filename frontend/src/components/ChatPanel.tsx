"use client";

import { useEffect, useRef, useState } from "react";
import Panel from "@/components/Panel";
import Button from "@/components/Button";
import type { ChatMessage } from "@/lib/protocol";
import styles from "./ChatPanel.module.css";

/** Ephemeral table chat: recent scrollback plus a send box. Open to
 *  everyone at the table in every phase. */
export default function ChatPanel({
  messages,
  onSend,
}: {
  messages: ChatMessage[];
  onSend: (text: string) => void;
}) {
  const [draft, setDraft] = useState("");
  const logRef = useRef<HTMLDivElement>(null);

  // keep the newest line in view
  useEffect(() => {
    const el = logRef.current;
    if (el) el.scrollTop = el.scrollHeight;
  }, [messages]);

  function send() {
    const text = draft.trim();
    if (!text) return;
    onSend(text);
    setDraft("");
  }

  return (
    <Panel title="Table chat" className={styles.panel}>
      <div className={styles.log} ref={logRef}>
        {messages.length === 0 ? (
          <p className={styles.empty}>No messages yet.</p>
        ) : (
          messages.map((m) => (
            <p key={m.id} className={styles.row}>
              <b className={styles.from}>{m.sender}</b>
              <span>{m.text}</span>
            </p>
          ))
        )}
      </div>

      <div className={styles.inputRow}>
        <input
          className={styles.input}
          value={draft}
          maxLength={500}
          placeholder="Say something…"
          onChange={(e) => setDraft(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") send();
          }}
        />
        <Button size="sm" disabled={!draft.trim()} onClick={send}>
          Send
        </Button>
      </div>
    </Panel>
  );
}
