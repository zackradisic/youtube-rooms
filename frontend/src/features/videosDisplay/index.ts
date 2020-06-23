import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { Video } from '../../api/youtube-rooms-API'

type VideoState = Array<Video>

const initialState: VideoState = []

const videosDisplaySlice = createSlice({
  name: 'videosDisplay',
  initialState,
  reducers: {
    addVideo (state, action: PayloadAction<Video>) {
      state.push(action.payload)
    }
  }
})

export const {
  addVideo
} = videosDisplaySlice.actions

export default videosDisplaySlice.reducer
