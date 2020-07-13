import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { RoomPreview } from '../../api/youtube-rooms-API'

export interface RoomPreviews {
    rooms: RoomPreview[]
}

const initialState: RoomPreviews = {
  rooms: []
}

const roomPreviewsSlice = createSlice({
  name: 'roomPreviews',
  initialState,
  reducers: {
    setRoomPreviews (state, action: PayloadAction<RoomPreview[]>) {
      state.rooms = action.payload
    },
    addRoomPreview (state, action: PayloadAction<RoomPreview>) {
      state.rooms.push(action.payload)
    },
    removeRoomPreview (state, action: PayloadAction<RoomPreview>) {
      const { id } = action.payload
      const i = state.rooms.findIndex(e => e.id === id)
      if (i === -1) return
      state.rooms.splice(i, 1)
    }
  }
})

export const {
  setRoomPreviews,
  addRoomPreview,
  removeRoomPreview
} = roomPreviewsSlice.actions

export default roomPreviewsSlice.reducer
