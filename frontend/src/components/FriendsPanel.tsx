"use client";

import Link from "next/link";
import { useAuth } from "@/hooks/useAuth";

/**
 * Right-hand sidebar on the home page. There's no friends backend yet,
 * so for now it only handles the guest case: a prompt to log in.
 */
export default function FriendsPanel() {
  const { user, loading } = useAuth();

  return (
    <aside className="fr">
      <div className="fr-hd">
        <span className="micro">friends</span>
      </div>

      {!loading && !user && (
        <div className="fr-guest">
          <div className="fr-msg">
            <b>Guests can&apos;t add friends.</b>
            <p>Log in to connect with other players.</p>
          </div>
          <div className="fr-cta">
            <Link href="/login" className="btn xs">
              Log in
            </Link>
            <Link href="/signup" className="btn xs pri">
              Sign up
            </Link>
          </div>
        </div>
      )}
    </aside>
  );
}
