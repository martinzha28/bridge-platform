"use client";

import Panel from "@/components/Panel";
import Button from "@/components/Button";
import { cx } from "@/lib/cx";
import { sideTricks } from "@/lib/view";
import { vulEdges } from "@/lib/table-view";
import type { PlayerView } from "@/lib/protocol";
import EditableName from "./EditableName";
import styles from "./TableHeader.module.css";

export default function TableHeader({
  view,
  name,
  onRename,
  canReplay,
  onReplay,
}: {
  view: PlayerView | null;
  name: string;
  onRename: (v: string) => void;
  canReplay: boolean;
  onReplay: () => void;
}) {
  const { axis, cross } = vulEdges(view);
  const { ours, theirs } = view ? sideTricks(view) : { ours: 0, theirs: 0 };

  return (
    <Panel>
      <div className={styles.header}>
        <div
          className={styles.chip}
          style={{
            borderTopColor: axis,
            borderBottomColor: axis,
            borderLeftColor: cross,
            borderRightColor: cross,
          }}
        >
          {view?.boardNumber ?? "–"}
        </div>
        <EditableName name={name} onRename={onRename} />
        <div className={styles.tricks}>
          <div className={cx(styles.tks, styles.us)} title="Tricks we've taken">
            <b>{ours}</b>
          </div>
          <div
            className={cx(styles.tks, styles.them)}
            title="Tricks they've taken"
          >
            <b>{theirs}</b>
          </div>
          <Button size="sm" disabled={!canReplay} onClick={onReplay}>
            Replay
          </Button>
        </div>
        <Button href="/" size="sm">
          Leave
        </Button>
      </div>
    </Panel>
  );
}
