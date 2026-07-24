import { Fragment } from "react";
import Link from "next/link";
import { cx } from "@/lib/cx";
import { RailIcon, SpadeIcon } from "@/components/icons";
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
  title: string;
  href?: string;
};

// Only Home and Play route anywhere; the rest are placeholders for
// features that don't exist yet (rendered greyed out).
const NAV_GROUPS: RailItem[][] = [
  [
    { key: "home", title: "Home", href: "/" },
    { key: "play", title: "Play", href: "/create" },
    { key: "history", title: "History" },
  ],
  [
    { key: "puzzles", title: "Puzzles" },
    { key: "learn", title: "Learn" },
    { key: "watch", title: "Watch" },
  ],
  [
    { key: "news", title: "News" },
    { key: "friends", title: "Friends" },
    { key: "profile", title: "Profile" },
  ],
];

const SETTINGS_ITEM: RailItem = { key: "settings", title: "Settings" };

function RailButton({ item, active }: { item: RailItem; active: boolean }) {
  const className = cx(
    styles.rb,
    active && styles.on,
    !item.href && styles.disabled,
  );
  const icon = <RailIcon name={item.key} />;

  if (item.href) {
    return (
      <Link href={item.href} className={className} title={item.title}>
        {icon}
      </Link>
    );
  }
  return (
    <div className={className} title={`${item.title} — coming soon`}>
      {icon}
    </div>
  );
}

export default function Rail({ active }: { active?: RailKey }) {
  return (
    <nav className={styles.rail}>
      <div className={styles.lg} title="Bridge++">
        <SpadeIcon width={22} height={22} />
      </div>
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
