import { redirect } from "next/navigation";

// A table always has an id now; go make one.
export default function TableIndexPage() {
  redirect("/create");
}
