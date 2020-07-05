import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { User } from '../../api/youtube-rooms-API'

const initialState: { users: User[] } = { users: [] }

interface RemoveUserPayload {
  keyName: keyof User,
  val: User[keyof User]
}

const usersDisplaySlice = createSlice({
  name: 'usersDisplay',
  initialState,
  reducers: {
    addUser (state, action: PayloadAction<User>) {
      state.users.push(action.payload)
    },
    removeUser (state, action: PayloadAction<RemoveUserPayload>) {
      const { keyName, val } = action.payload
      state.users = state.users.splice(state.users.findIndex(user => user[keyName] === val), 1)
    },
    setUsers (state, action: PayloadAction<User[]>) {
      state.users = action.payload
    }
  }
})

export const {
  addUser,
  removeUser,
  setUsers
} = usersDisplaySlice.actions

export default usersDisplaySlice.reducer
