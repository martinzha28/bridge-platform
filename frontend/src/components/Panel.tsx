import type { ReactNode } from "react";
import { cx } from "@/lib/cx";
import styles from "./Panel.module.css";

/** Hairline-bordered panel with an optional uppercase title bar. */
export default function Panel({
  title,
  aside,
  padded,
  className,
  children,
}: {
  title?: ReactNode;
  aside?: ReactNode;
  /** Wrap children in a padded, scrollable column. */
  padded?: boolean;
  className?: string;
  children?: ReactNode;
}) {
  return (
    <div className={cx(styles.box, className)}>
      {title != null && (
        <div className={styles.bar}>
          <span className={styles.title}>{title}</span>
          {aside != null && <span className={styles.aside}>{aside}</span>}
        </div>
      )}
      {padded ? <div className={styles.body}>{children}</div> : children}
    </div>
  );
}
