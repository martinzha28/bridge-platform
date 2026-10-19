"use client";

import { useEffect, useMemo, useState, type ReactNode } from "react";
import { useParams } from "next/navigation";
import Rail from "@/components/Rail";
import Button from "@/components/Button";
import ChatPanel from "@/components/ChatPanel";
import NamePrompt from "@/components/NamePrompt";
import Lobby from "./Lobby";
import Felt from "./Felt";
import TableHeader from "./TableHeader";
import Plates from "./Plates";
import AuctionPanel from "./AuctionPanel";
import HistoryPanel from "./HistoryPanel";
import { useTable } from "@/hooks/useTable";
import { useDisplayIdentity } from "@/hooks/useDisplayIdentity";
import type { JoinIdentity } from "@/lib/protocol";
import { DEFAULT_CONFIG, loadConfig } from "@/lib/tableConfig";
import feltStyles from "@/components/felt.module.css";

/** Full-screen "empty felt" used for join errors and connecting states. */
function FeltGate({ children }: { children: ReactNode }) {
  return (
    <div className="app">
      <Rail active="play" />
      <div className="center">
        <div className={feltStyles.feltBox}>
          <div className={feltStyles.felt}>
            <div className={feltStyles.msg}>{children}</div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default function TablePage() {
  const params = useParams<{ id: string }>();
  const identity = useDisplayIdentity();

  if (!identity.ready) return <FeltGate>Joining table…</FeltGate>;
  if (identity.needsPrompt) {
    return (
      <FeltGate>
        <NamePrompt initial={identity.name} onSubmit={identity.chooseName} />
      </FeltGate>
    );
  }
  return <TableRoom id={params.id} join={identity.join} />;
}

/** Mounted once we know who's joining — useTable connects on mount. */
function TableRoom({ id, join }: { id: string; join: JoinIdentity }) {
  // config comes from /create via sessionStorage; a bare visitor (invite
  // link) gets defaults. Who hosts is the server's call (useTable).
  const { config, isCreator } = useMemo(() => {
    const stored = loadConfig(id);
    return { config: stored ?? DEFAULT_CONFIG, isCreator: stored !== null };
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
    isHost,
    me,
    people,
    messages,
    bid,
    playCard,
    sitAt,
    stand,
    addBot,
    removeBot,
    moveTo,
    setName,
    sendChat,
    startGame,
  } = useTable({ tableId: id, config, isCreator, join });

  const tableName = tableState?.name ?? "Practice table";

  const [replayOpen, setReplayOpen] = useState(false);
  const canReplay = lastCompletedTrick.length > 0;
  useEffect(() => {
    if (!canReplay) setReplayOpen(false);
  }, [canReplay]);

  if (joinError) {
    return (
      <FeltGate>
        {joinError}
        <Button href="/create" size="sm" style={{ marginLeft: 10 }}>
          New table
        </Button>
      </FeltGate>
    );
  }

  // still opening the socket / joining — no lobby or game state yet
  if (!view && tableState == null) {
    return <FeltGate>{error ?? "Joining table…"}</FeltGate>;
  }

  // a visitor who arrives after the game started (and never took a seat)
  // can't join — the lobby is gone
  if (!view && tableState?.started && mySeat == null && !isCreator) {
    return (
      <FeltGate>
        This table is already in progress.
        <Button href="/create" size="sm" style={{ marginLeft: 10 }}>
          Make your own table
        </Button>
      </FeltGate>
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
          people={people}
          meId={me?.id ?? null}
          messages={messages}
          onSit={sitAt}
          onSitHere={moveTo}
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
              onRename={isHost ? setName : undefined}
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
            <Plates view={view} people={people} />
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
