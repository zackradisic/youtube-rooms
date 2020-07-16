import axios from 'axios'

const ytRegexp = /(?:youtube\.com\/(?:[^\/]+\/.+\/|(?:v|e(?:mbed)?)\/|.*[?&]v=)|youtu\.be\/)([^"&?\/\s]{11})/

export const extractYoutubeID = (url: string): string | null => url.match(ytRegexp) ? (url.match(ytRegexp) as string[])[1] : null

export const verifyRoomPassword = async (password: string, roomName: string): Promise<boolean> => {
  var axios = require('axios')
  var data = JSON.stringify({ roomName: "zack's room", password: 'sfsdf' })

  var config = {
    method: 'post',
    url: 'http://localhost/api/rooms/verify/',
    headers: {
      'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36',
      'Content-Type': 'application/json'
    },
    data: data
  }

  try {
    const res = await axios(config)
    return !!res.data.success
  } catch (err) {
    console.log(err)
    return false
  }
}
