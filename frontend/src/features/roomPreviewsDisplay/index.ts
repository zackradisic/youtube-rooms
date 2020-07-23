import axios from 'axios'
import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { RoomPreview } from '../../api/youtube-rooms-API'
import { AppThunk } from '../../app/store'

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

export const fetchRoomPreview = (name?: string): AppThunk => async dispatch => {
  const url = `https://${process.env.REACT_APP_API_URL}/api/rooms${name ? '?name=' + encodeURI(name as string) : ''}`
  try {
    const res = await axios.get(url)
    if (!res.data.rooms) return
    if (res.data.rooms.length === 0) return

    if ((res.data.rooms as RoomPreview[])[0].id) {
      dispatch(setRoomPreviews(res.data.rooms as RoomPreview[]))
    }
  } catch (err) {
    console.log(err)
  }
}

export const {
  setRoomPreviews,
  addRoomPreview,
  removeRoomPreview
} = roomPreviewsSlice.actions

export default roomPreviewsSlice.reducer
