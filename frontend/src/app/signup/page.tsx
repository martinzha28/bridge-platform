import type { Metadata } from "next";
import AuthForm from "@/components/AuthForm";

export const metadata: Metadata = { title: "Bridge++ — Sign up" };

export default function SignupPage() {
  return <AuthForm mode="signup" />;
}
