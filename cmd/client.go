package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
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

type cast struct {
	CastID     int    `json:"id"`
	Department string `json:"known_for_department"`
	Name       string `json:"name"`
}

type crew struct {
	CrewID             int    `json:"id"`
	KnownForDepartment string `json:"known_for_department"`
	Department         string `json:"department"`
	Job                string `json:"job"`
	Name               string `json:"name"`
}

type credits struct {
	Cast []cast `json:"cast"`
	Crew []crew `json:"crew"`
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

func getAll(url string) ([]movie, error) {
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

	movies, err := getAll(u)
	if err != nil {
		return movie{}, err
	}

	if len(movies) < 1 {
		return movie{}, fmt.Errorf("%w: Invalid results", ErrInvalid)
	}

	return movies[0], nil
}

func getDetails(movieID int) (credits, error) {
	u := fmt.Sprintf("%s/movie/%d/credits?api_key=%s&language=en-US&page=1", APIROOT, movieID, APIKEY)

	r, err := newClient().Get(u)
	if err != nil {
		return credits{}, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(r.Body)
		if err != nil {
			return credits{}, fmt.Errorf("Cannot read body: %w", err)
		}
		err = ErrInvalidResponse

		if r.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}

		return credits{}, fmt.Errorf("%w: %s", err, msg)
	}

	var resp credits

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return credits{}, err
	}

	file, _ := json.MarshalIndent(resp, "", "")
	os.WriteFile("obj.json", file, 0644)

	if len(resp.Cast) == 0 && len(resp.Crew) == 0 {
		return credits{}, fmt.Errorf("%w: no results found", ErrNotFound)
	}

	return resp, nil
}
