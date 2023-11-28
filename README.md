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

The command can then be run in your terminal window with the following command:

```bash
$ ./bin/albumcut-<os> --youtubeLink="https://youtu.be/<video_id>" --artist="Artist Name" --albumTitle="Album Title" --year 2023 --cover="cover.png"
```
