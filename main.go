// import the main package
package main

import (
    "bufio" // import bufio to scan in user input
    "fmt" // import the fmt package for printing
    "os"
    "KanjiFrequencyHelper/csvoperations"
    "KanjiFrequencyHelper/utils"
    "KanjiFrequencyHelper/kanji"
    "regexp"
    "sort"
    "strings"
    "sync"
)

type KeigoReadings struct {
    alreadyRead bool
    keigoMap map[string][]rune
    keigoenglishslice []string
    keigoromajislice []string
    keigojapaneseslice []string
    regex * regexp.Regexp
}

type Builder struct {
    strings.Builder
}

func (keigoOps* KeigoReadings) printmapkeigo(userinput string) {
    fmt.Println(userinput)
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Println(scanner.Text())
    userintputdashed := userinput + "-"

    if keigoOps.keigoMap[userinput] != nil {
        keigostring := string(keigoOps.keigoMap[userinput])
        keigostring = strings.ReplaceAll(keigostring, "*", "\n")
        fmt.Printf("\n%s: %s\n", keigostring, userinput)
    } else if keigoOps.keigoMap[userintputdashed] != nil {
        keigostring := string(keigoOps.keigoMap[userintputdashed])
        keigostring = strings.ReplaceAll(keigostring, "*", "\n")
        fmt.Printf("\n%s: %s\n", keigostring, userinput)
    } else {
        fmt.Printf("\n%s: DOES NOT EXIST\n", userinput)
    }
}


// create function to handle error
func handleError(err error, message string) {
    if err != nil {
        fmt.Println("Error: ", message)
        fmt.Println("Error: ", err)
    }
}

// Main function
func main() {

    // create kanji ops blank pointer
    kanjiOps := &kanji.KanjiReadings{}
    keigoOps := &KeigoReadings{
        alreadyRead: false,
    }

    // Create a scanner used to read user input/options
    scanner := bufio.NewScanner(os.Stdin)


    // Create a waitgroup
    var wg sync.WaitGroup

    // Create  string array of filepaths, this will help us save into proper map
    filePaths := []string {
        "./resources/KanjiFrequencyListOnyomi.csv",
        "./resources/KanjiFrequencyListKunyomi.csv",
        "./resources/KunyomiWithHiragana.csv",
        "./resources/KanjiMeanings.csv",
        "./resources/all_Readings_string.csv",
        "./resources/keigo_mapper.csv",
    }

    lenFiles := len(filePaths)

    wg.Add(lenFiles)

    // Iterate through all five
    for _, filePath := range filePaths {
        go func(filePath string) {
            defer wg.Done()
            kanjiOps.Lock()
            defer kanjiOps.Unlock()
            csvMap, err := csvoperations.ReadCSV(filePath)
            
            if err != nil {
                handleError(err, "Error with file path: " + filePath)
                return
            }

            switch filePath {
                case "./resources/KanjiFrequencyListOnyomi.csv":
                    kanjiOps.OnyomiMap = csvMap

                case "./resources/KanjiFrequencyListKunyomi.csv": 
                    kanjiOps.KunyomiMap = csvMap

                case "./resources/KunyomiWithHiragana.csv":
                    kanjiOps.KunyomiWithHiragana = csvMap

                case "./resources/KanjiMeanings.csv":
                    kanjiOps.KanjiMeanings = csvMap

                case "./resources/all_Readings_string.csv":
                    // setting the csvmap directly into the kanj.KanjiReadings map
                    kanjiOps.Readings = csvMap

                case "./resources/keigo_mapper.csv":
                    keigoOps.keigoMap = csvMap
            }

        }(filePath)
    }

    // Wait for all the go routines to finish, wait on all five files
    wg.Wait()

    kanjiOps.LoadFrequencies()

    // now back in sequential mode, we can 

    // Loop to keep the program running unless the user types in "exit"
    for {
        utils.ClearScreen() 
        fmt.Print("Select Function:\n1. Kanji Finder\n2. Keigo Finder\n3. Onyomi\n4. Kunyomi\n5. KunyomiWithHiragana\n6. Exit\nEnter Input: ")

        scanner.Scan()
        applicationSelector := scanner.Text()
        
        userInput := ""

        var Readings bool = true

        for userInput != "exit" {
            if applicationSelector == "1" {
                utils.ClearScreen()
                // create a bool to track Readings
                

                fmt.Println("KANJI ASSISTANT: Enter (hiragana, romaji, or katakana to get Readings")
                fmt.Println("Enter Input: ('exit' to quit, 'Readings' toggles verbosity: ")
                
                if Readings == true {
                    fmt.Println("Reading data enabled...")
                } else {
                    fmt.Println("Reading data silenced...")
                }

                scanner.Scan()
                userInput = scanner.Text()

                if userInput == "Readings" {
                    Readings = !Readings
                    fmt.Println("Reading data silenced...")
                    _ = bufio.NewScanner(os.Stdin)
                    continue
                }


                // Send each string into the PrintMap
                if kanjiOps.OnyomiMap != nil {
                    kanjiOps.PrintMap("Onyomi", kanjiOps.OnyomiMap[userInput], userInput, Readings)
                }

                if kanjiOps.KunyomiMap != nil {
                    kanjiOps.PrintMap("Kunyomi", kanjiOps.KunyomiMap[userInput], userInput, Readings)
                }

                if kanjiOps.KunyomiWithHiragana != nil {
                    kanjiOps.PrintMap("Kunyomiwithhiragana", kanjiOps.KunyomiWithHiragana[userInput], userInput, Readings)
                }

            } else if applicationSelector == "2" {
                utils.ClearScreen()
                fmt.Println("KEIGO ASSISTANT: Enter english word to get all keigo Readings ('exit' to quit)")

                if keigoOps.alreadyRead == false {
                    for key, _ := range keigoOps.keigoMap { 
                        englishpattern := regexp.MustCompile(`^[A-Za-z ]+\-$`)
                        romajipattern := regexp.MustCompile(`^[A-Za-z ]+$`)

                        if englishpattern.MatchString(key) {
                            keigoOps.keigoenglishslice = append(keigoOps.keigoenglishslice, key[:len(key)-1])
                        } else if romajipattern.MatchString(key) {
                            keigoOps.keigoromajislice = append(keigoOps.keigoromajislice, key)
                        } else {
                            keigoOps.keigojapaneseslice = append(keigoOps.keigojapaneseslice, key)
                        }
                    }
                }

                sort.Strings(keigoOps.keigoenglishslice)
                sort.Strings(keigoOps.keigoromajislice)

                keigoOps.alreadyRead = true

                fmt.Println(keigoOps.keigoenglishslice)
                fmt.Println(keigoOps.keigoromajislice)
                fmt.Println(keigoOps.keigojapaneseslice)

                scanner.Scan()
                userInput = scanner.Text()

                if keigoOps.keigoMap != nil {
                    keigoOps.printmapkeigo(userInput)
                }
            } else if applicationSelector == "3" {
                kanjiOps.FrequencyAnalysis("Onyomi")
                userInput = "exit"
            } else if applicationSelector == "4" {
                kanjiOps.FrequencyAnalysis("Kunyomi")
                userInput = "exit"
            } else if applicationSelector == "5" {
                kanjiOps.FrequencyAnalysis("Kunyomiwithhiragana")
                userInput = "exit"
            } 
            fmt.Println("Press Enter to continue...")
            
            scanner.Scan()
            _ = scanner.Text()
        }
    }
}

