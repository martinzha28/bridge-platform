import { cx } from "@/lib/cx";
import { SEATS } from "@/lib/protocol";
import styles from "./PlaceholderFelt.module.css";

const SLOT = { North: "n", East: "e", South: "s", West: "w" } as const;

function Backs({ count }: { count: number }) {
  return (
    <>
      {Array.from({ length: count }).map((_, i) => (
        <div key={i} className={styles.cb} />
      ))}
    </>
  );
}

/** Static face-down felt shown on /create — a preview, nothing interactive. */
export default function PlaceholderFelt() {
  return (
    <div className={styles.box}>
      <div className={styles.felt}>
        {SEATS.map((seat) => (
          <span key={seat} className={cx(styles.mk, styles[SLOT[seat]])}>
            {seat[0]}
          </span>
        ))}
        {SEATS.map((seat) => (
          <div key={seat} className={cx(styles.seat, styles[SLOT[seat]])}>
            <Backs count={13} />
          </div>
        ))}
      </div>
    </div>
  );
}
