'use server'

import { revalidatePath } from 'next/cache'
import { redirect } from 'next/navigation'

import { createClient } from '@/utils/supabase/server'

const handleSignUp = async (formData: FormData) => {

    const supabase = createClient()
    // PR: should be using the upabase.auth.signUp
    const { data, error } = await supabase.auth.signUp({
      email: formData.get('email') as string,
      password: formData.get('password') as string,
      phone: formData.get('number') as string
    })

    if (error) {
    
    //  setError(error.message)
    console.log('error', error)
    } else {
      console.log('data for sign up', data)
      redirect('/')
    }
  }

  export default handleSignUp