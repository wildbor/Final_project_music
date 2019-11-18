package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type xArtistStruct struct {
	ResultCount int `json:"resultCount"`
	Results     []struct {
		SongID      int    `json:"trackId"`
		ArtistName  string `json:"artistName"`
		AlbumName   string `json:"collectionName,omitempty"`
		SongName    string `json:"trackName"`
		SongViewURL string `json:"trackViewUrl"`
	} `json:"results"`
}

type xPlayerStruct struct {
	ID          int    `json:"id" form:"id"`
	SongID      int    `json:"trackId"`
	ArtistName  string `json:"artistName"`
	AlbumName   string `json:"collectionName,omitempty"`
	SongName    string `json:"trackName"`
	SongViewURL string `json:"trackViewUrl"`
}

type xLyricStruct struct {
	Lyric string `json:"lyrics"`
}

var xVarArtist xArtistStruct
var xVarPlayer []xPlayerStruct
var xVarLyric xLyricStruct

func main() {

	e := echo.New()
	e.GET("/track", GetTrack)
	e.GET("/player", FilterPlayerController)
	e.Logger.Fatal(e.Start(":8080"))
}

// get track by artist name from itunes
// get track by artist name & track name from itunes
func GetTrack(c echo.Context) error {
	xartistname := c.QueryParam("artistname")
	xsongname := c.QueryParam("songname")
	url := "https://itunes.apple.com/search?term=" + xartistname
	req, _ := http.NewRequest("GET", url, nil)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	json.Unmarshal(body, &xVarArtist)
	if xsongname == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"messages": "Success get all " + xartistname + " songs",
			"users":    xVarArtist,
		})
	} else {
		for _, xdata := range xVarArtist.Results {
			//xarray2 := strings.Fields(xdata)
			xSongnameArr := xdata.SongName
			//fmt.Println("xarray: ", xarray2)
			//fmt.Println("xid: ", xid)
			if xSongnameArr == xsongname {
				return c.JSON(http.StatusOK, map[string]interface{}{
					"messages": "Get selected " + xartistname + " song",
					"user":     xdata,
				})
			}

		}
	}
	return c.JSON(http.StatusNotFound, "Song not found")
}

func FilterPlayerController(c echo.Context) error {

	xid, _ := strconv.Atoi(c.QueryParam("id"))

	xPlayerbind := xPlayerStruct{}
	c.Bind(&xPlayerbind)

	if c.QueryParam("id") == "" {
		if len(xVarPlayer) > 0 {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"messages": "Success get all playerlist",
				"users":    xVarPlayer,
			})
		} else {
			return c.JSON(http.StatusNotFound, "Playerlist blank")
		}
	} else {
		for _, xdata := range xVarPlayer {
			//xarray2 := strings.Fields(xdata)
			xarray2 := xdata.ID
			//fmt.Println("xarray: ", xarray2)
			//fmt.Println("xid: ", xid)
			if xarray2 == xid {
				return c.JSON(http.StatusOK, map[string]interface{}{
					"messages": "Get selected track by ID",
					"user":     xdata,
				})
			}
		}

		return c.JSON(http.StatusNotFound, "Data not found")
	}

}
