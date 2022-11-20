/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// movieCmd represents the movie command
var movieCmd = &cobra.Command{
	Use:   "movie",
	Short: "Search for movie <movie>",
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

	return printMovie(out, movie)
}

func printMovie(out io.Writer, movie movie) error {
	_, err := fmt.Fprintf(out, "\nTitle:  %s\n\nOverview: %s\n\nRating: %v\n", movie.Title, movie.Overview, movie.VoteAverage)

	return err
}

func init() {
	rootCmd.AddCommand(movieCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// movieCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// movieCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
