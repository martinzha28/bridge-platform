"use client";

import { boardsInput, parseBoards, type TableConfig } from "@/lib/tableConfig";

function Pills<T extends string>({
  value,
  options,
  onChange,
}: {
  value: T;
  options: { value: T; label: string; disabled?: boolean }[];
  onChange: (v: T) => void;
}) {
  return (
    <div className="pills">
      {options.map((o) => (
        <button
          key={o.value}
          type="button"
          disabled={o.disabled}
          className={`pill${value === o.value ? " on" : ""}`}
          onClick={() => onChange(o.value)}
        >
          {o.label}
        </button>
      ))}
    </div>
  );
}

/** Step 1: table settings only. Seats are arranged in the lobby. */
export default function CreateTableForm({
  config,
  onChange,
  tableId,
  creating = false,
  onContinue,
}: {
  config: TableConfig;
  onChange: (c: TableConfig) => void;
  tableId: string | null;
  creating?: boolean;
  onContinue: () => void;
}) {
  return (
    <div className="box create-box">
      <div className="bt">
        <span className="t">Create a New Table</span>
      </div>

      <div className="create-form">
        <label className="cf-field">
          <span>Title</span>
          <input
            className="cf-input"
            maxLength={60}
            placeholder="Practice table"
            value={config.name}
            onChange={(e) => onChange({ ...config, name: e.target.value })}
          />
        </label>

        <label className="cf-field">
          <span>Description</span>
          <textarea
            className="cf-input cf-textarea"
            maxLength={280}
            rows={4}
            placeholder="Optional — what's this table for?"
            value={config.description}
            onChange={(e) => onChange({ ...config, description: e.target.value })}
          />
        </label>

        <label className="cf-field">
          <span>Mode</span>
          <Pills
            value={config.mode}
            onChange={(mode) => onChange({ ...config, mode })}
            options={[
              { value: "casual", label: "Casual" },
              { value: "competitive", label: "Competitive" },
            ]}
          />
        </label>

        <label className="cf-field">
          <span>Visibility</span>
          <Pills
            value={config.visibility}
            onChange={(visibility) => onChange({ ...config, visibility })}
            options={[
              { value: "public", label: "Public" },
              { value: "private", label: "Private" },
            ]}
          />
        </label>

        <label className="cf-field">
          <span>Scoring</span>
          <Pills
            value={config.scoring}
            onChange={(scoring) => onChange({ ...config, scoring })}
            options={[
              { value: "imps", label: "IMPs" },
              { value: "matchpoints", label: "Matchpoints", disabled: true },
            ]}
          />
        </label>

        <label className="cf-field">
          <span>Boards</span>
          <input
            className="cf-input"
            inputMode="numeric"
            placeholder="unlimited"
            value={boardsInput(config.boards)}
            onChange={(e) => onChange({ ...config, boards: parseBoards(e.target.value) })}
          />
        </label>

        <button
          type="button"
          className="btn pri cf-create"
          disabled={!tableId || creating}
          onClick={onContinue}
        >
          {creating ? "Creating…" : "Continue to seating →"}
        </button>
        <p className="cf-hint">
          {creating
            ? "Setting up your table…"
            : "Next: pick seats, add bots, and invite players."}
        </p>
      </div>
    </div>
  );
}
