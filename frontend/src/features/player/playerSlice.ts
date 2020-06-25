import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { Video } from '../../api/youtube-rooms-API'

export interface PlayerDetails {
    isPlaying: boolean,
    current: Video
}

const initialState: PlayerDetails = {
  isPlaying: false,
  current: {
    url: '',
    title: '',
    requester: ''
  }
}

const playerSlice = createSlice({
  name: 'player',
  initialState,
  reducers: {
    setPlaying (state, action: PayloadAction<boolean>) {
      state.isPlaying = action.payload
    },
    setCurrent (state, action: PayloadAction<Video>) {
      state.current = action.payload
    }
  }
})

export const {
  setPlaying,
  setCurrent
} = playerSlice.actions

export default playerSlice.reducer
