// import the main package
package main

import (
    "bufio" // import bufio to scan in user input
    "fmt" // import the fmt package for printing
    "net/url"
    "os"
    "KanjiFrequencyHelper/csvoperations"
    "KanjiFrequencyHelper/utils"
    "regexp"
    "sort"
    "strings"
    "sync"
)

// Please look into possibly using channels -> still need to learn how to do this, maybe a mutex alternative
type KanjiReadings struct {
    onyomiMap map[string][]rune
    kunyomiMap map[string][]rune
    kunyomiWithHiragana map[string][]rune
    kanjiMeanings map[string][]rune
    readings map[string][]rune

    onyomifrequencyslice  [][]string
    kunyomifrequencyslice [][]string
    kunyomiwithhiraganafrequencyslice [][]string

    strings.Builder
    regex *regexp.Regexp
    mu sync.Mutex
}

func (kr *KanjiReadings) Lock() {
    kr.mu.Lock()
}

func (kr *KanjiReadings) Unlock() {
    kr.mu.Unlock()
}

type KeigoReadings struct {
    alreadyRead bool
    keigoMap map[string][]rune
    keigoenglishslice []string
    keigoromajislice []string
    keigojapaneseslice []string
    regex * regexp.Regexp
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


func (kanjiOps* KanjiReadings) loadfrequencies() {
    kanjiOps.onyomifrequencyslice = make([][]string, 80)
    kanjiOps.kunyomifrequencyslice = make([][]string, 30)
    kanjiOps.kunyomiwithhiraganafrequencyslice = make([][]string, 25)

    hiraganaPattern := regexp.MustCompile(`[A-Za-z]`)

    for key, value := range kanjiOps.onyomiMap {

        currentreadfrequency := len(value)

        if hiraganaPattern.MatchString(key) {
            kanjiOps.onyomifrequencyslice[currentreadfrequency] = append(kanjiOps.onyomifrequencyslice[currentreadfrequency], key)
        }

    }
    
    for key, value := range kanjiOps.kunyomiMap {

        currentreadfrequency := len(value)

        if hiraganaPattern.MatchString(key) {
            kanjiOps.kunyomifrequencyslice[currentreadfrequency] = append(kanjiOps.kunyomifrequencyslice[currentreadfrequency], key)
        }

    }
    for key, value := range kanjiOps.kunyomiWithHiragana {

        currentreadfrequency := len(value)

        if hiraganaPattern.MatchString(key) {
            kanjiOps.kunyomiwithhiraganafrequencyslice[currentreadfrequency] = append(kanjiOps.kunyomiwithhiraganafrequencyslice[currentreadfrequency], key)
        }
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
                kanjiOps.WriteString(jishoBaseLink)
                kanjiOps.WriteString(string(escaped))
                kanjiOps.WriteString("%20%23kanji")

                kanjilink := kanjiOps.String()
                linkOutput := strings.ReplaceAll(kanjilink, "\\n", "\n")
                kanjiOps.Reset()
                fmt.Printf("\n\n%s (%s): %s -> %s\nKanji Link: %s", kanjiString, userInput, meaningString, readingString, linkOutput)

                escaped = url.QueryEscape("*" + kanjiString + "*")
                kanjiOps.WriteString(jishoBaseLink)
                kanjiOps.WriteString(string(escaped))
                kanjiOps.WriteString("%20%23common")

                commonwordlink := kanjiOps.String()
                linkOutput = strings.ReplaceAll(commonwordlink, "\\n", "\n")
                kanjiOps.Reset()
                fmt.Println("Words Link: %s\n", commonwordlink)

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

func (kanjiOps* KanjiReadings) frequencyAnalysis(userinput string) {
    utils.ClearScreen()

    if userinput == "onyomi" {
        fmt.Println("Onyomi Frequency Report: ")
        for i := len(kanjiOps.onyomifrequencyslice) - 1; i >= 0; i-- {
            if kanjiOps.onyomifrequencyslice[i] != nil{
                fmt.Println(i, kanjiOps.onyomifrequencyslice[i])
            }
        }
    } else if userinput == "kunyomi" {
        fmt.Println("Kunyomi Frequency Report: ")
        for i := len(kanjiOps.kunyomifrequencyslice) - 1; i >= 0; i-- {
            if kanjiOps.kunyomifrequencyslice[i] != nil{
                fmt.Println(i, kanjiOps.kunyomifrequencyslice[i])
            }
        }
    } else if userinput == "kunyomiwithhiragana" {
        fmt.Println("Kunyomi with Hiragana Frequency Report: ")
        for i := len(kanjiOps.kunyomiwithhiraganafrequencyslice) - 1; i >= 0; i-- {
            if kanjiOps.kunyomiwithhiraganafrequencyslice[i] != nil{
                fmt.Println(i, kanjiOps.kunyomiwithhiraganafrequencyslice[i])
            }
        }
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
            csvMap, err := csvoperations.ReadCSV(filePath)
            
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

    kanjiOps.loadfrequencies()


    // now back in sequential mode, we can 

    // Loop to keep the program running unless the user types in "exit"
    for {
        utils.ClearScreen() 
        fmt.Print("Select Function:\n1. Kanji Finder\n2. Keigo Finder\n3. Onyomi\n4. Kunyomi\n5. KunyomiWithHiragana\n6. Exit\nEnter Input: ")

        scanner.Scan()
        applicationSelector := scanner.Text()
        
        userInput := ""

        var readings bool = true

        for userInput != "exit" {
            if applicationSelector == "1" {
                utils.ClearScreen()
                // create a bool to track readings
                

                fmt.Println("KANJI ASSISTANT: Enter (hiragana, romaji, or katakana to get readings")
                fmt.Println("Enter Input: ('exit' to quit, 'readings' toggles verbosity: ")
                
                if readings == true {
                    fmt.Println("Reading data enabled...")
                } else {
                    fmt.Println("Reading data silenced...")
                }

                scanner.Scan()
                userInput = scanner.Text()

                if userInput == "readings" {
                    readings = !readings
                    fmt.Println("Reading data silenced...")
                    _ = bufio.NewScanner(os.Stdin)
                    continue
                }


                // Send each string into the printMap
                if kanjiOps.onyomiMap != nil {
                    kanjiOps.printMap("onyomi", kanjiOps.onyomiMap[userInput], userInput, readings)
                }

                if kanjiOps.kunyomiMap != nil {
                    kanjiOps.printMap("kunyomi", kanjiOps.kunyomiMap[userInput], userInput, readings)
                }

                if kanjiOps.kunyomiWithHiragana != nil {
                    kanjiOps.printMap("kunyomiwithhiragana", kanjiOps.kunyomiWithHiragana[userInput], userInput, readings)
                }

            } else if applicationSelector == "2" {
                utils.ClearScreen()
                fmt.Println("KEIGO ASSISTANT: Enter english word to get all keigo readings ('exit' to quit)")

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
                kanjiOps.frequencyAnalysis("onyomi")
                userInput = "exit"
            } else if applicationSelector == "4" {
                kanjiOps.frequencyAnalysis("kunyomi")
                userInput = "exit"
            } else if applicationSelector == "5" {
                kanjiOps.frequencyAnalysis("kunyomiwithhiragana")
                userInput = "exit"
            } 
            fmt.Println("Press Enter to continue...")
            
            scanner.Scan()
            _ = scanner.Text()
        }
    }
}

