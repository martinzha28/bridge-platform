import { cx } from "@/lib/cx";
import { SEAT_LETTER } from "@/lib/view";
import { plateOrder } from "@/lib/table-view";
import type { PlayerView, Seat } from "@/lib/protocol";
import styles from "./Plates.module.css";

export default function Plates({ view }: { view: PlayerView | null }) {
  const mySeat: Seat = view?.seat ?? "South";
  return (
    <div className={styles.plates}>
      {plateOrder(mySeat).map((seat) => {
        const me = seat === mySeat;
        const label = me ? "You" : `Bot ${SEAT_LETTER[seat]}`;
        const sub = view?.dummy === seat ? "dummy" : me ? "you" : "bot · easy";
        return (
          <div
            key={seat}
            className={cx(
              styles.plate,
              me && styles.me,
              view?.turn === seat && styles.turn,
            )}
          >
            <span className={styles.st}>{SEAT_LETTER[seat]}</span>
            <div>
              <b>{label}</b>
              <div className={styles.ck}>{sub}</div>
            </div>
          </div>
        );
      })}
    </div>
  );
}
