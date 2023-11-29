# AlbumCut

A tool specifically designed for downloading live concerts or similar from YouTube and chopping them into an album format with individual tracks. 

# Dependencies 

AlbumCut requires you to have the following packages installed on your system:

- `ffmpeg`: [Download and install here](https://ffmpeg.org/download.html)
- `yt-dlp`: [GitHub Page](https://github.com/yt-dlp/yt-dlp)


# Usage

Clone this repositiory to a directory of choice. The the executables are located in the `bin` directory. Here are versions for the following operating systems:

| OS      | Binary             |
|---------|--------------------|
| Linux   | `albumcut-linux`   |
| MacOS   | `albumcut-macos`   |
| Windows | `albumcut-windows` |

The program requires a `.csv` file containing the tracklist. This `.csv` file should be of the following format:

```csv
track_title, from, to,
...
...
```

An example file could look like this:

```csv
Maria También,00:00:01,00:03:16
August 10,00:03:17,00:07:46
White Gloves,00:07:47,
```
> 🖐️ Notice that the last track does not have a "to" timestamp. This ensures that the last track will go on till the end of the video. Optionally a end time stamp can be provided.

The command can then be run in your terminal window with the following command:

```bash
$ ./bin/albumcut-<os> --youtubeLink="https://youtu.be/<video_id>" --artist="Artist Name" --albumTitle="Album Title" --year 2023 --cover="cover.png"
```

> 🖐️ It is important that the YouTube link provided is the one obtained via the "Share" button!

This creates a folder at the current directory with the exported chopped youtube video.

> More help can be found using the `albumcut --help` command.

## Recommended workflow

- Rename the desired version of `albumcut-<os_version>` to `albumcut` 
- Make an alias for the new version, giving access to the command from anywhere
- Open a terminal window in the location you desire to download the album (for example `~/Music`)
- Download cover art and create the `input.csv` file at that location (example `~/Music/input.csv` and `~/Music/cover.png`)
- Run `albumcut` with fitting arguments. The command will clean up any files used (cover art, CSV file etc.)
