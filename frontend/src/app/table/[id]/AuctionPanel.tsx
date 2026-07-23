import Panel from "@/components/Panel";
import { cx } from "@/lib/cx";
import {
  SEAT_LETTER,
  auctionColumns,
  auctionRows,
  seatVulnerable,
} from "@/lib/view";
import type { PlayerView } from "@/lib/protocol";
import styles from "./ladder.module.css";

export default function AuctionPanel({ view }: { view: PlayerView | null }) {
  const columns = auctionColumns(view?.seat ?? "South");
  const rows = view ? auctionRows(view.calls, view.dealer, columns) : [];

  return (
    <Panel
      title="Auction"
      aside={view?.contract}
      className={styles.auctionPanel}
    >
      <div className={styles.scroll}>
        <table className={styles.table}>
          <thead>
            <tr>
              {columns.map((seat) => (
                <th
                  key={seat}
                  style={{
                    color:
                      view && seatVulnerable(seat, view.vulnerability)
                        ? "var(--red)"
                        : "var(--vul-safe)",
                  }}
                >
                  {SEAT_LETTER[seat]}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rows.map((row, i) => (
              <tr key={i}>
                {row.map((call, j) => (
                  <td key={j} className={cx(/[HD]/.test(call) && styles.red)}>
                    {call || " "}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </Panel>
  );
}
