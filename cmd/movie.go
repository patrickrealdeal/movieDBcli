/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// movieCmd represents the movie command
var movieCmd = &cobra.Command{
	Use:     "movie",
	Short:   "Search for movie <movie>",
	Aliases: []string{"m"},
	RunE: func(cmd *cobra.Command, args []string) error {
		apiRoot := viper.GetString("api-root")

		return movieAction(os.Stdout, apiRoot, args)
	},
}

func movieAction(out io.Writer, apiRoot string, args []string) error {
	req := strings.Join(args, "+")
	if req == "" {
		return fmt.Errorf("You must provide a movie request: movie <name>.")
	}

	movie, err := getMovie(apiRoot, req)
	if err != nil {
		return err
	}

	details, err := getDetails(movie.MovieID)
	if err != nil {
		return err
	}

	return printMovie(out, movie, details)
}

func printMovie(out io.Writer, movie movie, details credits) error {
	t := time.Now()

	dir := func() string {
		str := ""
		for _, v := range details.Crew {
			if v.Job == "Director" {
				str = str + v.Name + ", "
			}
		}
		str = strings.TrimSuffix(str, ", ")
		return str
	}()

	cine := func() string {
		str := ""
		for _, v := range details.Crew {
			if v.Job == "Director of Photography" {
				str = str + v.Name + ", "
			}
		}
		str = strings.TrimSuffix(str, ", ")
		return str
	}()

	sp := func() string {
		str := ""
		for _, v := range details.Crew {
			if v.Job == "Screenplay" {
				str = str + v.Name + ", "
			}
		}
		str = strings.TrimSuffix(str, ", ")
		return str
	}()

	if len(sp) == 0 {
		_, err := fmt.Fprintf(out, "\nTitle:  %s\nDirector: %s\nCinematography: %s\n\nOverview: %s\n\nRating: %v\n",
			movie.Title, dir, cine, movie.Overview, movie.VoteAverage)

		t1 := time.Since(t)
		fmt.Fprintf(out, "%f\n", t1.Seconds())
		return err
	}

	_, err := fmt.Fprintf(out, "\nTitle:  %s\nDirector: %s\nCinematography: %s\nScreenplay: %s\n\nOverview: %s\n\nRating: %v\n",
		movie.Title, dir, cine, sp, movie.Overview, movie.VoteAverage)

	t1 := time.Since(t)
	fmt.Fprintf(out, "%f\n", t1.Seconds())

	return err
}

func init() {
	searchCmd.AddCommand(movieCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// movieCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// movieCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
