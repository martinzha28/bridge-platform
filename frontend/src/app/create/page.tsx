"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Rail from "@/components/Rail";
import Panel from "@/components/Panel";
import CreateTableForm from "./CreateTableForm";
import PlaceholderFelt from "./PlaceholderFelt";
import { useTableDraft } from "@/hooks/useTableDraft";
import { DEFAULT_CONFIG, storeConfig, type TableConfig } from "@/lib/tableConfig";
import styles from "./CreateTableForm.module.css";

export default function CreatePage() {
  const { tableId, status, error, setName, setDescription } = useTableDraft();
  const [config, setConfig] = useState<TableConfig>(DEFAULT_CONFIG);
  const [creating, setCreating] = useState(false);
  const router = useRouter();

  // Step 1 of 2: this page only captures table settings. Seating happens
  // in the lobby. Flush the metadata to the server, then hand off.
  function toSeating() {
    if (!tableId || creating) return;
    setCreating(true);
    setName(config.name);
    setDescription(config.description);
    storeConfig(tableId, config);
    router.push(`/table/${tableId}`);
  }

  return (
    <div className="app">
      <Rail active="play" />

      <div className="center">
        <PlaceholderFelt />
      </div>

      <div className="side">
        {status === "error" ? (
          <Panel title="Create a New Table" padded className={styles.panel}>
            <p className={styles.hint}>
              {error ?? "Couldn't reach the server."}
            </p>
          </Panel>
        ) : (
          <CreateTableForm
            config={config}
            onChange={setConfig}
            tableId={tableId}
            creating={creating}
            onContinue={toSeating}
          />
        )}
      </div>
    </div>
  );
}
