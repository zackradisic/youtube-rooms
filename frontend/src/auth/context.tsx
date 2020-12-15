import React, { createContext, useState, useContext, useEffect } from 'react'
import Cookies from 'js-cookie'

import api, { setCredentials } from '../services/api'
import { useHistory } from 'react-router-dom'
import { User } from '../api/youtube-rooms-API'

export type AuthToken = {
  accessToken: string,
  accessExpiresAt: number,
  refreshToken: string,
  refreshExpiresAt: number
}

type AuthContextValue = {
  isAuthenticated: boolean,
  loading: boolean,
  user: User | null,
  setUser?: React.Dispatch<React.SetStateAction<User | null>>
  login: (email: string, password: string) => Promise<void>,
  logout: () => void,
  getUserInfo: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue>({
  isAuthenticated: false,
  loading: false,
  user: null,
  setUser: undefined,
  login: async () => {},
  logout: () => {},
  getUserInfo: async () => {}
})

export const AuthProvider = ({ children }: { children: React.ReactChild[] | React.ReactChild }) => {
  const [user, setUser] = useState<(User | null)>(null)
  const [loading, setLoading] = useState(true)
  const history = useHistory()

  useEffect(() => {
    const loadUserFromCookies = async () => {
      const token = Cookies.get('accessToken')
      if (token) {
        console.log(token)
        api.defaults.headers.Authorization = 'Bearer ' + token
        try {
          const { data: user } = await api.get('/me')
          if (user) setUser(user)
        } catch (err) {
          if (err.response) {
            if (err.response.status === 401) { // Access and refresh token are invalid
              alert('Your session has expired, please log back in.')
              logout()
              history.push('/signup')
            }
          }
        }
        // Add some logic to handle invalid or expired token
      }
      setLoading(false)
    }

    loadUserFromCookies()
  }, [])

  const login = async (email: string, password: string) => {
    const { data: token }: { data: AuthToken} = await api.post('/login', { email, password })
    console.log(token)
    if (token) {
      setCredentials(token.accessToken, token.accessExpiresAt, token.accessToken, token.refreshExpiresAt)
      const { data: user } = await api.get('/me')
      setUser(user)
    }
  }

  const logout = () => {
    Cookies.remove('accessToken')
    Cookies.remove('refreshToken')
    setUser(null)
  }

  const getUserInfo = async () => {
    const { data: user } = await api.get('/me')
    if (user) setUser(user)
  }

  return (
    <AuthContext.Provider value={{ isAuthenticated: !!user, user, login, loading, logout, setUser, getUserInfo }}>
      {children}
    </AuthContext.Provider>
  )
}

const useAuth = () => {
  const context = useContext(AuthContext)

  return context
}

/*
export const RequiresAuthentication = (Component: () => JSX.Element) => () => {
  const { isAuthenticated, loading } = useAuth()

  useEffect(() => {
    if (!isAuthenticated && !loading) router.push('/signup')
  }, [loading, isAuthenticated])

  return (
    <Component />
  )
}
*/

export default useAuth
