/*
Copyright © 2023 Mattis Kristensen <mattismoel@gmail.com>
*/
package cmd

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/mattismoel/albumcut/types"
	"github.com/spf13/cobra"
)

var (
	youtubeLink  string
	albumTitle   string
	artist       string
	year         int
	inputPath    string
	coverArtPath string
	format       string
)

var rootCmd = &cobra.Command{
	Use:   "albumcut",
	Short: "Cut audio files given an input CSV file.",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if youtubeLink != "" {
			_, err := exec.LookPath("yt-dlp")
			if err != nil {
				log.Fatalf("Could not get YouTube link: %v\n", err)
			}
			command := exec.Command("yt-dlp", "-f", "140", "-o", "output.m4a", youtubeLink)

			err = command.Run()
			if err != nil {
				log.Fatalf("Could not run command: %v\n", err)
			}

			fmt.Println("Successfully downloaded YouTube video.")
			fmt.Println("Reading CSV file...")

			songs := []*types.Song{}

			f, err := os.Open(inputPath)
			if err != nil {
				log.Fatalf("Could not read CSV input file %s: %v\n", inputPath, err)
			}
			defer f.Close()

			csvReader := csv.NewReader(f)
			records, err := csvReader.ReadAll()
			if err != nil {
				log.Fatalf("Unable to parse CSV file %s: %v\n", inputPath, err)
			}

			var songTitle string
			var from int
			var to int
			var duration time.Duration

			for line, record := range records {
				song := &types.Song{}
				songTitle = record[0]
				if songTitle == "" {
					fmt.Printf("Could not parse the song at line %v", line)
					continue
				}

				from, err = timeToSeconds(record[1])
				if err != nil {
					log.Fatalf("Could not parse the timestamp at line %v: %v\n", line, err)
				}

				if record[2] == "" {

				}
				to, err = timeToSeconds(record[2])
				if err != nil {
					if to == -1 {

					}
					log.Fatalf("Could not parse the end time stamp at line %v: %v\n", line, err)
				}

				duration = time.Duration(to)*time.Second - time.Duration(from)*time.Second

				song.Title = songTitle
				song.From = from
				song.Duration = duration

				songs = append(songs, song)
			}

			for track_number, song := range songs {
				fileName := fmt.Sprintf("%s.%s", song.Title, format)
				fromString := fmt.Sprintf("%d", song.From)
				durationString := fmt.Sprintf("%d", int(song.Duration.Seconds()))

				// commandString := fmt.Sprintf("-ss %d -i %s -t %d %s", song.From, inputPath)
				fmt.Println(fileName, fromString, durationString)
				cutCmd := exec.Command("ffmpeg", "-ss", fromString, "-i", "output.m4a", "-t", durationString, fileName)
				err = cutCmd.Run()
				if err != nil {
					log.Fatalf("Could not cut the audio file: %v\n", err)
				}

				// ffmpeg
				// -i "Song 1.mp3"
				// -c copy
				// -metadata author="Test Artist"
				// -metadata album="Test Album"
				// -metadata album_artist="Test Album Artist"
				// -metadata year=2020
				// -metadata title="Title"
				// "Song 1 cp.mp3"

				// ffmpeg
				// -i "out.mp3"
				// -i "images.jpg"
				// -map 0:0
				// -map 1:0
				// -c copy
				// -metadata author="Test Artist"
				// -metadata lbum="Test Album"
				// -metadata album_artist="Test Album Artist"
				// -metadata year=2020
				// -metadata title="Title"
				// -metadata track=1
				// -metadata:s:v
				// title="Album cover"
				// -metadata:s:v comment="Cover (front)" outcover.mp3
				fmt.Printf("Attempting to add metadata to %s...", fileName)
				var out bytes.Buffer
				var stderr bytes.Buffer
				metadataCmd := exec.Command(
					"ffmpeg",
					"-i", fileName,
					"-i", coverArtPath,
					"-map", "0:0",
					"-map", "1:0",
					"-c", "copy",
					// "id3v2_version 3",
					"-metadata", fmt.Sprintf("author=%s", artist),
					"-metadata", fmt.Sprintf("artist=%s", artist),
					"-metadata", fmt.Sprintf("composer=%s", artist),
					"-metadata", fmt.Sprintf("album=%s", albumTitle),
					"-metadata", fmt.Sprintf("album_artist=%s", artist),
					"-metadata", fmt.Sprintf("year=%d", year),
					"-metadata", fmt.Sprintf("title=%s", song.Title),
					"-metadata", fmt.Sprintf("track=%d", track_number+1),
					"-metadata", fmt.Sprintf("title=%s", "Album Cover"),
					"-metadata:s:v", fmt.Sprintf("comment=%s", "Cover (front)"),
					fmt.Sprintf("%d - %s.%s", track_number+1, song.Title, format),
				)

				metadataCmd.Stdout = &out
				metadataCmd.Stderr = &stderr

				err = metadataCmd.Run()
				if err != nil {
					fmt.Println(fmt.Sprint(err), ":", stderr.String())
				}
				fmt.Println("result: ", out.String())

				err := exec.Command("rm", fileName).Run()
				if err != nil {
					log.Fatalf("Could not remove temp files: %v\n", err)
				}

			}

		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&inputPath, "input", "i", "input.csv", "The CSV input file containing cut information.")
	rootCmd.Flags().StringVarP(&coverArtPath, "cover", "c", "cover.png", "The path to the cover art.")
	rootCmd.Flags().StringVarP(&albumTitle, "albumTitle", "t", "", "The title of the album that is to be cut.")
	rootCmd.Flags().IntVarP(&year, "year", "y", -1, "The release year of the album.")
	rootCmd.Flags().StringVarP(&artist, "artist", "a", "", "The artist of the album.")
	rootCmd.Flags().StringVarP(&youtubeLink, "youtubeLink", "l", "", "Link to a YouTube video.")
	rootCmd.Flags().StringVarP(&format, "format", "f", "mp3", "The desired output format of tracks.")

	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("albumTitle")
	rootCmd.MarkFlagRequired("artist")
	rootCmd.MarkFlagRequired("year")
}

func timeToSeconds(timestamp string) (int, error) {
	parts := strings.Split(timestamp, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid timestamp format: %s", timestamp)
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	second, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, err
	}

	totalSeconds := hour*3600 + minute*60 + second
	return totalSeconds, nil
}
