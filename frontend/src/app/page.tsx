import type { Metadata } from "next";
import Rail from "@/components/Rail";
import FriendsPanel from "./FriendsPanel";
import PlayModule from "./PlayModule";
import styles from "./home.module.css";

export const metadata: Metadata = {
  title: "Bridge++ — Home",
};

export default function HomePage() {
  return (
    <div className="app">
      <Rail active="home" />

      <div className="scroll">
        <div className={styles.top}>
          <h1>Bridge++ / home</h1>
        </div>

        <div className={styles.page}>
          <PlayModule />
        </div>
      </div>

      <FriendsPanel />
    </div>
  );
}
