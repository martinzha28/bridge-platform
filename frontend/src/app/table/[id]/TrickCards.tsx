import { cx } from "@/lib/cx";
import { cardFace } from "@/lib/view";
import { screenSlot } from "@/lib/table-view";
import type { PlayedCard, Seat } from "@/lib/protocol";
import styles from "@/components/felt.module.css";

export default function TrickCards({
  cards,
  mySeat,
}: {
  cards: PlayedCard[];
  mySeat: Seat;
}) {
  return (
    <div className={styles.trick}>
      {cards.map((pc) => {
        const face = cardFace(pc.card);
        return (
          <div
            key={pc.seat}
            className={cx(
              styles.pc,
              styles[screenSlot(pc.seat, mySeat)],
              face.red && styles.red,
            )}
          >
            {face.rank}
            <i>{face.suit}</i>
          </div>
        );
      })}
    </div>
  );
}
