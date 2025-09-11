import { playlist } from '@/services/audio/playlist.js';

const totalPlaylistDuration = playlist.reduce(
    (sum, song) => sum + song.duration,
    0
);

export function getCurrentTrack() {
    const now = new Date();

    const secondsSinceUTCMidnight =
        now.getUTCHours() * 3600 +
        now.getUTCMinutes() * 60 +
        now.getUTCSeconds() +
        now.getUTCMilliseconds() / 1000;

    const currentPlaylistPosition =
        secondsSinceUTCMidnight % totalPlaylistDuration;

    let accumulatedDuration = 0;
    for (const song of playlist) {
        const songEndTime = accumulatedDuration + song.duration;
        if (currentPlaylistPosition < songEndTime) {
            const currentTimeInSong =
                currentPlaylistPosition - accumulatedDuration;
            return {
                url: song.url,
                currentTime: currentTimeInSong,
            };
        }
        accumulatedDuration = songEndTime;
    }

    return { url: playlist[0].url, currentTime: 0 };
}
