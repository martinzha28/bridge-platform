import type { ReactNode } from "react";
import { cx } from "@/lib/cx";
import styles from "./Panel.module.css";

/** Hairline-bordered panel with an optional uppercase title bar. */
export default function Panel({
  title,
  aside,
  className,
  children,
}: {
  title?: ReactNode;
  aside?: ReactNode;
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
      {children}
    </div>
  );
}
