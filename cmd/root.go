/*
Copyright © 2023 Mattis Kristensen <mattismoel@gmail.com>
*/
package cmd

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	// "sync"
	"time"

	"github.com/fatih/color"
	"github.com/mattismoel/albumcut/types"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	youtubeLink  string
	albumTitle   string
	artist       string
	year         int
	inputPath    string
	coverArtPath string
	format       string
	outputDir    string
	clean        bool
)

var rootCmd = &cobra.Command{
	Use:   "albumcut",
	Short: "A tool specifically designed for downloading live concerts or similar from YouTube and chopping them into an album format with individual tracks.",
	Run: func(cmd *cobra.Command, args []string) {
		outputDir = strings.TrimSuffix(outputDir, "/")

		// If user does not wish to export to current folder
		if outputDir != "./" {
			outputDir += fmt.Sprintf("/%s - %s (%d)", artist, albumTitle, year)
		}

		// If no YouTube link is provided
		if youtubeLink == "" {
			return
		}

		err := downloadYoutubeVideo(youtubeLink)
		if err != nil {
			log.Fatalf("Could not download YouTube video: %v", err)
		}

		tracks, err := getTracksFromCSV(inputPath)
		if err != nil {
			log.Fatalf("Could not parse CSV file: %v\n", err)
		}

		// If output directory does not exist, create the directory
		if _, err := os.Stat(outputDir); errors.Is(err, os.ErrNotExist) {
			err = os.Mkdir(outputDir, os.ModePerm)
			if err != nil {
				log.Fatalf("Could not create output directory at %s: %v\n", outputDir, err)
			}
		}

		err = exportTracks(tracks, outputDir)
		if err != nil {
			log.Fatalf("Could not create tracks at location %s: %v\n", outputDir, err)
		}

		if clean {
			err := cleanUp()
			if err != nil {
				log.Fatalf("Could not clean up: %v\n", err)
			}
		}
	},
}

func cleanUp() error {
	err := os.Remove("output.m4a")
	if err != nil {
		return err
	}

	err = os.Remove(inputPath)
	if err != nil {
		return err
	}

	err = os.Remove(coverArtPath)
	if err != nil {
		return err
	}

	return nil
}

func exportTrack(track *types.Track, outPath string) error {
	args := []string{"-i", "output.m4a", "-ss", strconv.Itoa(track.From)}

	duration := getTrackDuration(track)
	fileName := fmt.Sprintf("%s/%s.%s", outPath, track.Title, format)

	if track.To == -1 {
		args = append(args, fileName)
	} else {
		args = append(args, "-t", strconv.Itoa(int(duration.Seconds())), fileName)
	}

	command := exec.Command("ffmpeg", args...)

	var stderr bytes.Buffer
	command.Stderr = &stderr

	err := command.Run()
	if err != nil {
		return fmt.Errorf("Could not process audio for file %s: %v: %v\n", fileName, err, stderr.String())
	}

	err = addMetadata(track)
	if err != nil {
		return fmt.Errorf("Could not add metadata: %v\n", err)
	}

	err = os.Remove(fileName)
	if err != nil {
		return err
	}

	// err = os.Remove("output.m4a")
	// if err != nil {
	// 	return err
	// }

	return nil
}

func exportTracks(tracks []*types.Track, outPath string) error {
	for _, track := range tracks {
		err := exportTrack(track, outPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func addMetadata(track *types.Track) error {
	defer fmt.Printf("Added metadata to %q successfully\n", track.Title)

	filename := fmt.Sprintf("%d - %s.%s", track.TrackNumber, track.Title, format)
	inputPath := fmt.Sprintf("%s/%s.%s", outputDir, track.Title, format)
	outputPath := fmt.Sprintf("%s/%s", outputDir, filename)

	args := []string{
		"-i", inputPath,
		"-i", coverArtPath,
		"-map", "0:0",
		"-map", "1:0",
		"-c", "copy",
		"-metadata", fmt.Sprintf("author=%s", artist),
		"-metadata", fmt.Sprintf("artist=%s", artist),
		"-metadata", fmt.Sprintf("composer=%s", artist),
		"-metadata", fmt.Sprintf("album=%s", albumTitle),
		"-metadata", fmt.Sprintf("album_artist=%s", artist),
		"-metadata", fmt.Sprintf("year=%d", year),
		"-metadata", fmt.Sprintf("title=%s", track.Title),
		"-metadata", fmt.Sprintf("track=%d", track.TrackNumber),
		"-metadata:s:v", fmt.Sprintf("comment=%s", "Cover (front)"),
		outputPath,
	}

	command := exec.Command("ffmpeg", args...)

	var stderr bytes.Buffer
	var out bytes.Buffer

	command.Stdout = &out
	command.Stderr = &stderr

	err := command.Run()
	if err != nil {
		return fmt.Errorf("%v: %v\n", err, stderr.String())
	}

	return nil
}

func downloadYoutubeVideo(link string) error {
	_, err := exec.LookPath("yt-dlp")
	if err != nil {
		return fmt.Errorf("no such command 'yt-dlp': %v\n", err)
	}

	defer log.Printf("Downloaded audio from YouTube video %q successfully\n", link)
	cmdArguments := []string{"-f", "140", "-o", "output.m4a", youtubeLink}

	command := exec.Command("yt-dlp", cmdArguments...)
	var stderr bytes.Buffer
	command.Stderr = &stderr
	err = command.Run()
	if err != nil {
		return fmt.Errorf("%v:%v", err, stderr.String())
	}

	return nil
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
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", "./", "The directory to where the tracks are to be exported. AlbumCut will create the album directory itself. Defaults to current directory.")
	rootCmd.Flags().BoolVar(&clean, "clean", true, "Clean up files after export (CSV file, cover art)")

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

func getTracksFromCSV(csvPath string) ([]*types.Track, error) {
	defer fmt.Printf("Successfully parsed CSV file %s\n", csvPath)
	tracks := []*types.Track{}

	// Attempt to open file and parse CSV contents
	f, err := os.Open(csvPath)
	if err != nil {
		log.Fatalf("Could not read CSV input file %s: %v\n", inputPath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var title string
	var from int
	var to int

	fmt.Printf("Found %d tracks:\n", len(records))
	for line, record := range records {
		track := &types.Track{}
		title = record[0]
		from, err = timeToSeconds(record[1])

		track.TrackNumber = line + 1
		track.From = from
		track.Title = title

		fmt.Printf("%d - %q\n", track.TrackNumber, track.Title)

		// If last track and end is not specified
		if record[2] == "" {
			track.To = -1
			tracks = append(tracks, track)
			break
		}

		from, err = timeToSeconds(record[1])
		if err != nil {
			log.Fatalf("Could not parse the timestamp at line %v: %v\n", line, err)
		}

		to, err = timeToSeconds(record[2])

		if err != nil {
			log.Fatalf("Could not parse date: %v\n", err)
		}

		track.Title = title
		track.TrackNumber = line + 1
		track.From = from
		track.To = to

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func getTrackDuration(track *types.Track) time.Duration {
	duration := time.Duration(track.To)*time.Second - time.Duration(track.From)*time.Second
	return duration
}
