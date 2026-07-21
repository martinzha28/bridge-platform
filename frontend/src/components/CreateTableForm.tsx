"use client";

import { useState } from "react";
import { SEATS, type Seat } from "@/lib/protocol";
import {
  boardsInput,
  configValid,
  cycleSeat,
  inviteUrl,
  parseBoards,
  type TableConfig,
} from "@/lib/tableConfig";

const SEAT_LETTER: Record<Seat, string> = { North: "N", East: "E", South: "S", West: "W" };
const ROLE_LABEL = { you: "You", bot: "Bot", open: "Open" } as const;

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

export default function CreateTableForm({
  config,
  onChange,
  tableId,
  onCreate,
}: {
  config: TableConfig;
  onChange: (c: TableConfig) => void;
  tableId: string | null;
  onCreate: () => void;
}) {
  const [copied, setCopied] = useState(false);
  const link =
    tableId && typeof window !== "undefined"
      ? inviteUrl(window.location.origin, tableId)
      : "";
  const valid = configValid(config);

  function copy() {
    if (!link) return;
    navigator.clipboard?.writeText(link).then(
      () => {
        setCopied(true);
        setTimeout(() => setCopied(false), 1500);
      },
      () => {},
    );
  }

  return (
    <div className="box create-box">
      <div className="bt">
        <span className="t">Create a Table</span>
      </div>

      <div className="create-form">
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

        <div className="cf-field">
          <span>Seats</span>
          <div className="cf-seats">
            {SEATS.map((seat) => (
              <button
                key={seat}
                type="button"
                className={`cf-seat role-${config.seats[seat]}`}
                onClick={() => onChange({ ...config, seats: cycleSeat(config.seats, seat) })}
              >
                <b>{SEAT_LETTER[seat]}</b>
                <span>{ROLE_LABEL[config.seats[seat]]}</span>
              </button>
            ))}
          </div>
          {!valid && <p className="cf-hint">Pick one seat for yourself.</p>}
        </div>

        <div className="cf-field">
          <span>Invite link</span>
          <div className="cf-link">
            <input className="cf-input" readOnly value={link || "generating…"} />
            <button type="button" className="btn xs" disabled={!link} onClick={copy}>
              {copied ? "Copied" : "Copy"}
            </button>
          </div>
        </div>

        <button
          type="button"
          className="btn pri cf-create"
          disabled={!valid || !tableId}
          onClick={onCreate}
        >
          Create Table
        </button>
      </div>
    </div>
  );
}
