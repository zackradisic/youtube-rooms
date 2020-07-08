import { combineReducers } from '@reduxjs/toolkit'

import videosDisplayReducer from '../features/videosDisplay'
import usersDisplayReducer from '../features/usersDisplay'
import playerReducer from '../features/player/playerSlice'
import roomCredentialsReducer from '../features/roomCredentials'

const rootReducer = combineReducers({
  videosDisplay: videosDisplayReducer,
  usersDisplay: usersDisplayReducer,
  player: playerReducer,
  roomCredentials: roomCredentialsReducer
})

export type RootState = ReturnType<typeof rootReducer>

export default rootReducer
