import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface RoomCredentials {
    name: string,
    password?: string
}

const initialState: RoomCredentials = {
  name: '',
  password: ''
}

const roomCredentialsSlice = createSlice({
  name: 'roomCredentials',
  initialState,
  reducers: {
    setCredentials (state, action: PayloadAction<RoomCredentials>) {
      state.name = action.payload.name
      state.password = action.payload.password
    }
  }
})

export const {
  setCredentials
} = roomCredentialsSlice.actions

export default roomCredentialsSlice.reducer
