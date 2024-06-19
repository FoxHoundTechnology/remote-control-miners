// "use client";
// import { useState } from "react";
// import { SubmitButton } from "./submit-button";

// export default function FA-Form({ searchParams }: { searchParams: { message: string } }) {
//   const [email, setEmail] = useState<string>('');
//   const [message, setMessage] = useState<string | null>(null);
//   const [error, setError] = useState<string | null>(null);

//   const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
//     e.preventDefault();

//     const response = await fetch('/api/auth/reset-password', {
//       method: 'POST',
//       headers: {
//         'Content-Type': 'application/json',
//       },
//       body: JSON.stringify({ email }),
//     });

//     if (!response.ok) {
//       const data = await response.json();
//       setError(data.error || 'Could not send reset email. Please try again.');
//     } else {
//       setMessage('Password reset email sent. Please check your inbox.');
//     }
//   };

//   return (
//     <div className="flex-1 flex flex-col w-full px-8 sm:max-w-md justify-center gap-2">
//       <h1 className="text-center text-3xl font-semibold mb-4 mt-20">Reset Password</h1>
//       <form onSubmit={handleSubmit} className="animate-in flex-1 flex flex-col w-full justify-center gap-2 text-foreground">
//         <label>
//           Enter registered email address
//         </label>
//         <input
//           className="rounded-md px-4 py-2 bg-inherit border mb-6"
//           type="email"
//           name="email"
//           placeholder="you@example.com"
//           value={email}
//           onChange={(e) => setEmail(e.target.value)}
//           required
//         />
//         <SubmitButton
//           className="bg-green-700 rounded-md px-4 py-2 text-foreground mb-2"
//           pendingText="Sending..."
//         >
//           Send Reset Email
//         </SubmitButton>
//       </form>
//       {message && <p className="text-green-500">{message}</p>}
//       {error && <p className="text-red-500">{error}</p>}
//     </div>
//   );
// }