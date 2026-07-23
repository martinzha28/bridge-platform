"use client";

import Panel from "@/components/Panel";
import Button from "@/components/Button";
import { Field, Input, Textarea } from "@/components/Field";
import { Pill, PillGroup } from "@/components/Pill";
import { boardsInput, parseBoards, type TableConfig } from "@/lib/tableConfig";
import styles from "./CreateTableForm.module.css";

function Choice<T extends string>({
  value,
  options,
  onChange,
}: {
  value: T;
  options: { value: T; label: string; disabled?: boolean }[];
  onChange: (v: T) => void;
}) {
  return (
    <PillGroup>
      {options.map((o) => (
        <Pill
          key={o.value}
          disabled={o.disabled}
          active={value === o.value}
          onClick={() => onChange(o.value)}
        >
          {o.label}
        </Pill>
      ))}
    </PillGroup>
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
    <Panel title="Create a New Table" padded className={styles.panel}>
      <Field label="Title">
          <Input
            maxLength={60}
            placeholder="Practice table"
            value={config.name}
            onChange={(e) => onChange({ ...config, name: e.target.value })}
          />
        </Field>

        <Field label="Description">
          <Textarea
            maxLength={280}
            rows={4}
            placeholder="Optional — what's this table for?"
            value={config.description}
            onChange={(e) => onChange({ ...config, description: e.target.value })}
          />
        </Field>

        <Field label="Mode">
          <Choice
            value={config.mode}
            onChange={(mode) => onChange({ ...config, mode })}
            options={[
              { value: "casual", label: "Casual" },
              { value: "competitive", label: "Competitive" },
            ]}
          />
        </Field>

        <Field label="Visibility">
          <Choice
            value={config.visibility}
            onChange={(visibility) => onChange({ ...config, visibility })}
            options={[
              { value: "public", label: "Public" },
              { value: "private", label: "Private" },
            ]}
          />
        </Field>

        <Field label="Scoring">
          <Choice
            value={config.scoring}
            onChange={(scoring) => onChange({ ...config, scoring })}
            options={[
              { value: "imps", label: "IMPs" },
              { value: "matchpoints", label: "Matchpoints", disabled: true },
            ]}
          />
        </Field>

        <Field label="Boards">
          <Input
            inputMode="numeric"
            placeholder="unlimited"
            value={boardsInput(config.boards)}
            onChange={(e) =>
              onChange({ ...config, boards: parseBoards(e.target.value) })
            }
          />
        </Field>

        <Button
          variant="primary"
          className={styles.submit}
          disabled={!tableId || creating}
          onClick={onContinue}
        >
          {creating ? "Creating…" : "Continue to seating →"}
        </Button>
        <p className={styles.hint}>
          {creating
            ? "Setting up your table…"
            : "Next: pick seats, add bots, and invite players."}
        </p>
    </Panel>
  );
}
