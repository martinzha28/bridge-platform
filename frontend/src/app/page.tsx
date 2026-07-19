import type { Metadata } from "next";
import Link from "next/link";
import Rail from "@/components/Rail";

export const metadata: Metadata = {
  title: "Bridge++ — Home",
};

export default function HomePage() {
  return (
    <div className="app">
      <Rail active="home" />

      <div className="scroll">
        <div className="top">
          <h1>Bridge++ / home</h1>
        </div>

        <div className="page">
          <div className="box">
            <div className="bt">
              <span className="t">Play</span>
            </div>
            <div className="play">
              <div className="lead">
                <h2>
                  Sit down and play.
                  <br />
                  &spades; <span className="red">&hearts;</span>{" "}
                  <span className="red">&diams;</span> &clubs;
                </h2>
                <p>
                  Casual or rated, IMPs or matchpoints, humans or bots. Take any
                  open seat and deal.
                </p>
                <div className="cta">
                  <Link href="/table" className="btn pri">
                    Create Table
                  </Link>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
