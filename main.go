package main

import (
	poePackage "PoENinjaData/poePackages"
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wcharczuk/go-chart"
)

func main() {

	
	updateCmd := flag.NewFlagSet("update", flag.PanicOnError)
	continiousCmd := flag.NewFlagSet("cont", flag.PanicOnError)
	continiousHardcore := continiousCmd.Bool("hc", false, "If league is hardcore or not")
	continiousStandard := continiousCmd.Bool("st", false, "If league is standard or not")
	singleCmd := flag.NewFlagSet("sing", flag.PanicOnError)
	previousLeagueJSONCmd := flag.NewFlagSet("json", flag.PanicOnError)

	if len(os.Args) < 2 {
        fmt.Println("expected subcommand")
        os.Exit(1)
    }
	switch os.Args[1] {
	case "update":
		updateCmd.Parse(os.Args[2:])
		update()
	case "cont":
		continiousCmd.Parse(os.Args[2:])
		fmt.Println(continiousCmd.Args())
		if len(continiousCmd.Args()) != 1{
			log.Println("Continous subcommand must contain at one league")
			panic("Continous subcommand must contain at one league")
		}
		itemName := continiousCmd.Arg(0)
		var leagueData [][]string
		if *continiousHardcore && *continiousStandard{
			fmt.Println("Hc & St")
			_, _, leagueData, _ = poePackage.ReadRawDirectory()
			createGraphValueOverTimeForOneTimeForOneTimeContinuous(itemName, leagueData, "HCST"+"ContiniousGraph" )
		}
		if !*continiousHardcore && *continiousStandard{
			leagueData, _, _, _ = poePackage.ReadRawDirectory()
			createGraphValueOverTimeForOneTimeForOneTimeContinuous(itemName, leagueData, "SCST"+"ContiniousGraph" )
			fmt.Println("Not Hc & St")
		}
		if *continiousHardcore && !*continiousStandard{
			_, _, _, leagueData = poePackage.ReadRawDirectory()
			createGraphValueOverTimeForOneTimeForOneTimeContinuous(itemName, leagueData, "HCTemp"+"ContiniousGraph" )
			fmt.Println("Hc & Not St")
		}
		if !*continiousHardcore && !*continiousStandard{
			_, leagueData, _, _ = poePackage.ReadRawDirectory()
			createGraphValueOverTimeForOneTimeForOneTimeContinuous(itemName, leagueData, "SCTemp"+"ContiniousGraph" )
			fmt.Println("Not Hc & Not St")
		}
	case "sing":
		singleCmd.Parse(os.Args[2:])
		if len(singleCmd.Args()) < 2{
			log.Println("Single subcommand must contain at least two variables, item and league")
			panic("Single subcommand must contain at least two variables, item and league")
		}

		itemName := singleCmd.Arg(0)
		leagueNames := singleCmd.Args()
		leagueNames = leagueNames[1:]
		var leagueData [][]string
		var leagueDataHC [][]string

		_, _, leagueData, leagueDataHC = poePackage.ReadRawDirectory()
		fmt.Println("Softcore League")
		var leagueDataList [][]string
		for a := 0; a < len(leagueData); a++{
			for b := 0; b < len(leagueNames); b++{
				if leagueData[a][0] == string(leagueNames[b]){
					leagueDataList = append(leagueDataList, leagueData[a])
				}
			}
		}
		for a := 0; a < len(leagueDataHC); a++{
			for b := 0; b < len(leagueNames); b++{
				if leagueDataHC[a][0] == string(leagueNames[b]){
					leagueDataList = append(leagueDataList, leagueDataHC[a])
				}
			}
		}
		tempName := ""
		for i := 0; i < len(leagueNames); i++{
			tempName += string(leagueNames[i])
		}
		createGraphValueFromLeagueStartForOneItem(itemName, leagueDataList, itemName + tempName)
		
	
	case "json":
		previousLeagueJSONCmd.Parse(os.Args[2:])
		JSONFile, err := ioutil.ReadFile("data/previousLeagues.json")
		check(err)
		fmt.Println(string(JSONFile))	
	}
}
func createGraphValueOverTimeForOneTimeForOneTimeContinuous(searchItemName string, leagueData [][]string, fileName string) {

	var valuesSlices [][]string
	dayCounter := 0

	var temp []chart.ContinuousSeries
	for leagueCounter := 0; leagueCounter < len(leagueData); leagueCounter++ {
		tempCSVCurrencyStruct := getItemsOfLeague(leagueData[leagueCounter])
		valuesSlices = append(valuesSlices, getValuesOfItemFromLeague(tempCSVCurrencyStruct, searchItemName))
		tempChart := chart.ContinuousSeries{}

		for valueCounter := 0; valueCounter < len(valuesSlices[leagueCounter]); valueCounter++ {
			s, _ := strconv.ParseFloat(valuesSlices[leagueCounter][valueCounter], 64)
			tempChart.YValues = append(tempChart.YValues, s)
			tempChart.XValues = append(tempChart.XValues, float64(dayCounter+valueCounter))
			dayCounter += valueCounter
		}
		temp = append(temp, tempChart)
	}

	graph := chart.Chart{Series: []chart.Series{}}
	for i := 0; i < len(temp); i++ {
		graph.Series = append(graph.Series, temp[i])
	}
	f, _ := os.Create("data/graphs/" + fileName + searchItemName + ".png")
	defer f.Close()
	_ = graph.Render(chart.PNG, f)
}

