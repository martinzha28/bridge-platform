import Panel from "@/components/Panel";
import Button from "@/components/Button";
import styles from "./PlayModule.module.css";

/** Home-page "Play" panel: the one call to action. */
export default function PlayModule() {
  return (
    <Panel title="Play">
      <div className={styles.play}>
        <div className={styles.lead}>
          <h2>
            Sit down and play.
            <br />
            &spades; <span className="red">&hearts;</span>{" "}
            <span className="red">&diams;</span> &clubs;
          </h2>
          <p>
            Casual or rated, IMPs or matchpoints, humans or bots. Take any open
            seat and deal.
          </p>
          <div className={styles.cta}>
            <Button href="/create" variant="primary">
              Create Table
            </Button>
          </div>
        </div>
      </div>
    </Panel>
  );
}
