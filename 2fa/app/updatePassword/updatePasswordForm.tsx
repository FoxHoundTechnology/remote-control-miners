// import Link from "next/link";
// import { headers } from "next/headers";
// import { createClient } from "@/utils/supabase/server";
// import { redirect } from "next/navigation";

// export default function Login({
//   searchParams,
// }: {
//   searchParams: { message: string };
// }) {
//   const signIn = async (formData: FormData) => {
//     "use server";

//     const email = formData.get("email") as string;
//     const password = formData.get("password") as string;
//     const supabase = createClient();

//     const { error } = await supabase.auth.signInWithPassword({
//       email,
//       password,
//     });

//     if (error) {
//       return redirect("/login?message=Could not authenticate user");
//     }

//     return redirect("/protected");
//   };

//   const signUp = async (formData: FormData) => {
//     "use server";

//     const origin = headers().get("origin");
//     const email = formData.get("email") as string;
//     const password = formData.get("password") as string;
//     const supabase = createClient();

//     const { error } = await supabase.auth.signUp({
//       email,
//       password,
//       options: {
//         emailRedirectTo: `${origin}/auth/callback`,
//       },
//     });

//     if (error) {
//       return redirect("/login?message=Could not authenticate user");
//     }

//     return redirect("/login?message=Check email to continue sign in process");
//   };
"use client";

import { useState } from "react";
import { SubmitButton } from "./submit-button";

export default function UpdatePasswordForm() {
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);

    const password = formData.get("password") as string;
    const confirmPassword = formData.get("confirm_password") as string;

    if (password !== confirmPassword) {
      setError("Passwords do not match");
      return;
    }

    if (password.length < 12) {
      setError("Password length must be at least 12 characters");
      return;
    }
  
    if (password.length > 64) {
      setError("Password length must be under 64 characters");
      return;
    }

    const response = await fetch("/api/updatePassword", {
      method: "POST",
      body: JSON.stringify({
        password,
        confirm_password: confirmPassword,
      }),
      headers: {
        "Content-Type": "application/json",
      },
    });

    const result = await response.json();

    if (result.error) {
      setError(result.error);
    } else {
      window.location.href = "/login";
    }
  };

  return (
    <div className="flex-1 flex flex-col w-full px-8 sm:max-w-md justify-center gap-2">
      <form onSubmit={handleSubmit} className="animate-in flex-1 flex flex-col w-full justify-center gap-2 text-foreground">
        {error && <p className="text-red-500">{error}</p>}
        <label className="text-md" htmlFor="password">
          Password
        </label>
        <input
          className="rounded-md px-4 py-2 bg-inherit border mb-1"
          type="password"
          name="password"
          placeholder="••••••••"
          required
        />
        <label className="text-md" htmlFor="confirm_password">
          Confirm Password
        </label>
        <input
          className="rounded-md px-4 py-2 bg-inherit border mb-1"
          type="password"
          name="confirm_password"
          placeholder="••••••••"
          required
        />
        <SubmitButton
          className="bg-green-700 border border-foreground/20 rounded-md px-4 py-2 text-foreground mb-2"
          pendingText="Updating..."
        >
          Update Password
        </SubmitButton>
      </form>
    </div>
  );
}
