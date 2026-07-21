"use client";

import { useEffect, useRef, useState } from "react";
import type { ChatMessage } from "@/lib/protocol";

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
    <div className="box chat-panel">
      <div className="bt">
        <span className="t">Table chat</span>
      </div>

      <div className="chat-log" ref={logRef}>
        {messages.length === 0 ? (
          <p className="chat-empty">No messages yet.</p>
        ) : (
          messages.map((m) => (
            <p key={m.id} className="chat-row">
              <b className="chat-from">{m.sender}</b>
              <span>{m.text}</span>
            </p>
          ))
        )}
      </div>

      <div className="chat-input-row">
        <input
          className="chat-input"
          value={draft}
          maxLength={500}
          placeholder="Say something…"
          onChange={(e) => setDraft(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") send();
          }}
        />
        <button type="button" className="btn xs" disabled={!draft.trim()} onClick={send}>
          Send
        </button>
      </div>
    </div>
  );
}
