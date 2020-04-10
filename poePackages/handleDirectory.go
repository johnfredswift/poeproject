package poePackages

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ReadRawDirectory() ([][]string, [][]string, [][]string, [][]string) { //Returns Perm, HC Perm, Temp & HC Temp Leagues
	rawFilePath := `data/raw`

	files, readingDirErr := ioutil.ReadDir(rawFilePath)
	if readingDirErr != nil {
		panic(readingDirErr)
	}
	
	var permLeagues, hcPermLeagues [][]string //[0] Name, [1] StartDate, [2] EndDate //Date Format YYYY-MM-DD
	var tempLeagues, hcTempLeagues [][]string
	for _, f := range files {
		temp := append(strings.Split(f.Name(), "."), f.Name())
		name := temp[0]
		if name == "Standard" {
			permLeagues = append(permLeagues, temp)
		} else if name == "Hardcore" {
			hcPermLeagues = append(hcPermLeagues, temp)
		} else if strings.Contains(name, "Hardcore") {
			hcTempLeagues = append(hcTempLeagues, temp)
		} else {
			tempLeagues = append(tempLeagues, temp)
		}
	}
	//Leagues -> LeagueType -> Name, Start, Finish, Dir Name
	return permLeagues, hcPermLeagues, tempLeagues, hcTempLeagues
}
func UpdateAllLeagueDirectories(permLeagues [][]string, hcPermLeagues [][]string, tempLeagues [][]string, hcTempLeagues [][]string) {
	err := LeagueFiles(permLeagues[0][0], permLeagues[0][1], permLeagues[0][2])
	if err != nil {
		panic(err)
	}
	err = LeagueFiles(hcPermLeagues[0][0], hcPermLeagues[0][1], hcPermLeagues[0][2])
	if err != nil {
		panic(err)
	}
	err = LeagueFiles(permLeagues[len(permLeagues)-1][0],
		permLeagues[len(permLeagues)-1][1],
		permLeagues[len(permLeagues)-1][2])
	if err != nil {
		panic(err)
	}
	err = LeagueFiles(hcPermLeagues[len(hcPermLeagues)-1][0],
		hcPermLeagues[len(hcPermLeagues)-1][1],
		hcPermLeagues[len(hcPermLeagues)-1][2])
	if err != nil {
		panic(err)
	}
	for _, l := range tempLeagues {
		err := LeagueFiles(l[0], l[1], l[2])
		if err != nil {
			panic(err)
		}
	}
	for _, l := range hcTempLeagues {
		err := LeagueFiles(l[0], l[1], l[2])
		if err != nil {
			panic(err)
		}
	}
}
func UpdateMultipleLeaguesInDirectoriesContents(leagues [][]string) {
	log.Println("Started Updating Leagues In Directories Content")
	for _, league := range leagues {
		UpdateLeaguesInDirectoriesContents(league)
	}
	log.Println("Finished Updating Leagues In Directories Content")
}
func UpdateLeaguesInDirectoriesContents(league []string) {
	log.Println("UpdatingDirectories Contents for :", league[0])
	rawPath := `data/raw/` + league[3]
	leagueDataPath := `data/leagueData/` + league[0]
	var files []string
	err := filepath.Walk(rawPath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for i, srcFile := range files {
		if i != 0 {
			fileName := strings.Split(srcFile, `\`)[3]
			dstFile := leagueDataPath + `/` + fileName
			//fmt.Println(srcFile)
			//fmt.Println(dstFile)
			err = CopyFile(srcFile, dstFile)
			if err != nil {
				panic(err)
			}
		}
	}

}
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
func LeagueFiles(leagueName string, startDate string, endDate string) error {
	dirPath := `data/leagueData/` + leagueName
	dirExists, dirExistsErr := PathExists(dirPath)
	if dirExistsErr != nil {
		panic(dirExistsErr)
	}
	if !dirExists {
		dirCreateErr := os.MkdirAll(dirPath, os.ModePerm)
		if dirCreateErr != nil {
			panic(dirCreateErr)
		}
	}
	filePath := `data/leagueData/` + leagueName + "/notes.txt"
	fileExists, fileExistsErr := PathExists(filePath)
	if fileExistsErr != nil {
		panic(fileExistsErr)
	}
	if !fileExists {
		f, fileExistsErr := os.Create(filePath)
		if fileExistsErr != nil {
			panic(fileExistsErr)
		}
		_, err := f.WriteString(
			"League Name: " + leagueName +
				"\nStart Date: " + startDate +
				"\nEnd Date: " + endDate)
		if err != nil {
			panic(err)
		}
		defer f.Close()
	} else {
		text, fileReadErr := ioutil.ReadFile(filePath)
		if fileReadErr != nil {
			panic(fileReadErr)
		}
		lines := strings.Split(string(text), "\n")
		lines[2] = "End Date: " + endDate
		output := strings.Join(lines, "\n")
		fileWriteErr := ioutil.WriteFile(filePath, []byte(output), 0666)
		if fileWriteErr != nil {
			panic(fileWriteErr)
		}
	}
	return nil
}
func deleteNotes(path string) {
	_ = os.Remove(path)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
