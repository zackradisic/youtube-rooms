import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { User } from '../../api/youtube-rooms-API'

type UsersState = Array<User>

const initialState: UsersState = []

interface RemoveUserPayload {
  keyName: keyof User,
  val: User[keyof User]
}

const usersDisplaySlice = createSlice({
  name: 'usersDisplay',
  initialState,
  reducers: {
    addUser (state, action: PayloadAction<User>) {
      state.push(action.payload)
    },
    removeUser (state, action: PayloadAction<RemoveUserPayload>) {
      const { keyName, val } = action.payload
      state = state.splice(state.findIndex(user => user[keyName] === val), 1)
    },
    setUsers (state, action: PayloadAction<User[]>) {
      state = action.payload
    }
  }
})

export const {
  addUser,
  removeUser,
  setUsers
} = usersDisplaySlice.actions

export default usersDisplaySlice.reducer
