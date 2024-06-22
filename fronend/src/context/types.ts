export type ErrCallbackType = (err: { [key: string]: string }) => void

export type LoginParams = {
  email: string
  password: string
  rememberMe?: boolean
}

export type UserDataType = {
  id?: number
  role?: string
  email: string

  fullName?: string
  username?: string
  password?: string
  lastLogin?: string

  // TODO: undo here
  //phoneNumber?: string
  //avatar?: string | null
}

export type AuthValuesType = {
  loading: boolean
  logout: () => void
  user: UserDataType | null
  setLoading: (value: boolean) => void
  setUser: (value: UserDataType | null) => void
  login: (params: LoginParams, errorCallback?: ErrCallbackType) => void
  register: (params: UserDataType, errorCallback?: ErrCallbackType) => void
}
