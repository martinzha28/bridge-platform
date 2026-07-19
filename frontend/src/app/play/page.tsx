import { redirect } from "next/navigation";

// The lobby has no backend yet; "/play" just returns to home.
export default function PlayPage() {
  redirect("/");
}
