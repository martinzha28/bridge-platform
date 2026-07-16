import { cardFace } from "@/lib/view";

interface CardProps {
  /** Card code like "SA"; omit for a face-down back. */
  code?: string;
  /** When true the card can be dragged or clicked to play. */
  playable?: boolean;
  onPlay?: (code: string) => void;
}

export default function Card({ code, playable = false, onPlay }: CardProps) {
  if (!code) return <div className="cb" />;

  const face = cardFace(code);
  const className = ["cf", face.red && "red", playable && "playable"]
    .filter(Boolean)
    .join(" ");

  return (
    <div
      className={className}
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
