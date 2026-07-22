"use client";

import Button from "@/components/Button";
import { useAuth } from "@/hooks/useAuth";
import styles from "./FriendsPanel.module.css";

/**
 * Right-hand sidebar on the home page. There's no friends backend yet,
 * so for now it only handles the guest case: a prompt to log in.
 */
export default function FriendsPanel() {
  const { user, loading } = useAuth();

  return (
    <aside className={styles.fr}>
      <div className={styles.head}>
        <span className="micro">friends</span>
      </div>

      {!loading && !user && (
        <div className={styles.guest}>
          <div className={styles.msg}>
            <b>Guests can&apos;t add friends.</b>
            <p>Log in to connect with other players.</p>
          </div>
          <div className={styles.cta}>
            <Button href="/login" size="sm">
              Log in
            </Button>
            <Button href="/signup" size="sm" variant="primary">
              Sign up
            </Button>
          </div>
        </div>
      )}
    </aside>
  );
}
