package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

//	Variables

var artist, err = getData()
var search = searchBar()
var locationFilter = filterLocation()
var copyArtist = artist

// Refresh de l'API toutes les cinq minutes

func RefreshingApi() {
	for {
		time.Sleep(5 * time.Minute)
		artist, err = getData()
		search = searchBar()
		locationFilter = filterLocation()
		copyArtist = artist
		fmt.Println("Api Refreshed")
	}

}

func redoWhenError() {
	artist, err = getData()
	search = searchBar()
	locationFilter = filterLocation()
	copyArtist = artist
}

// Structure pour API

type group struct {
	Artists   string
	Locations string
	Dates     string
	Relation  string
}

type artists struct {
	Id           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	Locations    string
	ConcertDates string
	Relations    string
}

type relation struct {
	Id             int
	DatesLocations map[string][]string
}

type location struct {
	Id        int
	Locations []string
	Dates     string
}

type dates struct {
	Id    int
	Dates []string
}

type artistFull struct {
	Id           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	Locations    []string
	ConcertDates []string
	Relations    map[string][]string
}

// Structure database pour searchbar

type SB struct {
	ArtistBandName  []string
	Members         []string
	Location        []string
	FirstAlbumDate  []string
	CreationDate    []int
	MinFirstAlbum   int
	MaxFirstAlbum   int
	MinCreationDate int
	MaxCreationDate int
	MaxMembers      []int
}

// Analyser données API

func readUrl(url string) ([]byte, bool) {
	res, err1 := http.Get(url)
	if err1 != nil {
		fmt.Println("Erreur URL")
		return nil, false
	}
	defer res.Body.Close()
	body, err1 := ioutil.ReadAll(res.Body)
	if err1 != nil {
		fmt.Println("Erreur URL")
		return nil, false
	}
	return body, true
}

func getLink(url string) (group, bool) {
	body, truth := readUrl(url)
	if truth == false {
		return group{}, false
	}
	links := group{}
	jsonUrl := json.Unmarshal(body, &links)
	if jsonUrl != nil {
		fmt.Println("Marche pas")
		return group{}, false
	}
	return links, true
}

func getArtists(url string) ([]artists, bool) {
	body, truth := readUrl(url)
	if truth == false {
		return nil, false
	}
	var art []artists
	jsonUrl := json.Unmarshal(body, &art)
	if jsonUrl != nil {
		fmt.Println("Marche pas")
		return nil, false
	}
	return art, true
}

func getRelation(url string) (relation, bool) {
	body, truth := readUrl(url)
	if truth == false {
		return relation{}, false
	}
	relations := relation{}
	jsonUrl := json.Unmarshal(body, &relations)
	if jsonUrl != nil {
		fmt.Println("Marche pas")
		return relation{}, false
	}
	return relations, true
}

func getLocations(url string) (location, bool) {
	body, truth := readUrl(url)
	if truth == false {
		return location{}, false
	}
	locations := location{}
	jsonUrl := json.Unmarshal(body, &locations)
	if jsonUrl != nil {
		fmt.Println("Marche pas")
		return location{}, false
	}
	return locations, true
}

func getDates(url string) (dates, bool) {
	body, truth := readUrl(url)
	if truth == false {
		return dates{}, false
	}
	date := dates{}
	jsonUrl := json.Unmarshal(body, &date)
	if jsonUrl != nil {
		fmt.Println("Marche pas")
		return dates{}, false
	}
	return date, true
}

func getData() ([]artistFull, bool) {
	category, truth := getLink("https://groupietrackers.herokuapp.com/api")
	if truth == false {
		return nil, false
	}
	artis, truth2 := getArtists(category.Artists)
	if truth2 == false {
		return nil, false
	}
	var listArtistComplete []artistFull
	for i := range artis {
		l, err1 := getLocations(artis[i].Locations)
		d, err2 := getDates(artis[i].ConcertDates)
		r, err3 := getRelation(artis[i].Relations)
		if err1 == false || err2 == false || err3 == false {
			return nil, false
		}
		var artistComplete artistFull
		artistComplete.Id = artis[i].Id
		artistComplete.Image = artis[i].Image
		artistComplete.Name = artis[i].Name
		artistComplete.Members = artis[i].Members
		artistComplete.CreationDate = artis[i].CreationDate
		artistComplete.FirstAlbum = artis[i].FirstAlbum
		artistComplete.Locations = l.Locations
		artistComplete.ConcertDates = d.Dates
		artistComplete.Relations = r.DatesLocations
		listArtistComplete = append(listArtistComplete, artistComplete)
	}
	return listArtistComplete, true
}

// Fonction pour les filtres

