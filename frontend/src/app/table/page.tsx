"use client";

import { useState } from "react";
import Rail from "@/components/Rail";
import BiddingBox from "@/components/BiddingBox";
import Card from "@/components/Card";
import { useTable, type TableStatus } from "@/hooks/useTable";
import type { PlayerView, Seat } from "@/lib/protocol";
import { SEATS } from "@/lib/protocol";
import {
  AUCTION_COLUMNS,
  SEAT_LETTER,
  auctionRows,
  boardResultText,
  canBid,
  cardFace,
  groupHandBySuit,
  handCount,
  historyRow,
  trumpSuit,
  type BoardResult,
} from "@/lib/view";

const FELT_POS: Record<Seat, string> = { North: "n", East: "e", South: "s", West: "w" };
const PLATE_ORDER: Seat[] = ["South", "West", "North", "East"];

function statusLabel(status: TableStatus): string {
  return status === "closed" ? "disconnected" : "connecting…";
}

function feltCaption(view: PlayerView | null): string {
  if (!view) return "";
  if (view.phase === "Auction") return `${view.turn} to call`;
  if (view.phase === "Play") return `${view.turn} to play`;
  if (view.result?.passedOut) return "passed out";
  if (view.result) {
    const s = view.result.score;
    return `${view.result.contract ?? "done"} · ${s >= 0 ? "+" : ""}${s}`;
  }
  return "complete";
}

export default function TablePage() {
  const { view, status, error, history, bid, playCard } = useTable();

  return (
    <div className="app">
      <Rail active="play" />

      <div className="center">
        <TableHeader view={view} />
        <Felt view={view} status={status} error={error} onBid={bid} onPlay={playCard} />
        <Plates view={view} />
      </div>

      <div className="side">
        <AuctionPanel view={view} />
        <HistoryPanel history={history} />
        <div className="box">
          <div className="bt">
            <span className="t">Table chat</span>
          </div>
        </div>
      </div>
    </div>
  );
}

function TableHeader({ view }: { view: PlayerView | null }) {
  return (
    <div className="box">
      <div className="thead">
        <div className="bchip">{view?.boardNumber ?? "–"}</div>
        <h1>Practice table</h1>
        <div className="sc">
          <span>vul</span> {view?.vulnerability ?? "–"} &nbsp;
          <span>contract</span> {view?.contract ?? "–"}
        </div>
      </div>
    </div>
  );
}

function Felt({
  view,
  status,
  error,
  onBid,
  onPlay,
}: {
  view: PlayerView | null;
  status: TableStatus;
  error: string | null;
  onBid: (call: string) => void;
  onPlay: (card: string) => void;
}) {
  const [dropActive, setDropActive] = useState(false);
  const legal = view?.legalCards ?? [];

  function handleDrop(e: React.DragEvent) {
    e.preventDefault();
    setDropActive(false);
    const code = e.dataTransfer.getData("text/card");
    if (code && legal.includes(code)) onPlay(code);
  }

  return (
    <div className="box felt-box">
      <div className="bt">
        <span className="t">
          Board {view?.boardNumber ?? "–"}
          {view ? ` · ${view.phase}` : ""}
        </span>
        <span className="r">{feltCaption(view)}</span>
      </div>

      <div
        className={`felt${dropActive ? " drop" : ""}`}
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
            className={`mk ${FELT_POS[seat]}${view?.turn === seat ? " on" : ""}`}
          >
            {SEAT_LETTER[seat]}
          </span>
        ))}

        {view &&
          SEATS.filter((s) => s !== "South").map((seat) => (
            <OpponentSeat key={seat} view={view} seat={seat} legal={legal} onPlay={onPlay} />
          ))}

        {view && view.currentTrick.length > 0 && (
          <div className="trick">
            {view.currentTrick.map((pc) => {
              const face = cardFace(pc.card);
              return (
                <div
                  key={pc.seat}
                  className={`pc ${FELT_POS[pc.seat]}${face.red ? " red" : ""}`}
                >
                  {face.rank}
                  <i>{face.suit}</i>
                </div>
              );
            })}
          </div>
        )}

        {view && (
          <div className="seat s">
            {groupHandBySuit(view.hand, trumpSuit(view.contract)).map((group, i) => (
              <div className="grp" key={i}>
                {group.map((code) => (
                  <Card
                    key={code}
                    code={code}
                    playable={legal.includes(code)}
                    onPlay={onPlay}
                  />
                ))}
              </div>
            ))}
          </div>
        )}

        {view && canBid(view) && (
          <BiddingBox legalCalls={view.legalCalls ?? []} onBid={onBid} />
        )}

        {!view && (
          <div className="felt-msg">{error ?? statusLabel(status)}</div>
        )}

        {view?.phase === "Complete" && (
          <div className="felt-msg result">{boardResultText(view)} · next board…</div>
        )}
      </div>
    </div>
  );
}

