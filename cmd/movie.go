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

	jobs := []crew{}
	for _, c := range details.Crew {
		if c.Job == "Director" || c.Job == "Director of Photography" || c.Job == "Screenplay" {
			jobs = append(jobs, c)
		}
	}

	return printMovie(out, movie, jobs)
}

func printMovie(out io.Writer, movie movie, details []crew) error {
	t := time.Now()

	jobs := []string{"Director", "Director of Photography", "Screenplay"}

	res := []string{}

	f := func(s string) string {
		str := ""
		for _, v := range details {
			if v.Job == s {
				str = str + v.Name + ", "
			}
		}
		str = strings.TrimSuffix(str, ", ")
		return str
	}

	for _, j := range jobs {
		str := f(j)
		res = append(res, str)
	}

	dir := res[0]
	cine := res[1]
	sp := res[2]

	t1 := time.Since(t)
	fmt.Fprintf(out, "%s\n", t1.String())

	if len(sp) == 0 {
		_, err := fmt.Fprintf(out, "\nTitle:  %s\nDirector: %s\nCinematography: %s\n\nOverview: %s\n\nRating: %v\n",
			movie.Title, dir, cine, movie.Overview, movie.VoteAverage)
		return err
	}

	_, err := fmt.Fprintf(out, "\nTitle:  %s\nDirector: %s\nCinematography: %s\nScreenplay: %s\n\nOverview: %s\n\nRating: %v\n",
		movie.Title, dir, cine, sp, movie.Overview, movie.VoteAverage)

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
