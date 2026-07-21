import type { ButtonHTMLAttributes } from "react";
import { cx } from "@/lib/cx";
import styles from "./Button.module.css";

type Props = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "default" | "primary";
  size?: "md" | "sm";
};

export default function Button({
  variant = "default",
  size = "md",
  className,
  type = "button",
  ...rest
}: Props) {
  return (
    <button
      type={type}
      className={cx(
        styles.btn,
        variant === "primary" && styles.primary,
        size === "sm" && styles.sm,
        className,
      )}
      {...rest}
    />
  );
}
