import { cardFace } from "@/lib/view";
import { cx } from "@/lib/cx";
import styles from "./felt.module.css";

interface CardProps {
  /** Card code like "SA"; omit for a face-down back. */
  code?: string;
  /** When true the card can be dragged or clicked to play. */
  playable?: boolean;
  onPlay?: (code: string) => void;
}

export default function Card({ code, playable = false, onPlay }: CardProps) {
  if (!code) return <div className={styles.back} />;

  const face = cardFace(code);

  return (
    <div
      className={cx(styles.face, face.red && styles.red, playable && styles.playable)}
      draggable={playable}
      onDragStart={
        playable
          ? (e) => e.dataTransfer.setData("text/card", code)
          : undefined
      }
      onClick={playable && onPlay ? () => onPlay(code) : undefined}
    >
      {face.rank}
      <i>{face.suit}</i>
    </div>
  );
}
