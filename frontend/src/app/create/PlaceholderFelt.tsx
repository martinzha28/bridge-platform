import { cx } from "@/lib/cx";
import { SEATS } from "@/lib/protocol";
import felt from "@/components/felt.module.css";

const SLOT = { North: "n", East: "e", South: "s", West: "w" } as const;

function Backs({ count }: { count: number }) {
  return (
    <>
      {Array.from({ length: count }).map((_, i) => (
        <div key={i} className={felt.back} />
      ))}
    </>
  );
}

/** Static face-down felt shown on /create — a preview, nothing interactive. */
export default function PlaceholderFelt() {
  return (
    <div className={felt.feltBox}>
      <div className={felt.felt}>
        {SEATS.map((seat) => (
          <span key={seat} className={cx(felt.mk, felt[SLOT[seat]])}>
            {seat[0]}
          </span>
        ))}
        {SEATS.map((seat) => (
          <div key={seat} className={cx(felt.seat, felt[SLOT[seat]])}>
            <Backs count={13} />
          </div>
        ))}
      </div>
    </div>
  );
}
