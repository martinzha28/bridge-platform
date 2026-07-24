import type { ReactNode } from "react";
import styles from "./WipNotice.module.css";

/** Small "more to come" disclaimer. Drop it anywhere a section is a
 *  preview of features still being built. */
export default function WipNotice({ children }: { children?: ReactNode }) {
  return (
    <div className={styles.wip}>
      <svg
        viewBox="0 0 24 24"
        width={18}
        height={18}
        fill="none"
        stroke="currentColor"
        strokeWidth={1.6}
        strokeLinecap="round"
        strokeLinejoin="round"
        aria-hidden="true"
      >
        <path d="M12 9v4M12 17h.01" />
        <path d="M10.3 3.9 1.8 18a2 2 0 0 0 1.7 3h17a2 2 0 0 0 1.7-3L13.7 3.9a2 2 0 0 0-3.4 0Z" />
      </svg>
      <p>
        {children ?? (
          <>
            <b>Work in progress.</b> Rated play, tournaments, puzzles, and friends
            are on the way.
          </>
        )}
      </p>
    </div>
  );
}
