import type { ButtonHTMLAttributes, ReactNode } from "react";
import Link from "next/link";
import { cx } from "@/lib/cx";
import styles from "./Button.module.css";

type Props = Omit<ButtonHTMLAttributes<HTMLButtonElement>, "className"> & {
  variant?: "default" | "primary";
  size?: "md" | "sm";
  className?: string;
  children?: ReactNode;
  /** Render as a Next.js <Link> that looks like a button. */
  href?: string;
};

export default function Button({
  variant = "default",
  size = "md",
  className,
  href,
  type,
  children,
  ...rest
}: Props) {
  const cls = cx(
    styles.btn,
    variant === "primary" && styles.primary,
    size === "sm" && styles.sm,
    className,
  );

  if (href != null) {
    return (
      <Link href={href} className={cls}>
        {children}
      </Link>
    );
  }

  return (
    <button type={type ?? "button"} className={cls} {...rest}>
      {children}
    </button>
  );
}
