import Panel from "@/components/Panel";
import type { Person } from "@/lib/protocol";
import { SEAT_LETTER } from "@/lib/view";
import styles from "./PeoplePanel.module.css";

/** Everyone connected to the table — players and spectators. */
export default function PeoplePanel({
  people,
  meId,
}: {
  people: Person[];
  meId: number | null;
}) {
  return (
    <Panel title="At this table" aside={people.length} padded>
      {people.map((p) => (
        <div key={p.id} className={styles.row}>
          <span className={styles.name}>
            {p.name}
            {p.host && <span className={styles.you}> · host</span>}
            {p.id === meId && <span className={styles.you}> (you)</span>}
          </span>
          <span className={styles.where}>
            {p.seat ? SEAT_LETTER[p.seat] : "watching"}
          </span>
        </div>
      ))}
    </Panel>
  );
}
