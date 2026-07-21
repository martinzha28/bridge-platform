"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Rail from "@/components/Rail";
import CreateTableForm from "@/components/CreateTableForm";
import { useTableDraft } from "@/hooks/useTableDraft";
import { SEATS } from "@/lib/protocol";
import { DEFAULT_CONFIG, storeConfig, type TableConfig } from "@/lib/tableConfig";

const FELT_POS = { North: "n", East: "e", South: "s", West: "w" } as const;

function CardBacks({ count }: { count: number }) {
  return (
    <>
      {Array.from({ length: count }).map((_, i) => (
        <div key={i} className="cb" />
      ))}
    </>
  );
}

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
        <div className="felt-box">
          <div className="felt">
            {SEATS.map((seat) => (
              <span key={seat} className={`mk ${FELT_POS[seat]}`}>
                {seat[0]}
              </span>
            ))}
            <div className="seat n">
              <CardBacks count={13} />
            </div>
            <div className="seat w">
              <CardBacks count={13} />
            </div>
            <div className="seat e">
              <CardBacks count={13} />
            </div>
            <div className="seat s">
              <CardBacks count={13} />
            </div>
          </div>
        </div>
      </div>

      <div className="side">
        {status === "error" ? (
          <div className="box create-box">
            <div className="bt">
              <span className="t">Create a New Table</span>
            </div>
            <div className="create-form">
              <p className="cf-hint">{error ?? "Couldn't reach the server."}</p>
            </div>
          </div>
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
