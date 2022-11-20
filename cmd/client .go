package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const APIKEY = "86df293aaead73693b8da7fd28b3549e"

// const APIURL = "https://api.themoviedb.org/3/search/movie?api_key=" + APIKEY + "&query=" + m + "&language=en-US&page=1"
const APIROOT = "https://api.themoviedb.org/3"

var (
	ErrConnection      = errors.New("Connection error")
	ErrNotFound        = errors.New("Not found")
	ErrInvalidResponse = errors.New("Invalid server response")
	ErrInvalid         = errors.New("Invalid data")
	ErrNotNumber       = errors.New("Not a number")
)

type movie struct {
	MovieID     int     `json:"id"`
	Title       string  `json:"title"`
	ReleaseDate string  `json:"release_date"`
	Language    string  `json:"original_language"`
	Adult       bool    `json:"adult"`
	Image       string  `json:"poster_path"`
	Overview    string  `json:"overview"`
	VoteAverage float32 `json:"vote_average"`
}

type movieList struct {
	List []movie `json:"results"`
}

func newClient() *http.Client {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}

	return c
}

func getMovies(url string) ([]movie, error) {
	r, err := newClient().Get(url)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnection, err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, fmt.Errorf("Cannot read body: %w", err)
		}
		err = ErrInvalidResponse

		if r.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}

		return nil, fmt.Errorf("%w: %s", err, msg)
	}

	var resp movieList

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return nil, err
	}

	if len(resp.List) == 0 {
		return nil, fmt.Errorf("%w: no results found", ErrNotFound)
	}

	return resp.List, nil
}

func getMovie(apiRoot, req string) (movie, error) {
	u := fmt.Sprintf("%s/search/movie?api_key=%s&query=%s&language=en-US&page=1", apiRoot, APIKEY, req)

	movies, err := getMovies(u)
	if err != nil {
		return movie{}, err
	}

	if len(movies) < 1 {
		return movie{}, fmt.Errorf("%w: Invalid results", ErrInvalid)
	}

	return movies[0], nil
}