func createGraphValueFromLeagueStartForOneItem(searchItemName string, leagueData [][]string, fileName string) {

	var valuesSlices [][]string

	for i := 0; i < len(leagueData); i++ {
		tempCSVCurrencyStruct := getItemsOfLeague(leagueData[i])
		valuesSlices = append(valuesSlices, getValuesOfItemFromLeague(tempCSVCurrencyStruct, searchItemName))

	}
	var doubleResults [][][]float64
	for leagueResultsCounter := 0; leagueResultsCounter < len(valuesSlices); leagueResultsCounter++ {
		var tempDoubleResults [][]float64
		for resultCounter := 0; resultCounter < len(valuesSlices[leagueResultsCounter]); resultCounter++ {
			s, _ := strconv.ParseFloat(valuesSlices[leagueResultsCounter][resultCounter], 64)
			tempDoubleResult := []float64{s, float64(resultCounter)}
			tempDoubleResults = append(tempDoubleResults, tempDoubleResult)
		}
		doubleResults = append(doubleResults, tempDoubleResults)
	}
	var temp []chart.ContinuousSeries

	for leagueCounter := 0; leagueCounter < len(doubleResults); leagueCounter++ {
		tempChart := chart.ContinuousSeries{}
		for valueCounter := 0; valueCounter < len(doubleResults[leagueCounter]); valueCounter++ {
			tempChart.YValues = append(tempChart.YValues, doubleResults[leagueCounter][valueCounter][0])
			tempChart.XValues = append(tempChart.XValues, doubleResults[leagueCounter][valueCounter][1])
		}
		temp = append(temp, tempChart)

	}

	graph := chart.Chart{Series: []chart.Series{}}
	for i := 0; i < len(temp); i++ {
		graph.Series = append(graph.Series, temp[i])
	}
	f, _ := os.Create("data/graphs/" + fileName + searchItemName + ".png")
	defer f.Close()
	_ = graph.Render(chart.PNG, f)
}

func getValuesOfItemFromLeague(items []CSVCurrencyStruct, searchItemName string) []string {
	var returnValues []string
	counter := 1
	for _, i := range items {
		if i.name == searchItemName {
			returnValues = append(returnValues, i.value)
			counter++
		}
	}
	return returnValues
}

func getItemsOfLeague(tempLeague []string) []CSVCurrencyStruct {
	rawPath := `data/leagueData/` + tempLeague[0] + `/` + tempLeague[3] + `.items.csv`
	csvFile, err := os.Open(rawPath)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	defer csvFile.Close()
	var returnStruct []CSVCurrencyStruct
	scanner := bufio.NewScanner(csvFile)
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		tempSlice := strings.Split(scanner.Text(), ";")

		var tempLineStruct CSVCurrencyStruct
		if len(tempSlice) == 9 {
			tempLineStruct = CSVCurrencyStruct{
				tempSlice[0], tempSlice[1], tempSlice[2], tempSlice[3], tempSlice[4],
				tempSlice[5], "", tempSlice[6], tempSlice[7], tempSlice[8]}

		} else if len(tempSlice) == 10 {
			tempLineStruct = CSVCurrencyStruct{
				tempSlice[0], tempSlice[1], tempSlice[2], tempSlice[3], tempSlice[4],
				tempSlice[5], tempSlice[6], tempSlice[7], tempSlice[8], tempSlice[9]}
		}
		returnStruct = append(returnStruct, tempLineStruct)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return returnStruct
}

