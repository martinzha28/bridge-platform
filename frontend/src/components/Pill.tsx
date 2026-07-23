import type { ButtonHTMLAttributes, ReactNode } from "react";
import { cx } from "@/lib/cx";
import styles from "./Pill.module.css";

export function PillGroup({
  className,
  children,
}: {
  className?: string;
  children: ReactNode;
}) {
  return <div className={cx(styles.group, className)}>{children}</div>;
}

export function Pill({
  active,
  compact,
  className,
  ...rest
}: ButtonHTMLAttributes<HTMLButtonElement> & {
  active?: boolean;
  /** Shrink to content instead of filling the row. */
  compact?: boolean;
}) {
  return (
    <button
      type="button"
      className={cx(
        styles.pill,
        active && styles.on,
        compact && styles.compact,
        className,
      )}
      {...rest}
    />
  );
}
