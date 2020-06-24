
const ytRegexp = /(?:youtube\.com\/(?:[^\/]+\/.+\/|(?:v|e(?:mbed)?)\/|.*[?&]v=)|youtu\.be\/)([^"&?\/\s]{11})/

export const extractYoutubeID = (url: string): string | null => url.match(ytRegexp) ? (url.match(ytRegexp) as string[])[1] : null