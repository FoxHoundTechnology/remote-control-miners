import { createContext, useEffect, useState, ReactNode } from 'react'

import { useRouter } from 'next/router'

import axios from 'axios'

import authConfig from 'src/configs/auth'

import { AuthValuesType, LoginParams, ErrCallbackType, UserDataType } from './types'

const defaultProvider: AuthValuesType = {
  user: null,
  loading: true,
  setUser: () => null,
  setLoading: () => Boolean,
  login: () => Promise.resolve(),
  logout: () => Promise.resolve(),
  register: () => Promise.resolve()
}

const AuthContext = createContext(defaultProvider)

// base URLs for axios
const userAuthUrl = process.env.REMOTE_CONTROL_SERVER_URL
export const authAxios = axios.create({
  baseURL: `${userAuthUrl}/users`
})

type Props = {
  children: ReactNode
}

const AuthProvider = ({ children }: Props) => {
  const [user, setUser] = useState<UserDataType | null>(defaultProvider.user)
  const [loading, setLoading] = useState<boolean>(defaultProvider.loading)

  const router = useRouter()

  useEffect(() => {
    const initAuth = async (): Promise<void> => {
      const storedToken = window.localStorage.getItem('accessToken')!
      if (storedToken) {
        setLoading(true)
        await authAxios
          .post(
            '/jwt/header',
            {},
            {
              headers: {
                Authorization: storedToken
              }
            }
          )
          .then(async response => {
            // success case
            setLoading(false)
            // TODO: modify the backend
            const userData = response?.data?.data

            setUser({
              id: userData?.ID,
              username: userData?.user_name,
              email: userData?.email,
              role: userData?.role || 'technician'
            })
          })
          .catch(err => {
            setLoading(false)
            localStorage.removeItem('userData')
            localStorage.removeItem('refreshToken')
            localStorage.removeItem('accessToken')

            setUser(null)

            if (!router.pathname.includes('login')) {
              router.replace('/login')
            }
          })
      } else {
        setLoading(false)
      }
    }

    initAuth()
  }, [])

  // *WIP
  // NOTE: this should be done by only admins
  const handleRegisterUser = async (params: UserDataType, errorCallback?: ErrCallbackType) => {
    authAxios
      .post('', {
        user_name: params?.username,
        password: params?.password,
        email: params?.email
      })
      .then(async response => {
        const userData = response?.data?.data
        const returnUrl = '/login'

        setUser({
          id: userData?.ID,
          username: userData?.user_name,
          email: userData?.email,
          // TODO: need the logic to set up a default role (i.e. client)
          //       in order for clients to be able to
          //       put some access control constraints on their own hosting clients
          role: 'technician'
        })

        window.localStorage.setItem(
          'userData',
          JSON.stringify({
            id: userData?.ID,
            username: userData?.user_name,
            email: userData?.email,
            role: 'technician'
          })
        )

        // const redirectURL = returnUrl ? returnUrl : '/'
        router.replace('/')
      })
      .catch(err => {
        if (errorCallback) errorCallback(err)
      })
  }

  const handleLogin = (params: LoginParams, errorCallback?: ErrCallbackType) => {
    authAxios
      .post('/login', params)
      .then(async response => {
        const userData = response?.data?.data

        params.rememberMe ? window.localStorage.setItem('accessToken', response?.data?.accessToken) : null
        // setLoading(false)

        const returnUrl = router.query.returnUrl

        // NOTE: it has to follow the same structure as the user data
        // TODO: redo the create account UI once again
        setUser({
          id: userData?.ID,
          username: userData?.user_name,
          email: userData?.email,
          role: userData?.role || 'technician'
        })

        params.rememberMe
          ? window.localStorage.setItem(
              'userData',
              JSON.stringify({
                id: userData?.ID,
                username: userData?.user_name,
                email: userData?.email,
                role: userData?.role || 'technician'
              })
            )
          : null

        const redirectURL = returnUrl && returnUrl !== '/' ? returnUrl : '/'
        router.replace(redirectURL as string)
      })

      .catch(err => {
        if (errorCallback) errorCallback(err)
      })
  }

  const handleLogout = () => {
    // TODO: add the logout endpoint

    setUser(null)
    window.localStorage.removeItem('userData')
    window.localStorage.removeItem(authConfig.storageTokenKeyName)
    router.push('/login')
  }

  const values = {
    user,
    loading,
    setUser,
    setLoading,
    register: handleRegisterUser,
    login: handleLogin,
    logout: handleLogout
  }

  return <AuthContext.Provider value={values}>{children}</AuthContext.Provider>
}

export { AuthContext, AuthProvider }