func filters(w http.ResponseWriter, r *http.Request) {
	copyArtist = nil
	isFilterMember := false

	for i := 1; i <= 7; i++ {
		if r.FormValue(strconv.Itoa(i)) != "" {
			isFilterMember = true
			break
		}
	}
	for i := 0; i < len(artist); i++ {
		isLocation := false
		if r.FormValue("locationFilter") != "" {
			for _, loc := range artist[i].Locations {
				if strings.Contains(loc, r.FormValue("locationFilter")) {
					isLocation = true
					break
				}
			}
		}
		creaDate, _ := strconv.Atoi(r.FormValue("CreationDate"))
		member := strconv.Itoa(len(artist[i].Members))
		firstAlbum := artist[i].FirstAlbum[6:]
		checkboxMember := r.FormValue(member)

		//debug
		fmt.Printf("creaDate: %v\n", creaDate)
		fmt.Printf("CreationDate from artist: %v\n", artist[i].CreationDate)
		fmt.Printf("FirstAlbum: %v\n", firstAlbum)

		//debug pour condtions
		fmt.Printf("isFilterMember: %v\n", isFilterMember)
		fmt.Printf("checkboxMember: %v\n", checkboxMember)
		fmt.Printf("isLocation: %v\n", isLocation)
		fmt.Printf("r.FormValue(\"FirstAlbum\"): %v\n", r.FormValue("FirstAlbum"))


		if artist[i].CreationDate >= creaDate && firstAlbum >= r.FormValue("FirstAlbum") {
			if isFilterMember == true && checkboxMember != "" || !isFilterMember {
				if isLocation || r.FormValue("locationFilter") == "" {
					copyArtist = append(copyArtist, artist[i])
				}
			}
		}
	}
}

func filterLocation() []string {
	locations := search.Location
	var locationsFinal []string
	for i := 0; i < len(locations); i++ {
		index := strings.Index(locations[i], "-")
		country := locations[i][index+1:]
		if !AlreadyInSlice(locationsFinal, country, 0, []int{}) {
			locationsFinal = append(locationsFinal, country)
		}
	}
	sort.Slice(locationsFinal, func(i, j int) bool {
		return locationsFinal[i] < locationsFinal[j]
	})
	return locationsFinal
}

// Fonction searchbar

func searchBar() SB {
	var sortArtistData SB
	maxMembers := 0
	if !err || artist == nil {
		return SB{}
	}
	for i := 0; i < len(artist); i++ {
		sortArtistData.ArtistBandName = append(sortArtistData.ArtistBandName, artist[i].Name)
		for j := 0; j < len(artist[i].Members); j++ {
			if !AlreadyInSlice(sortArtistData.Members, artist[i].Members[j], 0, []int{}) {
				sortArtistData.Members = append(sortArtistData.Members, artist[i].Members[j])
			}
		}
		if len(artist[i].Members) > maxMembers {
			maxMembers = len(artist[i].Members)
		}
		for j := 0; j < len(artist[i].Locations); j++ {
			if !AlreadyInSlice(sortArtistData.Location, artist[i].Locations[j], 0, []int{}) {
				sortArtistData.Location = append(sortArtistData.Location, artist[i].Locations[j])
			}
		}
		sortArtistData.FirstAlbumDate = append(sortArtistData.FirstAlbumDate, artist[i].FirstAlbum)
		if i == 0 {
			sortArtistData.MinFirstAlbum, _ = strconv.Atoi(artist[i].FirstAlbum[6:])
			sortArtistData.MaxFirstAlbum, _ = strconv.Atoi(artist[i].FirstAlbum[6:])
		} else {
			date, _ := strconv.Atoi(artist[i].FirstAlbum[6:])
			if sortArtistData.MinFirstAlbum > date {
				sortArtistData.MinFirstAlbum = date
			} else if sortArtistData.MaxFirstAlbum < date {
				sortArtistData.MaxFirstAlbum = date
			}
		}
		if !AlreadyInSlice(nil, "", artist[i].CreationDate, sortArtistData.CreationDate) {
			sortArtistData.CreationDate = append(sortArtistData.CreationDate, artist[i].CreationDate)
		}
	}
	for i := 1; i < 8; i++ {
		sortArtistData.MaxMembers = append(sortArtistData.MaxMembers, i)
	}
	sort.Slice(sortArtistData.Members, func(i, j int) bool {
		return sortArtistData.Members[i] < sortArtistData.Members[j]
	})
	sort.Slice(sortArtistData.ArtistBandName, func(i, j int) bool {
		return sortArtistData.ArtistBandName[i] < sortArtistData.ArtistBandName[j]
	})
	sort.Slice(sortArtistData.FirstAlbumDate, func(i, j int) bool {
		return sortArtistData.FirstAlbumDate[i] < sortArtistData.FirstAlbumDate[j]
	})
	sort.Slice(sortArtistData.Location, func(i, j int) bool {
		return sortArtistData.Location[i] < sortArtistData.Location[j]
	})
	sort.Slice(sortArtistData.CreationDate, func(i, j int) bool {
		return sortArtistData.CreationDate[i] < sortArtistData.CreationDate[j]
	})
	sortArtistData.MinCreationDate, sortArtistData.MaxCreationDate = sortArtistData.CreationDate[0], sortArtistData.CreationDate[len(sortArtistData.CreationDate)-1]

	return sortArtistData
}

