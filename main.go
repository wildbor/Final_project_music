package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
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
	e.GET("/player/lyric", GetLyric)
	e.POST("/player", CreatePlayerListController)
	e.PUT("/player", DeleteController)
	e.GET("/play", OpenLinkController)
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

func OpenLinkController(c echo.Context) error {
	xid, _ := strconv.Atoi(c.QueryParam("id"))
	if c.QueryParam("id") == "" {
		return c.JSON(http.StatusNotFound, "ID Not Found")
	} else {
		for _, xdata := range xVarPlayer {
			//xarray2 := strings.Fields(xdata)
			xarray2 := xdata.ID
			//fmt.Println("xarray: ", xarray2)
			//fmt.Println("xid: ", xid)
			if xarray2 == xid {
				xUrl := xdata.SongViewURL

				var args []string
				switch runtime.GOOS {
				case "darwin":
					args = []string{"open", xUrl}
				case "windows":
					args = []string{"cmd", "/c", "start", xUrl}
				default:
					args = []string{"xdg-open", xUrl}
				}
				cmd := exec.Command(args[0], args[1:]...)
				err := cmd.Run()
				if err != nil {
					log.Printf("openinbrowser: %v\n", err)
				}

			}

		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"messages": "Open selected track by ID",
		})
	}
}

//get all playerlist
//get playerlist by ID
//get lyric

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

//create new track on player
func CreatePlayerListController(c echo.Context) error {

	xPlayerBind := xPlayerStruct{}
	c.Bind(&xPlayerBind)
	var xStatus int = 0

	url := "https://itunes.apple.com/search?term=" + xPlayerBind.ArtistName
	req, _ := http.NewRequest("GET", url, nil)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	json.Unmarshal(body, &xVarArtist)

	for _, xdata := range xVarArtist.Results {
		//xarray2 := strings.Fields(xdata)
		xSongnameArr := xdata.SongName
		//fmt.Println("xarray: ", xarray2)
		//fmt.Println("xid: ", xid)
		if xSongnameArr == xPlayerBind.SongName {
			xStatus = 1

			//"messages": "Get selected " + xartistname + " song",
			//"user":     xdata,
		}
	}
	if xStatus == 1 {
		//fmt.Println("ADA")
		if len(xVarPlayer) == 0 {
			xPlayerBind.ID = 1
		} else {
			newID := xVarPlayer[len(xVarPlayer)-1].ID + 1
			xPlayerBind.ID = newID
		}
		xVarPlayer = append(xVarPlayer, xPlayerBind)

	} else {
		//fmt.Println("GAK ADA")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"messages": "Failed create user, because track name can't found on iTunes",
			"user":     xPlayerBind,
			//"total":    len(xUsers),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"messages": "success create user",
		"user":     xPlayerBind,
		//"total":    len(xUsers),
	})

}

func GetLyric(c echo.Context) error {
	xid, _ := strconv.Atoi(c.QueryParam("id"))

	//xUsers[xid-1].UserName = xuserbind.UserName
	for _, xdata := range xVarPlayer {
		//xarray2 := strings.Fields(xdata)
		xIDPlayer := xdata.ID

		//fmt.Println("xarray: ", xarray2)
		//fmt.Println("xid: ", xid)
		if xIDPlayer == xid {
			xArtisPlayer := xdata.ArtistName
			xTrackPlayer := xdata.SongName

			url := "https://api.lyrics.ovh/v1/" + xArtisPlayer + "/" + xTrackPlayer

			req, _ := http.NewRequest("GET", url, nil)

			res, _ := http.DefaultClient.Do(req)

			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(body, &xVarLyric)
			//fmt.Println(res)
			//fmt.Println(string(body))

			return c.JSON(http.StatusOK, map[string]interface{}{
				"messages": "Lyric found",
				"lyric":    xVarLyric,
			})
		}

	}
	return c.JSON(http.StatusNotFound, "Lyric didn't found")

}

func DeleteController(c echo.Context) error {
	xid, _ := strconv.Atoi(c.QueryParam("id"))

	if xid == len(xVarPlayer) {
		xVarPlayer[xid-1] = xPlayerStruct{}
		xVarPlayer = xVarPlayer[:len(xVarPlayer)-1]
	} else {
		xVarPlayer[xid-1] = xVarPlayer[len(xVarPlayer)-1]
		xVarPlayer[len(xVarPlayer)-1] = xPlayerStruct{}
		xVarPlayer = xVarPlayer[:len(xVarPlayer)-1]
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"messages": "Data already deleted",
	})
}
