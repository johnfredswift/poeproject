package main

import (
	poePackage "PoENinjaData/poePackages"
	"testing"
	"os"
)

func TestReadRawDirectory(t *testing.T) {
	_, _, testLeagueData, _ := poePackage.ReadRawDirectory()
	if string(testLeagueData[0][0][0]) == ""{
		t.Error("No League Files Found ,expected test files to be filled, Please ensure data files are downloaded and correctly stored")
	}
}

func TestCreateGraphValueOverTimeForOneTimeForOneTimeContinuous(t *testing.T){
	_, _, testLeagueData, _ := poePackage.ReadRawDirectory()
	createGraphValueFromLeagueStartForOneItem("TheDoctor", testLeagueData, "TEST")

	fileStats, err := os.Stat("data/graphs/TESTThe Doctor.png")
	if err == nil {
	  } else {
		t.Error("File Not Created, Please ensure file structure is correctly made") 
	}

	if fileStats.Size() < 1{
		t.Error("File Size is too small, Assumed the graph was not created, please ensure league data is correctly downloaded and located")
	}

}