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

```
track_title, from, to,
...
```

An example file could look like this:

```
Maria También,00:00:01,00:03:16
August 10,00:03:17,00:07:46
White Gloves,00:07:47,

```
> Notice that the last track does not have a "to" timestamp. This ensures that the last track will go on till the end of the video. Optionally a end time stamp can be provided.

The command can then be run in your terminal window with the following command:

```bash
$ ./bin/albumcut-<os> --youtubeLink="https://youtu.be/<video_id>" --artist="Artist Name" --albumTitle="Album Title" --year 2023 --cover="cover.png"
```

This creates a folder at the current directory with the exported chopped youtube video.

> More help can be found using the `albumcut --help` command.
