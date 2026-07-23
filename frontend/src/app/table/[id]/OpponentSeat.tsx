import Card from "@/components/Card";
import { cx } from "@/lib/cx";
import { groupHandBySuit, handCount, trumpSuit } from "@/lib/view";
import { screenSlot } from "@/lib/table-view";
import type { PlayerView, Seat } from "@/lib/protocol";
import styles from "@/components/felt.module.css";

export default function OpponentSeat({
  view,
  seat,
  mySeat,
  legal,
  onPlay,
}: {
  view: PlayerView;
  seat: Seat;
  mySeat: Seat;
  legal: string[];
  onPlay: (card: string) => void;
}) {
  const pos = screenSlot(seat, mySeat);

  if (view.dummy === seat && view.dummyHand) {
    return (
      <div className={cx(styles.seat, styles[pos], styles.dummy)}>
        {groupHandBySuit(view.dummyHand, trumpSuit(view.contract)).map(
          (group, i) => (
            <div className={styles.grp} key={i}>
              {group.map((code) => (
                <Card
                  key={code}
                  code={code}
                  playable={legal.includes(code)}
                  onPlay={onPlay}
                />
              ))}
            </div>
          ),
        )}
      </div>
    );
  }

  return (
    <div className={cx(styles.seat, styles[pos])}>
      {Array.from({ length: handCount(view, seat) }).map((_, i) => (
        <Card key={i} />
      ))}
    </div>
  );
}