func update() {
	updatePreviousLeagues()
	//fmt.Println(permLeagues, hcPermLeagues, tempLeagues, hcTempLeagues)
	permLeagues, hcPermLeagues, tempLeagues, hcTempLeagues := poePackage.ReadRawDirectory()
	poePackage.UpdateAllLeagueDirectories(permLeagues, hcPermLeagues, tempLeagues, hcTempLeagues)
	poePackage.UpdateMultipleLeaguesInDirectoriesContents(permLeagues)
	poePackage.UpdateMultipleLeaguesInDirectoriesContents(tempLeagues)
	poePackage.UpdateMultipleLeaguesInDirectoriesContents(hcPermLeagues)
	poePackage.UpdateMultipleLeaguesInDirectoriesContents(hcTempLeagues)
}
func updatePreviousLeagues() {
	//CurrentActiveLeagues := getCurrentlyActivePublicLeagues()
	//fmt.Println(CurrentActiveLeagues)
	previousLeaguesFile := loadPreviousLeaguesJSONFile()
	defer previousLeaguesFile.Close()
	currentLeagues := getCurrentlyActivePublicLeagues()
	fileLeagues := readLeaguesFromFile(previousLeaguesFile)
	updatedLeagues := PoELeagues{}

	for i := 0; i < len(currentLeagues); i++ {
		found := false
		for c := 0; c < len(fileLeagues); c++ {
			if currentLeagues[i].ID == fileLeagues[c].ID {
				found = true
			}
		}
		if !found {
			updatedLeagues = append(updatedLeagues, currentLeagues[i])
		}
	}
	updatedLeagues = append(fileLeagues, updatedLeagues...)
	writeJSONToFile(updatedLeagues, previousLeaguesFile)

}

func readLeaguesFromFile(previousLeaguesFile *os.File) (leagues PoELeagues) {
	data, err := ioutil.ReadAll(previousLeaguesFile)
	if err != nil {
		log.Panic(err)
	}
	fileReader := bytes.NewReader(data)
	leagues = unmarshalLeagues(fileReader)
	return leagues
}

func writeJSONToFile(data PoELeagues, file *os.File) {
	dataBytes, marshallingErr := json.Marshal(data)
	if marshallingErr != nil {
		panic(marshallingErr)
	}
	_, writingToFileErr := file.Write(dataBytes)
	if writingToFileErr != nil {
		panic(writingToFileErr)
	}
}
func unmarshalLeagues(data io.Reader) PoELeagues {
	result := PoELeagues{}
	jsonErr := json.NewDecoder(data).Decode(&result)
	if jsonErr != nil {
		log.Println(jsonErr, "File was either empty, or data was corrupted, Adding current leagues")
		return getCurrentlyActivePublicLeagues()
	}
	return result
}

func loadPreviousLeaguesJSONFile() *os.File {
	const fileLocation = `data/previousLeagues.json`
	//Checks for files existence
	if _, err := os.Stat(fileLocation); err == nil {
		log.Println("Exists")
	} else if os.IsNotExist(err) {
		f, createFileErr := os.Create(fileLocation)
		if createFileErr != nil {
			panic(createFileErr)
		} else {
			log.Println("File Created:", fileLocation)
		}
		return f
	} else { //Should the file exists but cannot be accessed for whatever reason
		panic(err)
	}
	f, err := os.OpenFile(fileLocation, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	} else {
		log.Println("File Opened:", fileLocation) //updatePreviousLeagues
	}
	return f
}
func stripeStandardLeaguesFromLeagues(leagues PoELeagues) PoELeagues {
	currentTempLeague := ""
	for _, i := range leagues {
		if !strings.Contains(i.ID, " ") {
			if !(i.ID == "Standard" || i.ID == "Hardcore") {
				currentTempLeague = i.ID
				break
			}
		}
	}
	var temp PoELeagues
	for _, i := range leagues {
		if strings.Contains(i.ID, currentTempLeague) {
			temp = append(temp, i)
		}
	}
	return temp
}

func getCurrentlyActivePublicLeagues() PoELeagues {
	url := "http://api.pathofexile.com/leagues"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	body := resp.Body
	return unmarshalLeagues(body)
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type PoELeagues []struct {
	ID          string        `json:"id"`
	Realm       string        `json:"realm"`
	Description string        `json:"description"`
	URL         string        `json:"url"`
	StartAt     time.Time     `json:"startAt"`
	EndAt       interface{}   `json:"endAt"`
	DelveEvent  bool          `json:"delveEvent"`
	Rules       []interface{} `json:"rules"`
	RegisterAt  time.Time     `json:"registerAt,omitempty"`
}

//League;Date;Id;Name;BaseType;Variant;Links;Value;Confidence
//League;Date;Id;Variant;Name;BaseType;;Links;Value;Confidence

type CSVCurrencyStruct struct {
	league     string
	date       string
	ID         string
	variant    string
	name       string
	baseType   string
	quality    string
	links      string
	value      string
	confidence string
}
