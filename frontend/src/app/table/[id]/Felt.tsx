"use client";

import { useState } from "react";
import BiddingBox from "@/components/BiddingBox";
import Card from "@/components/Card";
import { cx } from "@/lib/cx";
import {
  SEAT_LETTER,
  boardResultText,
  canBid,
  groupHandBySuit,
  trumpSuit,
} from "@/lib/view";
import { screenSlot, statusLabel, vulEdges } from "@/lib/table-view";
import type { TableStatus } from "@/hooks/useTable";
import { SEATS, type PlayedCard, type PlayerView, type Seat } from "@/lib/protocol";
import OpponentSeat from "./OpponentSeat";
import TrickCards from "./TrickCards";
import ReplayPopup from "./ReplayPopup";
import styles from "@/components/felt.module.css";

export default function Felt({
  view,
  status,
  error,
  paused,
  replayCards,
  onCloseReplay,
  onBid,
  onPlay,
}: {
  view: PlayerView | null;
  status: TableStatus;
  error: string | null;
  paused: boolean;
  replayCards: PlayedCard[] | null;
  onCloseReplay: () => void;
  onBid: (call: string) => void;
  onPlay: (card: string) => void;
}) {
  const [dropActive, setDropActive] = useState(false);
  // nothing is playable while a finished trick is held on the table
  const legal = paused ? [] : (view?.legalCards ?? []);
  const { axis, cross } = vulEdges(view);
  const mySeat: Seat = view?.seat ?? "South";
  // the live trick, or the just-finished one during the auto-pause
  const trick =
    view && view.currentTrick.length > 0
      ? view.currentTrick
      : (view?.lastTrick ?? []);

  function handleDrop(e: React.DragEvent) {
    e.preventDefault();
    setDropActive(false);
    const code = e.dataTransfer.getData("text/card");
    if (code && legal.includes(code)) onPlay(code);
  }

  return (
    <div className={styles.feltBox}>
      <div
        className={cx(styles.felt, dropActive && styles.drop)}
        style={{
          borderTopColor: axis,
          borderBottomColor: axis,
          borderLeftColor: cross,
          borderRightColor: cross,
        }}
        onDragOver={(e) => {
          e.preventDefault();
          setDropActive(true);
        }}
        onDragLeave={() => setDropActive(false)}
        onDrop={handleDrop}
      >
        {SEATS.map((seat) => (
          <span
            key={seat}
            className={cx(
              styles.mk,
              styles[screenSlot(seat, mySeat)],
              view?.turn === seat && styles.on,
            )}
          >
            {SEAT_LETTER[seat]}
          </span>
        ))}

        {view &&
          SEATS.filter((s) => s !== mySeat).map((seat) => (
            <OpponentSeat
              key={seat}
              view={view}
              seat={seat}
              mySeat={mySeat}
              legal={legal}
              onPlay={onPlay}
            />
          ))}

        {trick.length > 0 && <TrickCards cards={trick} mySeat={mySeat} />}

        {view && (
          <div className={cx(styles.seat, styles.s)}>
            {groupHandBySuit(view.hand, trumpSuit(view.contract)).map(
              (group, i) => (
                <div className={styles.grp} key={i}>
                  {group.map((code) => (
                    <Card
                      key={code}
                      code={code}
                      playable={legal.includes(code)}
                      onPlay={onPlay}
                    />
                  ))}
                </div>
              ),
            )}
          </div>
        )}

        {view && canBid(view) && (
          <BiddingBox legalCalls={view.legalCalls ?? []} onBid={onBid} />
        )}

        {!view && (
          <div className={styles.msg}>{error ?? statusLabel(status)}</div>
        )}

        {view?.phase === "Complete" && (
          <div className={cx(styles.msg, styles.result)}>
            {boardResultText(view)} · next board…
          </div>
        )}

        {replayCards && replayCards.length > 0 && (
          <ReplayPopup
            cards={replayCards}
            mySeat={mySeat}
            onClose={onCloseReplay}
          />
        )}
      </div>
    </div>
  );
}
