"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import Rail from "@/components/Rail";
import BiddingBox from "@/components/BiddingBox";
import Card from "@/components/Card";
import ChatPanel from "@/components/ChatPanel";
import Lobby from "./Lobby";
import EditableName from "./EditableName";
import { useTable, type TableStatus } from "@/hooks/useTable";
import type {
  ChatMessage,
  PlayedCard,
  PlayerView,
  Seat,
  TableState,
} from "@/lib/protocol";
import { SEATS } from "@/lib/protocol";
import { DEFAULT_CONFIG, inviteUrl, loadConfig } from "@/lib/tableConfig";
import {
  SEAT_LETTER,
  auctionColumns,
  auctionRows,
  boardResultText,
  canBid,
  cardFace,
  groupHandBySuit,
  handCount,
  historyRow,
  seatVulnerable,
  sideTricks,
  trumpSuit,
  type BoardResult,
} from "@/lib/view";
import {
  plateOrder,
  screenSlot,
  statusLabel,
  vulEdges,
} from "@/lib/table-view";

export default function TablePage() {
  const params = useParams<{ id: string }>();
  const id = params.id;

  // config + host flag come from /create via sessionStorage; a bare
  // visitor (invite link) gets defaults and is not the host.
  const { config, isHost } = useMemo(() => {
    const stored = loadConfig(id);
    return { config: stored ?? DEFAULT_CONFIG, isHost: stored !== null };
  }, [id]);

  const {
    view,
    tableState,
    status,
    error,
    joinError,
    history,
    paused,
    lastCompletedTrick,
    mySeat,
    messages,
    bid,
    playCard,
    sitAt,
    stand,
    takeSeat,
    addBot,
    removeBot,
    moveTo,
    setName,
    sendChat,
    startGame,
  } = useTable({ tableId: id, config, isHost });

  const tableName = tableState?.name ?? "Practice table";

  const [replayOpen, setReplayOpen] = useState(false);
  const canReplay = lastCompletedTrick.length > 0;
  useEffect(() => {
    if (!canReplay) setReplayOpen(false);
  }, [canReplay]);

  if (joinError) {
    return (
      <div className="app">
        <Rail active="play" />
        <div className="center">
          <div className="felt-box">
            <div className="felt">
              <div className="felt-msg">
                {joinError}
                <Link href="/create" className="btn xs" style={{ marginLeft: 10 }}>
                  New table
                </Link>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // still opening the socket / joining — no lobby or game state yet
  if (!view && tableState == null) {
    return (
      <div className="app">
        <Rail active="play" />
        <div className="center">
          <div className="felt-box">
            <div className="felt">
              <div className="felt-msg">{error ?? "Joining table…"}</div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // a visitor who arrives after the game started (and never took a seat)
  // can't join — the lobby is gone
  if (!view && tableState?.started && mySeat == null && !isHost) {
    return (
      <div className="app">
        <Rail active="play" />
        <div className="center">
          <div className="felt-box">
            <div className="felt">
              <div className="felt-msg">
                This table is already in progress.
                <Link href="/create" className="btn xs" style={{ marginLeft: 10 }}>
                  Make your own table
                </Link>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  const inLobby = !view && tableState != null && !tableState.started;

  return (
    <div className="app">
      <Rail active="play" />

      {inLobby ? (
        <Lobby
          tableId={id}
          tableState={tableState!}
          name={tableName}
          mySeat={mySeat}
          isHost={isHost}
          messages={messages}
          onSit={sitAt}
          onSitHere={moveTo}
          onTakeSeat={takeSeat}
          onStand={stand}
          onAddBot={addBot}
          onRemoveBot={removeBot}
          onStart={startGame}
          onRename={setName}
          onChat={sendChat}
        />
      ) : (
        <>
          <div className="center">
            <TableHeader
              view={view}
              name={tableName}
              onRename={setName}
              canReplay={canReplay}
              onReplay={() => setReplayOpen(true)}
            />
            <Felt
              view={view}
              status={status}
              error={error}
              paused={paused}
              replayCards={replayOpen ? lastCompletedTrick : null}
              onCloseReplay={() => setReplayOpen(false)}
              onBid={bid}
              onPlay={playCard}
            />
            <Plates view={view} />
          </div>

          <div className="side">
            <AuctionPanel view={view} />
            <HistoryPanel history={history} />
            <ChatPanel messages={messages} onSend={sendChat} />
          </div>
        </>
      )}
    </div>
  );
}

function TableHeader({
  view,
  name,
  onRename,
  canReplay,
  onReplay,
}: {
  view: PlayerView | null;
  name: string;
  onRename: (v: string) => void;
  canReplay: boolean;
  onReplay: () => void;
}) {
  const { axis, cross } = vulEdges(view);
  const { ours, theirs } = view ? sideTricks(view) : { ours: 0, theirs: 0 };
  return (
    <div className="box">
      <div className="thead">
        <div
          className="bchip"
          style={{
            borderTopColor: axis,
            borderBottomColor: axis,
            borderLeftColor: cross,
            borderRightColor: cross,
          }}
        >
          {view?.boardNumber ?? "–"}
        </div>
        <EditableName name={name} onRename={onRename} />
        <div className="tricks">
          <div className="tks tks-us" title="Tricks we've taken">
            <b>{ours}</b>
          </div>
          <div className="tks tks-them" title="Tricks they've taken">
            <b>{theirs}</b>
          </div>
          <button className="xs" disabled={!canReplay} onClick={onReplay}>
            Replay
          </button>
        </div>
        <Link href="/" className="btn xs">
          Leave
        </Link>
      </div>
    </div>
  );
}

function TrickCards({ cards, mySeat }: { cards: PlayedCard[]; mySeat: Seat }) {
  return (
    <div className="trick">
      {cards.map((pc) => {
        const face = cardFace(pc.card);
        return (
          <div
            key={pc.seat}
            className={`pc ${screenSlot(pc.seat, mySeat)}${face.red ? " red" : ""}`}
          >
            {face.rank}
            <i>{face.suit}</i>
          </div>
        );
      })}
    </div>
  );
}

function Felt({
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
    <div className="felt-box">
      <div
        className={`felt${dropActive ? " drop" : ""}`}
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
            className={`mk ${screenSlot(seat, mySeat)}${view?.turn === seat ? " on" : ""}`}
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

        {replayCards && replayCards.length > 0 && (
          <>
            <div className="replay-backdrop" onClick={onCloseReplay} />
            <div className="replay-box">
              <TrickCards cards={replayCards} mySeat={mySeat} />
            </div>
          </>
        )}
      </div>
    </div>
  );
}

function OpponentSeat({
  view,
  seat,
  mySeat,
  legal,
  onPlay,
}: {
  view: PlayerView;
  seat: Seat;
  mySeat: Seat;
  legal: string[];
  onPlay: (card: string) => void;
}) {
  const pos = screenSlot(seat, mySeat);

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
  const mySeat: Seat = view?.seat ?? "South";
  return (
    <div className="plates">
      {plateOrder(mySeat).map((seat) => {
        const me = seat === mySeat;
        const label = me ? "You" : `Bot ${SEAT_LETTER[seat]}`;
        const sub = view?.dummy === seat ? "dummy" : me ? "you" : "bot · easy";
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
  const columns = auctionColumns(view?.seat ?? "South");
  const rows = view ? auctionRows(view.calls, view.dealer, columns) : [];

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
              {columns.map((seat) => (
                <th
                  key={seat}
                  style={{
                    color:
                      view && seatVulnerable(seat, view.vulnerability)
                        ? "var(--red)"
                        : "var(--vul-safe)",
                  }}
                >
                  {SEAT_LETTER[seat]}
                </th>
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
