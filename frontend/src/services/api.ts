import axios, { AxiosError } from 'axios'
import Cookies from 'js-cookie'

const urls = {
  development: 'http://localhost:80/'
}

const api = axios.create({
  baseURL: `${process.env.REACT_APP_API_URL}/api/`,
  headers: {
    Accept: 'application/json',
    'Content-Type': 'application/json'
  }
})

export const setCredentials = (accessToken: string, accessExpiresAt: number, refreshToken: string, refreshExpiresAt: number) => {
  Cookies.set('accessToken', accessToken, { expires: new Date(accessExpiresAt * 1000) })
  api.defaults.headers.Authorization = 'Bearer ' + accessToken

  if (refreshToken && refreshExpiresAt) Cookies.set('refreshToken', refreshToken, { expires: new Date(refreshExpiresAt * 1000) })
}

export const handleErr = (response: any, logout: () => void, history: any) => {
  if (response) {
    if (response.status === 401) {
      logout()
      alert('Your session has expired. Please log back in.')
      history.push('/login')
      return true
    }
  }
}

export default api
