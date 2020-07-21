import axios from 'axios'

const ytRegexp = /(?:youtube\.com\/(?:[^\/]+\/.+\/|(?:v|e(?:mbed)?)\/|.*[?&]v=)|youtu\.be\/)([^"&?\/\s]{11})/

export const extractYoutubeID = (url: string): string | null => url.match(ytRegexp) ? (url.match(ytRegexp) as string[])[1] : null

export const verifyRoomPassword = async (password: string, roomName: string): Promise<boolean> => {
  var axios = require('axios')
  var data = JSON.stringify({ roomName: roomName, password: password })

  var config = {
    method: 'post',
    url: 'https://api.theatre.theradisic.com/api/rooms/verify/',
    data: data,
    validateStatus: () => true
  }

  const res = await axios(config)
  if (res.status !== 200) {
    throw new Error(res.data.error)
  }
  return !!res.data.success
}
