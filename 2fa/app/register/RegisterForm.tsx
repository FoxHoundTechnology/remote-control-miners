"use client"; // will run in the browser

import { useState } from "react";
import { SubmitButton } from "./submit-button";

export default function RegisterForm() {
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);

    const response = await fetch("/api/register", {
      method: "POST",
      body: JSON.stringify({
        email: formData.get("email"),
        number: formData.get("number"),
        password: formData.get("password"),
        confirm_password: formData.get("confirm_password"),
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
    <form
      className="animate-in flex-1 flex flex-col w-full justify-center gap-2 text-foreground"
      onSubmit={handleSubmit}
    >
      <label className="text-md" htmlFor="email">
        Email
      </label>
      <input
        className="rounded-md px-4 py-2 bg-inherit border mb-6"
        name="email"
        placeholder="you@example.com"
        required
      />
      <label className="text-md">Phone Number</label>
      <input
        className="rounded-md px-4 py-2 bg-inherit border mb-6"
        name="number"
        placeholder="+1 (XXX) XXX-XXXX"
        required
      />
      <label className="text-md" htmlFor="password">
        Password
      </label>
      <input
        className="rounded-md px-4 py-2 bg-inherit border mb-6"
        type="password"
        name="password"
        placeholder="••••••••"
        required
      />
      <label className="text-md" htmlFor="confirm_password">
        Confirm Password
      </label>
      <input
        className="rounded-md px-4 py-2 bg-inherit border mb-6"
        type="password"
        name="confirm_password"
        placeholder="••••••••"
        required
      />
      {/* <button
        type="submit"
        className="bg-green-700 border border-foreground/20 rounded-md px-4 py-2 text-foreground mb-2"
      >
        Sign Up
      </button> */}
         <SubmitButton
           formAction={RegisterForm}
           className="bg-green-700 border border-foreground/20 rounded-md px-4 py-2 text-foreground mb-2"
           pendingText="Signing Up..."
         >
           Sign Up
         </SubmitButton>
      <div className="mt-4 text-center">
        <span>Already have an account? </span>
        <a href="/login" className="text-blue-500 hover:underline">
          Sign in.
        </a>
      </div>
      {error && (
        <p className="mt-4 p-4 bg-red-100 text-red-700 text-center">
          {error}
        </p>
      )}
    </form>
  );
}
