'use client'

import { createClient } from '@/utils/supabase/server'
import { useState } from 'react'
import Link from 'next/link'

import handleSignUp from './action'
import { SubmitButton } from './submit-button'

export default function Register({ searchParams }: { searchParams: { message: string } }) {
  const [error, setError] = useState<string | null>(null)

  return (
    <div className='flex-1 flex flex-col w-full px-8 sm:max-w-md justify-center gap-2'>
      <Link
        href='/'
        className='absolute left-8 top-8 py-2 px-4 rounded-md no-underline text-foreground bg-btn-background hover:bg-btn-background-hover flex items-center group text-sm'
      >
        <svg
          xmlns='http://www.w3.org/2000/svg'
          width='24'
          height='24'
          viewBox='0 0 24 24'
          fill='none'
          stroke='currentColor'
          strokeWidth='2'
          strokeLinecap='round'
          strokeLinejoin='round'
          className='mr-2 h-4 w-4 transition-transform group-hover:-translate-x-1'
        >
          <polyline points='15 18 9 12 15 6' />
        </svg>{' '}
        Back
      </Link>
      <form className='animate-in flex-1 flex flex-col w-full justify-center gap-2 text-foreground'>
        <label className='text-md' htmlFor='email'>
          Email
        </label>
        <input
          className='rounded-md px-4 py-2 bg-inherit border mb-6'
          name='email'
          placeholder='you@example.com'
          required
        />
        <label className='text-md'>Phone Number</label>
        <input
          className='rounded-md px-4 py-2 bg-inherit border mb-6'
          name='number'
          placeholder='+1 (XXX) XXX-XXXX'
          required
        />
        <label className='text-md' htmlFor='password'>
          Password
        </label>
        <input
          className='rounded-md px-4 py-2 bg-inherit border mb-6'
          type='password'
          name='password'
          placeholder='••••••••'
          required
        />
        <label className='text-md' htmlFor='confirm_password'>
          Confirm Password
        </label>
        <input
          className='rounded-md px-4 py-2 bg-inherit border mb-6'
          type='password'
          name='confirm_password'
          placeholder='••••••••'
          required
        />
        {/* <button
        type="submit"
        className="bg-green-700 border border-foreground/20 rounded-md px-4 py-2 text-foreground mb-2"
      >
        Sign Up
      </button> */}
        <SubmitButton
          formAction={handleSignUp}
          className='bg-green-700 border border-foreground/20 rounded-md px-4 py-2 text-foreground mb-2'
          pendingText='Signing Up...'
        >
          Sign Up
        </SubmitButton>
        <div className='mt-4 text-center'>
          <span>Already have an account? </span>
          <a href='/login' className='text-blue-500 hover:underline'>
            Sign in.
          </a>
        </div>
        {error && <p className='mt-4 p-4 bg-red-100 text-red-700 text-center'>{error}</p>}
      </form>{' '}
      {searchParams?.message && (
        <p className='mt-4 p-4 bg-foreground/10 text-foreground text-center'>{searchParams.message}</p>
      )}
    </div>
  )
}
