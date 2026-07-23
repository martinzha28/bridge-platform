"use client";

import Panel from "@/components/Panel";
import Button from "@/components/Button";
import { Pill, PillGroup } from "@/components/Pill";
import { SEATS, type Seat, type TableState } from "@/lib/protocol";
import { SEAT_LETTER } from "@/lib/view";
import styles from "./SeatsPanel.module.css";

export default function SeatsPanel({
  tableState,
  mySeat,
  isHost,
  onSit,
  onSitHere,
  onTakeSeat,
  onStand,
  onAddBot,
  onRemoveBot,
  onStart,
}: {
  tableState: TableState;
  mySeat: Seat | null;
  isHost: boolean;
  onSit: (dir: Seat) => void;
  onSitHere: (dir: Seat) => void;
  onTakeSeat: (dir: Seat) => void;
  onStand: () => void;
  onAddBot: (dir: Seat) => void;
  onRemoveBot: (dir: Seat) => void;
  onStart: () => void;
}) {
  const open = SEATS.filter((s) => tableState.seats[s] === "");
  const filled = open.length === 0;
  const startHint = filled
    ? "All seats filled — ready when you are."
    : `Waiting on ${open.map((s) => SEAT_LETTER[s]).join(", ")} — add a bot or a player.`;

  return (
    <Panel title="Seats" padded>
      {SEATS.map((seat) => {
        const who = tableState.seats[seat];
        const isMine = mySeat === seat;
        return (
          <div key={seat} className={styles.row}>
            <span className={styles.dir}>{SEAT_LETTER[seat]}</span>
            {isMine ? (
              <>
                <span className={styles.tag}>You</span>
                {!isHost && (
                  <Button
                    size="sm"
                    className={styles.pushRight}
                    onClick={onStand}
                  >
                    Stand up
                  </Button>
                )}
              </>
            ) : who === "human" ? (
              <span className={styles.tag}>Player</span>
            ) : isHost ? (
              <>
                <PillGroup className={styles.pills}>
                  <Pill
                    compact
                    active={who === "bot"}
                    onClick={() => who !== "bot" && onAddBot(seat)}
                  >
                    Bot
                  </Pill>
                  <Pill
                    compact
                    active={who === ""}
                    onClick={() => who !== "" && onRemoveBot(seat)}
                  >
                    Open
                  </Pill>
                </PillGroup>
                <Button
                  size="sm"
                  className={styles.pushRight}
                  onClick={() => onSitHere(seat)}
                >
                  Sit here
                </Button>
              </>
            ) : who === "bot" ? (
              <Button
                size="sm"
                className={styles.pushRight}
                onClick={() => onTakeSeat(seat)}
              >
                Take seat
              </Button>
            ) : (
              <Button
                size="sm"
                className={styles.pushRight}
                onClick={() => onSit(seat)}
              >
                Sit here
              </Button>
            )}
          </div>
        );
      })}

      <Button
        variant="primary"
        className={styles.start}
        disabled={!filled}
        onClick={onStart}
      >
        Start game
      </Button>
      <p className={styles.hint}>{startHint}</p>
    </Panel>
  );
}