function OpponentSeat({
  view,
  seat,
  legal,
  onPlay,
}: {
  view: PlayerView;
  seat: Seat;
  legal: string[];
  onPlay: (card: string) => void;
}) {
  const pos = FELT_POS[seat];

  if (view.dummy === seat && view.dummyHand) {
    return (
      <div className={`seat ${pos} dummy`}>
        {groupHandBySuit(view.dummyHand, trumpSuit(view.contract)).map((group, i) => (
          <div className="grp" key={i}>
            {group.map((code) => (
              <Card
                key={code}
                code={code}
                playable={legal.includes(code)}
                onPlay={onPlay}
              />
            ))}
          </div>
        ))}
      </div>
    );
  }

  return (
    <div className={`seat ${pos}`}>
      {Array.from({ length: handCount(view, seat) }).map((_, i) => (
        <Card key={i} />
      ))}
    </div>
  );
}

function Plates({ view }: { view: PlayerView | null }) {
  return (
    <div className="plates">
      {PLATE_ORDER.map((seat) => {
        const me = seat === "South";
        const label = me ? "You" : `Bot ${SEAT_LETTER[seat]}`;
        const sub = view?.dummy === seat ? "dummy" : me ? "south" : "bot · easy";
        return (
          <div
            key={seat}
            className={`plate${me ? " me" : ""}${view?.turn === seat ? " turn" : ""}`}
          >
            <span className="st">{SEAT_LETTER[seat]}</span>
            <div>
              <b>{label}</b>
              <div className="ck">{sub}</div>
            </div>
          </div>
        );
      })}
    </div>
  );
}

function AuctionPanel({ view }: { view: PlayerView | null }) {
  const rows = view ? auctionRows(view.calls, view.dealer) : [];

  return (
    <div className="box">
      <div className="bt">
        <span className="t">Auction</span>
        {view?.contract && <span className="r">{view.contract}</span>}
      </div>
      <div className="scrolly">
        <table className="auc">
          <thead>
            <tr>
              {AUCTION_COLUMNS.map((seat) => (
                <th key={seat}>{SEAT_LETTER[seat]}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rows.map((row, i) => (
              <tr key={i}>
                {row.map((call, j) => (
                  <td key={j} className={/[HD]/.test(call) ? "red" : undefined}>
                    {call || " "}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

function HistoryPanel({ history }: { history: BoardResult[] }) {
  const rows = history.map(historyRow);
  return (
    <div className="box">
      <div className="bt">
        <span className="t">History</span>
        {rows.length > 0 && <span className="r">{rows.length} played</span>}
      </div>
      <div className="scrolly">
        <table className="auc">
          <thead>
            <tr>
              <th>#</th>
              <th>Result</th>
              <th>Score NS</th>
              <th>Score EW</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((r) => (
              <tr key={r.board}>
                <td>{r.board}</td>
                <td className={r.red ? "red" : undefined}>{r.result}</td>
                <td>{r.scoreNS}</td>
                <td>{r.scoreEW}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
