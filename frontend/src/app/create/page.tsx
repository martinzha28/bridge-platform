"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Rail from "@/components/Rail";
import Panel from "@/components/Panel";
import NamePrompt from "@/components/NamePrompt";
import CreateTableForm from "./CreateTableForm";
import PlaceholderFelt from "./PlaceholderFelt";
import { useTableDraft } from "@/hooks/useTableDraft";
import {
  useDisplayIdentity,
  type DisplayIdentity,
} from "@/hooks/useDisplayIdentity";
import { DEFAULT_CONFIG, storeConfig, type TableConfig } from "@/lib/tableConfig";
import styles from "./CreateTableForm.module.css";

export default function CreatePage() {
  const identity = useDisplayIdentity();

  return (
    <div className="app">
      <Rail active="play" />

      <div className="center">
        <PlaceholderFelt />
      </div>

      <div className="side">
        {!identity.ready ? null : identity.needsPrompt ? (
          <NamePrompt
            initial={identity.name}
            onSubmit={identity.chooseName}
          />
        ) : (
          <CreateDraft identity={identity} />
        )}
      </div>
    </div>
  );
}

/** Mounted only once we know who's playing: it opens the draft socket,
 *  which creates the table under that name. */
function CreateDraft({ identity }: { identity: DisplayIdentity }) {
  const { tableId, status, error, setName, setDescription } = useTableDraft(
    identity.join,
  );
  const [config, setConfig] = useState<TableConfig>(DEFAULT_CONFIG);
  const [creating, setCreating] = useState(false);
  const [renaming, setRenaming] = useState(false);
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

  if (status === "error") {
    return (
      <Panel title="Create a New Table" padded className={styles.panel}>
        <p className={styles.hint}>{error ?? "Couldn't reach the server."}</p>
      </Panel>
    );
  }

  // A rename applies when the lobby connects; the draft socket keeps the
  // table alive meanwhile, so it stays mounted.
  if (renaming) {
    return (
      <NamePrompt
        initial={identity.name}
        onSubmit={(name) => {
          identity.chooseName(name);
          setRenaming(false);
        }}
        onCancel={() => setRenaming(false)}
      />
    );
  }

  return (
    <>
      <CreateTableForm
        config={config}
        onChange={setConfig}
        tableId={tableId}
        creating={creating}
        onContinue={toSeating}
      />
      <p className={styles.playingAs}>
        Playing as <b>{identity.name}</b>
        {!identity.user && (
          <>
            {" · "}
            <button
              type="button"
              className={styles.change}
              onClick={() => setRenaming(true)}
            >
              change
            </button>
          </>
        )}
      </p>
    </>
  );
}
