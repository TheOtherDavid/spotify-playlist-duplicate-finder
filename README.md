# spotify-playlist-duplicate-finder
For a given Spotify playlist, find any duplicates, based on the track name as a string, including remixes.

One of my friends has a great Spotify playlist, with like 1700 great songs. But the issue is, she's convinced there are duplicates in there, but it would be too hard to sort through all 1700 songs to find them. Even worse, sometimes a song appears on multiple different albums. Sometimes a band will release a song as a single, and then later release it again on their full album. Because these are officially different tracks, Spotify itself can't detect the duplicate based on track ID alone. By comparing songs based on their name as a string, and not based on the internal Spotify ID, these potential duplicates can be found and dealt with. Because we all know how programs can go out of control, I am not automatically deleting these duplicates, but simply writing them to a file. Deletion can be done manually.

This app uses a few environment variables. These are attributes of your Spotify API account, and should be self-explanatory.

SPOTIFY_ID

SPOTIFY_SECRET

REFRESH_TOKEN

USER_ID

The app takes in two arguments, a Spotify playlist ID and the desired filename. If the playlist ID is left unspecified, the program will generate an error. If the filename is left unspecified, a default of "duplicate-spotify-track" will be used.

To locate your Spotify playlist ID, go to the playlist in your Spotify app, More Options > Share > Copy Link to Playlist. This will copy something like this to your clipboard:

https://open.spotify.com/playlist/413xP6b7H2IBnWvf6BSRSY?si=242fdaecf6e568f1
(link deliberately changed to be broken)

The key piece of this is the playlistId, which in this example, is 413xP6b7H2IBnWvf6BSRSY. This is the proper argument to use for this app.
