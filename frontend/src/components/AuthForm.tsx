"use client";

import { useState, type FormEvent } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import Button from "@/components/Button";
import { login, signup } from "@/lib/auth";
import styles from "./AuthForm.module.css";

const CONFIG = {
  login: {
    title: "Log in",
    submit: "Log in",
    action: login,
    helper: null as string | null,
    passwordAutoComplete: "current-password",
    alt: { prompt: "Need an account?", href: "/signup", label: "Sign up" },
  },
  signup: {
    title: "Sign up",
    submit: "Create account",
    action: signup,
    helper:
      "Password: at least 8 characters, one uppercase letter, and one special character.",
    passwordAutoComplete: "new-password",
    alt: { prompt: "Already have an account?", href: "/login", label: "Log in" },
  },
} as const;

export default function AuthForm({ mode }: { mode: "login" | "signup" }) {
  const cfg = CONFIG[mode];
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");
    setBusy(true);
    const res = await cfg.action(username.trim(), password);
    if (res.ok) {
      router.push("/");
    } else {
      setError(res.error);
      setBusy(false);
    }
  }

  return (
    <div className={styles.auth}>
      <form className={styles.card} onSubmit={onSubmit}>
        <h1>{cfg.title}</h1>

        <label className={styles.field}>
          <span>Username</span>
          <input
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            autoComplete="username"
            required
          />
        </label>

        <label className={styles.field}>
          <span>Password</span>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            autoComplete={cfg.passwordAutoComplete}
            required
          />
        </label>

        {cfg.helper && <p className={styles.helper}>{cfg.helper}</p>}
        {error && <p className={styles.error}>{error}</p>}

        <Button
          type="submit"
          variant="primary"
          className={styles.submit}
          disabled={busy}
        >
          {busy ? "…" : cfg.submit}
        </Button>

        <p className={styles.alt}>
          {cfg.alt.prompt} <Link href={cfg.alt.href}>{cfg.alt.label}</Link>
        </p>
        <Link href="/" className={styles.back}>
          &larr; Back to Bridge++
        </Link>
      </form>
    </div>
  );
}
