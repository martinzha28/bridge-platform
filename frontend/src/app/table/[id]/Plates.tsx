import { cx } from "@/lib/cx";
import { SEAT_LETTER } from "@/lib/view";
import { peopleBySeat, plateOrder } from "@/lib/table-view";
import type { Person, PlayerView, Seat } from "@/lib/protocol";
import styles from "./Plates.module.css";

export default function Plates({
  view,
  people,
}: {
  view: PlayerView | null;
  people: Person[];
}) {
  const mySeat: Seat = view?.seat ?? "South";
  const seated = peopleBySeat(people);
  return (
    <div className={styles.plates}>
      {plateOrder(mySeat).map((seat) => {
        const me = seat === mySeat;
        // a seat with no connected person is a bot (seats can't sit empty mid-game)
        const human = seated[seat];
        const label = me ? "You" : human ? human.name : "Bot";
        const sub =
          view?.dummy === seat
            ? "dummy"
            : me
              ? "you"
              : human
                ? human.guest
                  ? "guest"
                  : "player"
                : "bot · easy";
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
