import { Fragment } from "react";
import Link from "next/link";
import { cx } from "@/lib/cx";
import styles from "./Rail.module.css";

type RailKey =
  | "home"
  | "play"
  | "history"
  | "puzzles"
  | "learn"
  | "watch"
  | "news"
  | "friends"
  | "profile"
  | "settings";

type RailItem = {
  key: RailKey;
  label: string;
  title: string;
  href?: string;
};

// Only Home and Play route anywhere; the rest are placeholders for
// features that don't exist yet.
const NAV_GROUPS: RailItem[][] = [
  [
    { key: "home", label: "Hm", title: "Home", href: "/" },
    { key: "play", label: "Pl", title: "Play", href: "/create" },
    { key: "history", label: "Hi", title: "History" },
  ],
  [
    { key: "puzzles", label: "Pz", title: "Puzzles" },
    { key: "learn", label: "Ln", title: "Learn" },
    { key: "watch", label: "Wa", title: "Watch" },
  ],
  [
    { key: "news", label: "Nw", title: "News" },
    { key: "friends", label: "Fr", title: "Friends" },
    { key: "profile", label: "Me", title: "Profile" },
  ],
];

const SETTINGS_ITEM: RailItem = { key: "settings", label: "Cf", title: "Settings" };

function RailButton({ item, active }: { item: RailItem; active: boolean }) {
  const className = cx(styles.rb, active && styles.on);
  if (item.href) {
    return (
      <Link href={item.href} className={className} title={item.title}>
        {item.label}
      </Link>
    );
  }
  return (
    <div className={className} title={item.title}>
      {item.label}
    </div>
  );
}

export default function Rail({ active }: { active?: RailKey }) {
  return (
    <nav className={styles.rail}>
      <div className={styles.lg}>&spades;</div>
      {NAV_GROUPS.map((group, i) => (
        <Fragment key={i}>
          {i > 0 && <div className={styles.sep} />}
          {group.map((item) => (
            <RailButton key={item.key} item={item} active={active === item.key} />
          ))}
        </Fragment>
      ))}
      <div className={styles.spacer} />
      <RailButton item={SETTINGS_ITEM} active={active === "settings"} />
    </nav>
  );
}
