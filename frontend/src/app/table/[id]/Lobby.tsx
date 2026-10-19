"use client";

import Panel from "@/components/Panel";
import Button from "@/components/Button";
import Card from "@/components/Card";
import ChatPanel from "@/components/ChatPanel";
import {
  SEATS,
  type ChatMessage,
  type Person,
  type Seat,
  type TableState,
} from "@/lib/protocol";
import { SEAT_LETTER } from "@/lib/view";
import { peopleBySeat, plateOrder, screenSlot } from "@/lib/table-view";
import { cx } from "@/lib/cx";
import EditableName from "./EditableName";
import SeatsPanel from "./SeatsPanel";
import InvitePanel from "./InvitePanel";
import PeoplePanel from "./PeoplePanel";
import styles from "./Lobby.module.css";
import felt from "@/components/felt.module.css";
import plates from "./Plates.module.css";

export default function Lobby({
  tableId,
  tableState,
  name,
  mySeat,
  isHost,
  people,
  meId,
  messages,
  onSit,
  onSitHere,
  onStand,
  onAddBot,
  onRemoveBot,
  onStart,
  onRename,
  onChat,
}: {
  tableId: string;
  tableState: TableState;
  name: string;
  mySeat: Seat | null;
  isHost: boolean;
  people: Person[];
  meId: number | null;
  messages: ChatMessage[];
  onSit: (dir: Seat) => void;
  onSitHere: (dir: Seat) => void;
  onStand: () => void;
  onAddBot: (dir: Seat) => void;
  onRemoveBot: (dir: Seat) => void;
  onStart: () => void;
  onRename: (v: string) => void;
  onChat: (text: string) => void;
}) {
  const feltSeat = mySeat ?? "South";
  const seated = peopleBySeat(people);

  function occupant(seat: Seat): string {
    if (mySeat === seat) return "You";
    const who = tableState.seats[seat];
    if (who === "bot") return "Bot";
    if (who === "human") return seated[seat]?.name ?? "Player";
    return "Open";
  }

  function sub(seat: Seat): string {
    if (tableState.seats[seat] === "") return "waiting";
    return seated[seat]?.guest ? "guest" : "";
  }

  return (
    <>
      <div className="center">
        <Panel>
          <div className={styles.header}>
            <div className={styles.chip}>–</div>
            <EditableName name={name} onRename={isHost ? onRename : undefined} />
            <Button href="/" size="sm">
              Leave
            </Button>
          </div>
        </Panel>

        <div className={felt.feltBox}>
          <div className={felt.felt}>
            {SEATS.map((seat) => (
              <span
                key={seat}
                className={cx(felt.mk, felt[screenSlot(seat, feltSeat)])}
              >
                {SEAT_LETTER[seat]}
              </span>
            ))}
            {SEATS.map((seat) => (
              <div
                key={seat}
                className={cx(felt.seat, felt[screenSlot(seat, feltSeat)])}
              >
                {Array.from({ length: 13 }).map((_, i) => (
                  <Card key={i} />
                ))}
              </div>
            ))}
          </div>
        </div>

        <div className={plates.plates}>
          {plateOrder(feltSeat).map((seat) => (
            <div
              key={seat}
              className={cx(plates.plate, mySeat === seat && plates.me)}
            >
              <span className={plates.st}>{SEAT_LETTER[seat]}</span>
              <div>
                <b>{occupant(seat)}</b>
                <div className={plates.ck}>{sub(seat)}</div>
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className="side">
        <SeatsPanel
          tableState={tableState}
          mySeat={mySeat}
          isHost={isHost}
          names={seated}
          onSit={onSit}
          onSitHere={onSitHere}
          hostName={people.find((p) => p.host)?.name ?? null}
          onStand={onStand}
          onAddBot={onAddBot}
          onRemoveBot={onRemoveBot}
          onStart={onStart}
        />
        <InvitePanel tableId={tableId} description={tableState.description} />
        <PeoplePanel people={people} meId={meId} />
        <ChatPanel messages={messages} onSend={onChat} />
      </div>
    </>
  );
}
