import type { PlayedCard, Seat } from "@/lib/protocol";
import TrickCards from "./TrickCards";
import styles from "@/components/felt.module.css";

export default function ReplayPopup({
  cards,
  mySeat,
  onClose,
}: {
  cards: PlayedCard[];
  mySeat: Seat;
  onClose: () => void;
}) {
  return (
    <>
      <div className={styles.replayBackdrop} onClick={onClose} />
      <div className={styles.replayBox}>
        <TrickCards cards={cards} mySeat={mySeat} />
      </div>
    </>
  );
}
