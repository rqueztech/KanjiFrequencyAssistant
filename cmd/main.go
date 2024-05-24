// import the main package
package main

import (
    "bufio" // import bufio to scan in user input
    "fmt" // import the fmt package for printing
    "net/url"
    "os"
    "os/exec"
    "regexp"
    "runtime"
    "strings"
    "sync"
)

// Read the from the os/exec package and create custom clearscreen
func clearScreen() {
    // Create a pointer to the exec.Cmd struct, which is used to build the command to be used.
    var cmd *exec.Cmd

    // Check the runtime of the OS and create the command accordingly
    if runtime.GOOS == "windows" {
        cmd = exec.Command("cmd", "/c", "cls")
    } else {
        cmd = exec.Command("clear")
    }

    // Set the Stdout to the os.Stdout
    cmd.Stdout = os.Stdout

    // Run the command that was created
    cmd.Run()
}

// Please look into possibly using channels -> still need to learn how to do this, maybe a mutex alternative
type KanjiReadings struct {
    sync.Mutex
    onyomiMap, kunyomiMap, kunyomiWithHiragana, kanjiMeanings, readings map[string][]rune
    strings.Builder
    regex *regexp.Regexp
}

type KeigoReadings struct {
    keigoMap map[string][]rune
}

func (keigoOps* KeigoReadings) printmapkeigo(userinput string) {
    if keigoOps.keigoMap[userinput] != nil {
        keigostring := string(keigoOps.keigoMap[userinput])
        keigostring = strings.ReplaceAll(keigostring, "*", "\n")
        fmt.Printf("\n%s: %s\n", keigostring, userinput)
    }
}

func (kanjiOps* KanjiReadings) printMap(title string, map_result []rune, userInput string, readings bool) {
    // Print out the name of the function

    // Jisho link string baseline
    jishoBaseLink := "https://www.jisho.org/search/"
    
    if map_result != nil {
        fmt.Printf("============ %s ================", title)

        // Check if the value exists in the map, if not prints out DOES NOT EXIST
        for _, currentKanjiRune := range(map_result) {
            kanjiString := string(currentKanjiRune)
            escaped:= url.QueryEscape(kanjiString)

            currentKanji := string(currentKanjiRune)
            readingString := string(kanjiOps.readings[currentKanji])
            readingString = strings.ReplaceAll(readingString, "\\n", "\n")
            meaningString := string(kanjiOps.kanjiMeanings[currentKanji])

            if readings == true {
                kanjiOps.WriteString("\nLink: ")
                kanjiOps.WriteString(jishoBaseLink)
                kanjiOps.WriteString(string(escaped))
                kanjiOps.WriteString("%20%23kanji")

                kanjilink := kanjiOps.String()
                linkOutput := strings.ReplaceAll(kanjilink, "\\n", "\n")
                kanjiOps.Reset()
                    fmt.Printf("\n%s (%s): %s -> %s %s\n", kanjiString, userInput, meaningString, readingString, linkOutput)
                } else {
                    kanjiOps.WriteString(jishoBaseLink)
                    kanjiOps.WriteString(string(escaped))
                    kanjiOps.WriteString("%20%23kanji")

                    linkOutput := kanjiOps.String()
                    fmt.Printf("\n%s (%s): %s -> %s", kanjiString, userInput, meaningString, linkOutput)
                    kanjiOps.Reset()
                }
        } 

        fmt.Printf("\nNumber of [[%s]] readings --> : %d\n", userInput, len(map_result))
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
    kanjiOps := &KanjiReadings{}
    keigoOps := &KeigoReadings{}


    // Create a scanner used to read user input/options
    scanner := bufio.NewScanner(os.Stdin)

    // create a bool to track readings
    var readings bool = true

    // Create a waitgroup
    var wg sync.WaitGroup

    // Create  string array of filepaths, this will help us save into proper map
    filePaths := []string {
        "./resources/KanjiFrequencyListOnyomi.csv",
        "./resources/KanjiFrequencyListKunyomi.csv",
        "./resources/KunyomiWithHiragana.csv",
        "./resources/KanjiMeanings.csv",
        "./resources/all_readings_string.csv",
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
            csvMap, err := ReadCSV(filePath)
            
            if err != nil {
                handleError(err, "Error with file path: " + filePath)
                return
            }

            switch filePath {
                case "./resources/KanjiFrequencyListOnyomi.csv":
                    kanjiOps.onyomiMap = csvMap

                case "./resources/KanjiFrequencyListKunyomi.csv": 
                    kanjiOps.kunyomiMap = csvMap

                case "./resources/KunyomiWithHiragana.csv":
                    kanjiOps.kunyomiWithHiragana = csvMap

                case "./resources/KanjiMeanings.csv":
                    kanjiOps.kanjiMeanings = csvMap

                case "./resources/all_readings_string.csv":
                    // setting the csvmap directly into the KanjiReadings map
                    kanjiOps.readings = csvMap

                case "./resources/keigo_mapper.csv":
                    keigoOps.keigoMap = csvMap
            }

        }(filePath)
    }

    // Wait for all the go routines to finish, wait on all five files
    wg.Wait()

    // now back in sequential mode, we can 

    // Loop to keep the program running unless the user types in "exit"
    for {
        clearScreen() 
        fmt.Print("Select Function:\n1. Kanji Finder\n2. Keigo Finder\n3. Exit\nEnter Input: ")

        scanner.Scan()
        applicationSelector := scanner.Text()
        
        userInput := ""

        for userInput != "exit" {
            if applicationSelector == "1" {
                clearScreen()
                fmt.Println("KANJI ASSISTANT: Enter (hiragana, romaji, or katakana to get readings")
                fmt.Println("Enter Input: ('exit' to quit, 'readings' toggles verbosity: ")
            } else if applicationSelector == "2" {
                clearScreen()
                fmt.Println("KEIGO ASSISTANT: Enter english word to get all keigo readings ('exit' to quit)")

                for key, _ := range keigoOps.keigoMap {
                    fmt.Printf("%s, ", key)
                }
                    
            }

            scanner.Scan()
            userInput = scanner.Text()

            if readings == true {
                fmt.Println("Reading data enabled...")
            } else {
                fmt.Println("Reading data silenced...")
            }

            
            if userInput == "exit" {
                fmt.Println("Exiting the program...")
                break
            } else if userInput == "readings" {
                readings = !readings
                fmt.Println("Reading data silenced...")
                _ = bufio.NewScanner(os.Stdin)
                continue
            }

            if applicationSelector == "1" {
                // Send each string into the printMap
                if kanjiOps.onyomiMap != nil {
                    kanjiOps.printMap("onyomi", kanjiOps.onyomiMap[userInput], userInput, readings)
                }

                if kanjiOps.kunyomiMap != nil {
                    kanjiOps.printMap("kunyomi", kanjiOps.kunyomiMap[userInput], userInput, readings)
                }

                if kanjiOps.kunyomiWithHiragana != nil {
                    kanjiOps.printMap("kunyomi with hiragana", kanjiOps.kunyomiWithHiragana[userInput], userInput, readings)
                }
            } else if applicationSelector == "2" {
                if keigoOps.keigoMap != nil {
                    keigoOps.printmapkeigo(userInput)
                }
            } 
            fmt.Println("Press Enter to continue...")
            fmt.Scanln()
        }
    }
}
