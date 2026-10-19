"use client";

import { useState } from "react";
import Link from "next/link";
import Panel from "@/components/Panel";
import Button from "@/components/Button";
import { Field, Input } from "@/components/Field";
import { NAME_MAX, NAME_MIN, isValidName } from "@/lib/guestName";
import styles from "./NamePrompt.module.css";

/** "Play as" — a guest confirms or edits their nickname before joining. */
export default function NamePrompt({
  initial,
  onSubmit,
  onCancel,
}: {
  initial: string;
  onSubmit: (name: string) => void;
  /** Shown as a Cancel button when changing an already-chosen name. */
  onCancel?: () => void;
}) {
  const [draft, setDraft] = useState(initial);
  const valid = isValidName(draft);

  function submit() {
    if (valid) onSubmit(draft);
  }

  return (
    <Panel title="Play as" padded className={styles.panel}>
      <Field label="Display name">
        <Input
          autoFocus
          maxLength={40}
          value={draft}
          onChange={(e) => setDraft(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") submit();
          }}
        />
      </Field>
      <p className={valid ? styles.hint : styles.error}>
        {NAME_MIN}–{NAME_MAX} letters, numbers, spaces, _ or -.{" "}
        <Link href="/login">Log in</Link> to play under your username.
      </p>
      <div className={styles.actions}>
        {onCancel && (
          <Button size="sm" onClick={onCancel}>
            Cancel
          </Button>
        )}
        <Button variant="primary" disabled={!valid} onClick={submit}>
          Continue
        </Button>
      </div>
    </Panel>
  );
}
