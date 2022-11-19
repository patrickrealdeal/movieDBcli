/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// movieCmd represents the movie command
var movieCmd = &cobra.Command{
	Use:   "movie",
	Short: "Search for movie <movie>",
	RunE: func(cmd *cobra.Command, args []string) error {
		// apiRoot := viper.GetString("api-root")

		return listMovie(os.Stdout, APIURL)
	},
}

func listMovie(out io.Writer, apiRoot string) error {
	movies, err := get(apiRoot)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(out, movies[0].Title)

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
