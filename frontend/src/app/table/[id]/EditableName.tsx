"use client";

import { useEffect, useState } from "react";
import { cx } from "@/lib/cx";
import styles from "./EditableName.module.css";

/** Click-to-edit table title. Used in the lobby header and the play
 *  header. Read-only when `onRename` is omitted (non-hosts). */
export default function EditableName({
  name,
  onRename,
}: {
  name: string;
  onRename?: (v: string) => void;
}) {
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState(name);

  useEffect(() => {
    if (!editing) setDraft(name);
  }, [name, editing]);

  if (!onRename) {
    return <h1 className={styles.name}>{name}</h1>;
  }
  const rename = onRename;

  if (!editing) {
    return (
      <h1
        className={cx(styles.name, styles.editable)}
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
    if (v && v !== name) rename(v);
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
