import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { User } from '../../api/youtube-rooms-API'

type UsersState = Array<User>

const initialState: UsersState = []

interface RemoveUserPayload {
  keyName: keyof User,
  val: User[keyof User]
}

const usersDisplaySlice = createSlice({
  name: 'userssDisplay',
  initialState,
  reducers: {
    addUser (state, action: PayloadAction<User>) {
      state.push(action.payload)
    },
    removeUser (state, action: PayloadAction<RemoveUserPayload>) {
      const { keyName, val } = action.payload
      state = state.splice(state.findIndex(user => user[keyName] === val), 1)
    }
  }
})

export const {
  addUser,
  removeUser
} = usersDisplaySlice.actions

export default usersDisplaySlice.reducer