func searchBarCalculation(w http.ResponseWriter, r *http.Request) {
	key := strings.ToLower(r.FormValue("data"))
	copyArtist = nil
	if strings.Contains(key, "artist/groupe") {
		key = strings.Join(strings.Split(key, " - artist/groupe"), "")
		for i := 0; i < len(artist); i++ {
			if key == strings.ToLower(artist[i].Name) {
				copyArtist = append(copyArtist, artist[i])
			}
		}
	} else if strings.Contains(key, "membre") {
		key = strings.Join(strings.Split(key, " - membre"), "")
		for i := 0; i < len(artist); i++ {
			for j := range artist[i].Members {
				if key == strings.ToLower(artist[i].Members[j]) {
					copyArtist = append(copyArtist, artist[i])
				}
			}
		}
	} else if strings.Contains(key, "location") {
		key = strings.Join(strings.Split(key, " - location"), "")
		for i := 0; i < len(artist); i++ {
			for j := range artist[i].Locations {
				if key == strings.ToLower(artist[i].Locations[j]) {
					copyArtist = append(copyArtist, artist[i])
				}
			}
		}
	} else if strings.Contains(key, "date de création") {
		key = strings.Join(strings.Split(key, " - date de création"), "")
		for i := 0; i < len(artist); i++ {
			if key == strconv.Itoa(artist[i].CreationDate) {
				copyArtist = append(copyArtist, artist[i])
			}
		}
	} else if strings.Contains(key, "date du premier album") {
		key = strings.Join(strings.Split(key, " - date du premier album"), "")
		for i := 0; i < len(artist); i++ {
			if key == artist[i].FirstAlbum {
				copyArtist = append(copyArtist, artist[i])
			}
		}
	} else if key == "" {
		copyArtist = artist
	} else {
		for i := range artist {
			if strings.Contains(strings.ToLower(artist[i].Name), key) || key == strconv.Itoa(artist[i].CreationDate) || key == artist[i].FirstAlbum {
				copyArtist = append(copyArtist, artist[i])
				continue
			} else {
				truth := false
				for j := range artist[i].Members {
					if strings.Contains(strings.ToLower(artist[i].Members[j]), key) {
						copyArtist = append(copyArtist, artist[i])
						truth = true
						break
					}
				}
				if !truth {
					for k := range artist[i].Locations {
						if strings.Contains(strings.ToLower(artist[i].Locations[k]), key) {
							copyArtist = append(copyArtist, artist[i])
							break
						}
					}
				}

			}
		}
	}
}

func AlreadyInSlice(sliceString []string, wordCompared string, intCompared int, sliceInt []int) bool {
	if sliceString == nil {
		for already := range sliceInt {
			if intCompared == sliceInt[already] {
				return true
			}
		}
	} else {
		for already := range sliceString {
			if wordCompared == sliceString[already] {
				return true
			}
		}
	}
	return false
}

// Fonction pour les pages


func homePage(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("CreationDate") != "" {
		filters(w, r)
	}
	if r.FormValue("data") != "" {
		searchBarCalculation(w, r)
	}
	if err == false {
		redoWhenError()
		tmpl := template.Must(template.ParseFiles("../templates/500.gohtml"))
		_ = tmpl.Execute(w, nil)
	} else if r.URL.Path != "/" {
		tmpl := template.Must(template.ParseFiles("../templates/404.gohtml"))
		_ = tmpl.Execute(w, struct {
			Search SB
		}{Search: search})
	} else {
		tmpl := template.Must(template.ParseFiles("../templates/index.gohtml"))
		_ = tmpl.Execute(w, struct {
			Artist   []artistFull
			Search   SB
			Location []string
		}{Artist: copyArtist, Search: search, Location: locationFilter})
		copyArtist = artist
	}
}

func artistPage(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	idString := url[8:]
	id, _ := strconv.Atoi(idString)

	if id-1 >= len(artist) || id-1 < 0 {
		tmpl := template.Must(template.ParseFiles("../templates/404.gohtml"))
		_ = tmpl.Execute(w, nil)
	} else {
		tmpl := template.Must(template.ParseFiles("../templates/artist.gohtml"))
		_ = tmpl.Execute(w, struct {
			Artist artistFull
			Search SB
		}{Artist: artist[id-1], Search: search})
	}
}

// Gestion requêtes

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/artist/", artistPage)
}

// Fonction main server

func main() {
	go RefreshingApi()
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets"))))
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("../templates"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("../images"))))
	handleRequests()

	fmt.Println("Server started \"http://localhost:8080\"")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
