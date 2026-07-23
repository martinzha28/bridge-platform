import { cx } from "@/lib/cx";
import {
  BID_LEVELS,
  BID_STRAINS,
  hasCall,
  strainIsRed,
  strainLabel,
} from "@/lib/view";
import styles from "./BiddingBox.module.css";

interface BiddingBoxProps {
  legalCalls: string[];
  onBid: (call: string) => void;
}

/**
 * Overlay shown on the felt while it is this player's turn to call.
 * Every bid is rendered; only the legal ones are enabled.
 */
export default function BiddingBox({ legalCalls, onBid }: BiddingBoxProps) {
  return (
    <div className={styles.box} role="group" aria-label="Bidding box">
      <div className={styles.grid}>
        {BID_LEVELS.map((level) => (
          <div className={styles.row} key={level}>
            {BID_STRAINS.map((strain) => {
              const call = `${level}${strain}`;
              return (
                <button
                  key={strain}
                  type="button"
                  className={cx(styles.btn, strainIsRed(strain) && styles.red)}
                  disabled={!hasCall(legalCalls, call)}
                  onClick={() => onBid(call)}
                >
                  {level}
                  <i>{strainLabel(strain)}</i>
                </button>
              );
            })}
          </div>
        ))}
      </div>
      <div className={styles.row}>
        <button
          type="button"
          className={cx(styles.btn, styles.wide)}
          disabled={!hasCall(legalCalls, "P")}
          onClick={() => onBid("P")}
        >
          Pass
        </button>
        <button
          type="button"
          className={cx(styles.btn, styles.wide)}
          disabled={!hasCall(legalCalls, "X")}
          onClick={() => onBid("X")}
        >
          Dbl
        </button>
        <button
          type="button"
          className={cx(styles.btn, styles.wide)}
          disabled={!hasCall(legalCalls, "XX")}
          onClick={() => onBid("XX")}
        >
          Rdbl
        </button>
      </div>
    </div>
  );
}
