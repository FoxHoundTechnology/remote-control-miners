import { NextApiRequest, NextApiResponse } from "next";
import { createClient } from "@/utils/supabase/server";
import { NextResponse } from 'next/server';

export async function POST(req: Request) {
  const { email, number, password, confirm_password } = await req.json();

  if (number.length < 8) {
    return NextResponse.json({ error: "Phone number must be at least 8 digits" }, { status: 400 });
  }

  if (password.length < 12) {
    return NextResponse.json({ error: "Password must be at least 12 characters long" }, { status: 400 });
  }

  if (password.length > 64) {
    return NextResponse.json({ error: "Password must be less than 64 characters" }, { status: 400 });
  }

  if (password !== confirm_password) {
    return NextResponse.json({ error: "Passwords do not match" }, { status: 400 });
  }


  const supabase = createClient();

  const { data, error: fetchError } = await supabase
    .from('users')
    .select('id')
    .eq('email',email)
    .single();

  if (data) {
    return NextResponse.json({ error: 'Email already registered' }, { status: 400 });
  }

  const { error } = await supabase.auth.signUp({
    email,
    password,
    options: {
      emailRedirectTo: `${req.headers.get('origin')}/auth/callback`,
    },
  });

  if (error) {
    return NextResponse.json({ error: "Could not authenticate user" }, { status: 400 });
  }

  return NextResponse.json({ success: true });
}
