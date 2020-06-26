import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { Video } from '../../api/youtube-rooms-API'

export interface PlayerDetails {
    isPlaying: boolean,
    current: Video,
    seekTo: number
}

const initialState: PlayerDetails = {
  isPlaying: false,
  current: {
    url: '',
    title: '',
    requester: ''
  },
  seekTo: -1
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
    },
    seekTo (state, action: PayloadAction<number>) {
      state.seekTo = action.payload
    }
  }
})

export const {
  setPlaying,
  setCurrent,
  seekTo
} = playerSlice.actions

export default playerSlice.reducer
