import Panel from "@/components/Panel";
import { cx } from "@/lib/cx";
import { historyRow, type BoardResult } from "@/lib/view";
import styles from "./ladder.module.css";

export default function HistoryPanel({ history }: { history: BoardResult[] }) {
  const rows = history.map(historyRow);
  return (
    <Panel
      title="History"
      aside={rows.length > 0 ? `${rows.length} played` : undefined}
      className={styles.historyPanel}
    >
      <div className={styles.scroll}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th>#</th>
              <th>Result</th>
              <th>Score NS</th>
              <th>Score EW</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((r) => (
              <tr key={r.board}>
                <td>{r.board}</td>
                <td className={cx(r.red && styles.red)}>{r.result}</td>
                <td>{r.scoreNS}</td>
                <td>{r.scoreEW}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </Panel>
  );
}
