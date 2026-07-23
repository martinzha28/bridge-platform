"use client";

import { useState } from "react";
import Panel from "@/components/Panel";
import Button from "@/components/Button";
import { inviteUrl } from "@/lib/tableConfig";
import styles from "./InvitePanel.module.css";

export default function InvitePanel({
  tableId,
  description,
}: {
  tableId: string;
  description: string;
}) {
  const [copied, setCopied] = useState(false);
  const link =
    typeof window !== "undefined"
      ? inviteUrl(window.location.origin, tableId)
      : "";

  function copy() {
    navigator.clipboard?.writeText(link).then(
      () => {
        setCopied(true);
        setTimeout(() => setCopied(false), 1500);
      },
      () => {},
    );
  }

  return (
    <Panel title="Invite" padded>
      {description && <p className={styles.desc}>{description}</p>}
      <div className={styles.linkRow}>
        <input readOnly value={link} />
        <Button size="sm" onClick={copy}>
          {copied ? "Copied" : "Copy"}
        </Button>
      </div>
    </Panel>
  );
}
