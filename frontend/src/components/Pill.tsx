import type { ButtonHTMLAttributes, ReactNode } from "react";
import { cx } from "@/lib/cx";
import styles from "./Pill.module.css";

export function PillGroup({ children }: { children: ReactNode }) {
  return <div className={styles.group}>{children}</div>;
}

export function Pill({
  active,
  className,
  ...rest
}: ButtonHTMLAttributes<HTMLButtonElement> & { active?: boolean }) {
  return (
    <button
      type="button"
      className={cx(styles.pill, active && styles.on, className)}
      {...rest}
    />
  );
}
