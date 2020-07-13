import { combineReducers } from '@reduxjs/toolkit'

import videosDisplayReducer from '../features/videosDisplay'
import usersDisplayReducer from '../features/usersDisplay'
import playerReducer from '../features/player/playerSlice'
import roomCredentialsReducer from '../features/roomCredentials'
import roomPreviewsReducer from '../features/roomPreviewsDisplay'

const rootReducer = combineReducers({
  videosDisplay: videosDisplayReducer,
  usersDisplay: usersDisplayReducer,
  player: playerReducer,
  roomCredentials: roomCredentialsReducer,
  roomPreviews: roomPreviewsReducer
})

export type RootState = ReturnType<typeof rootReducer>

export default rootReducer
