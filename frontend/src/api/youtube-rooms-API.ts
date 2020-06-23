
export interface User {
    discordID: string,
    discordUsername: string,
    discordDiscriminator: string,
    isOwner: boolean
}

export interface Video {
    url: string,
    title: string,
    requester: User['discordID']
}
