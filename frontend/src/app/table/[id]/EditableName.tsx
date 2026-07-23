"use client";

import { useEffect, useState } from "react";
import styles from "./EditableName.module.css";

/** Click-to-edit table title. Used in the lobby header and the play header. */
export default function EditableName({
  name,
  onRename,
}: {
  name: string;
  onRename: (v: string) => void;
}) {
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState(name);

  useEffect(() => {
    if (!editing) setDraft(name);
  }, [name, editing]);

  if (!editing) {
    return (
      <h1
        className={styles.name}
        title="Click to rename"
        onClick={() => {
          setDraft(name);
          setEditing(true);
        }}
      >
        {name}
      </h1>
    );
  }

  function commit() {
    setEditing(false);
    const v = draft.trim();
    if (v && v !== name) onRename(v);
  }

  return (
    <input
      className={styles.input}
      autoFocus
      maxLength={60}
      value={draft}
      onChange={(e) => setDraft(e.target.value)}
      onBlur={commit}
      onKeyDown={(e) => {
        if (e.key === "Enter") commit();
        if (e.key === "Escape") {
          setDraft(name);
          setEditing(false);
        }
      }}
    />
  );
}
