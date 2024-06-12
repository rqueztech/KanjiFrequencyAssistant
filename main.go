// import the main package
package main

// Import all necessary packages 
import (
    "bufio" // import bufio to scan in user input
    "fmt" // import the fmt package for printing
    "os"
    "net/url"
    "KanjiFrequencyHelper/csvoperations"
    "KanjiFrequencyHelper/utils"
    "KanjiFrequencyHelper/kanji"
    "KanjiFrequencyHelper/keigo"
    "strings"
    "sync"
)

// Builder struct to be used by the program
type Builder struct {
    strings.Builder
}

// Main function
func main() {

    // create kanji ops blank pointer
    kanjiOps := &kanji.KanjiReadings{}
    keigoOps := &keigo.KeigoReadings{
        AlreadyRead: false,
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
        "./resources/FullDetailsBoth.csv",
        "./resources/FullDetailsKunyomi.csv",
        "./resources/FullDetailsOnyomi.csv",
        "./resources/KunyomiByEndings.csv",
        "./resources/TranslatorMap.csv",
        "./resources/KunyomiTransatives.csv",
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
                fmt.Println("File Not Found")
                return
            }

            // Create switch statement for importing all csv files
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
                    keigoOps.KeigoMap = csvMap

                case "./resources/FullDetailsBoth.csv":
					kanjiOps.FullDetailsBoth = csvMap

                case "./resources/FullDetailsKunyomi.csv":
					kanjiOps.FullDetailsBoth = csvMap

                case "./resources/FullDetailsOnyomi.csv":
					kanjiOps.FullDetailsBoth = csvMap

                case "./resources/KunyomiByEndings.csv":
                    kanjiOps.KunyomiByEndings = csvMap

                case "./resources/TranslatorMap.csv":
                    kanjiOps.TranslatorMap = csvMap

                case "./resources/KunyomiTransatives.csv":
                    kanjiOps.KunyomiTransatives = csvMap
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
        fmt.Print("Select Function:\n1. Kanji Finder\n2. Keigo Finder\n3. Onyomi\n4. Kunyomi\n5. KunyomiWithHiragana\n6. Kanji Only\n7. Enter Phrase to link \n8. KunyomiByEndings\nEnter Input: ")

        scanner.Scan()
        applicationSelector := scanner.Text()
        
        userInput := ""

        
        for userInput != "exit" {
            if applicationSelector == "1" {
                utils.ClearScreen()
                // create a bool to track Readings
                

                fmt.Println("KANJI ASSISTANT: Enter (hiragana, romaji, or katakana to get Readings")
                fmt.Println("Enter Input: ('exit' to quit, 'Readings' toggles verbosity: ")
                
                if kanjiOps.ShowReadings == true {
                    fmt.Println("Reading data enabled...")
                } else {
                    fmt.Println("Reading data silenced...")
                }

                scanner.Scan()
                userInput = scanner.Text()

                if userInput == "readings" {
                    kanjiOps.ShowReadings = !kanjiOps.ShowReadings
                    _ = bufio.NewScanner(os.Stdin)
                    continue
                }
                    
                // Send each string into the PrintMap
                if kanjiOps.OnyomiMap != nil && kanjiOps.OnyomiMap[userInput] != nil{
                    kanjiOps.PrintMap("Onyomi", kanjiOps.OnyomiMap[userInput], userInput)
                }

                if kanjiOps.KunyomiMap != nil && kanjiOps.KunyomiMap[userInput] != nil{
                    kanjiOps.PrintMap("Kunyomi", kanjiOps.KunyomiMap[userInput], userInput)
                }

                if kanjiOps.KunyomiWithHiragana != nil && kanjiOps.KunyomiWithHiragana[userInput] != nil{
                    kanjiOps.PrintMap("Kunyomiwithhiragana", kanjiOps.KunyomiWithHiragana[userInput], userInput)
                }

            } else if applicationSelector == "2" {
                keigoOps.ReadKeigo(userInput)
                scanner.Scan()
                userInput = scanner.Text()

                if keigoOps.KeigoMap != nil {
                    keigoOps.PrintMapKeigo(userInput)
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
            } else if applicationSelector == "6" {
                fmt.Println("Enter Kanji Here: ")
                scanner.Scan()
                userInput = scanner.Text()

                charMap := make(map[rune]bool)

                userInput = utils.GetPatternCleaning().RemoveNonKanji(userInput)

                for _, char := range userInput {
                    charMap[char] = true
                }

                var removeduplicates string
                for char := range charMap {
                    removeduplicates += string(char)
                }

                for _, char := range removeduplicates {
                    currentKanji := string(char)
                    QueryEscaped := url.QueryEscape(currentKanji)
                    kanjiDefinition := string(kanjiOps.KanjiMeanings[string(char)])

                    numReadings := string(kanjiOps.Readings[string(char)])
                    slicedreadings := strings.Split(numReadings, "*")
                    
                    fmt.Println(slicedreadings[1])

                    if slicedreadings[1] == "Both" {
                        fmt.Println(slicedreadings[2])
                        fmt.Println(slicedreadings[4])
                    } else if slicedreadings[1] == "Kunyomi" {
                        fmt.Println("NUMBER OF READINGS")
                        fmt.Println(slicedreadings[4])
                    } else if slicedreadings[1] == "Onyomi" {
                        fmt.Println("NUMBER OF READIGNS")
                        fmt.Println(slicedreadings[4])
                    }

                    fmt.Printf("%s -> \t%s : \nLink: https://www.jisho.org/search/%s%20%23kanji", currentKanji, kanjiDefinition, QueryEscaped)
                }
            } else if applicationSelector == "7" {
                fmt.Println("Enter Kanji Here: ")
                userInput = ""

                for userInput != "exit" {
                    fmt.Println("Enter Input: ('exit' to quit)")
                    scanner.Scan()
                    userInput = scanner.Text()
                    _ = scanner.Text()

                    if utils.GetPatternCleaning().IsRomajiPattern(userInput) {
                        fmt.Println("Please Enter Japense Characters")
                    } else {
                        QueryEscaped := url.QueryEscape(userInput)

                        fmt.Printf("%s -> Link: https://www.jisho.org/search/%s\n",  userInput, QueryEscaped)
                    }
                }


                userInput = "exit"
            } else if applicationSelector == "8" {
                utils.ClearScreen()

                fmt.Println("Enter Kunyomi Word Ending Here: ")
                scanner.Scan()
                userInput = scanner.Text()

                hiraganatranslation := kanjiOps.TranslatorMap[userInput]


                if userInput == "clear" {
                    utils.ClearScreen()
                    continue
                }

                if userInput == "exit" {
                    userInput = "exit"
                    break
                }
                
                transativecount := 0
                intransativecount := 0
                endingscount := 0
                iAdjCount := 0
                naAdjCount := 0
                AdverbCount := 0
                ConjunctionCount := 0

                for _, currentkanji := range(kanjiOps.KunyomiByEndings[userInput]) {
                    jointword := url.QueryEscape(string(currentkanji) + string(hiraganatranslation))
                    jointstring := string(currentkanji) + string(hiraganatranslation)
                    transatives := kanjiOps.KunyomiTransatives[jointstring]
                    
                    parts := strings.Split(string(transatives), "*")
                    
                    if len(parts) > 2 {
                        wordtype := parts[0]
                        definition := parts[1]
                        hiraganized := parts[2]
                        

                        fmt.Printf("%s %s (%s) -> %s |%s| https://www.jisho.org/search/%s\n", string(currentkanji), string(hiraganatranslation), hiraganized, wordtype, definition, jointword)

                        if wordtype == "Transative" {
                            transativecount++
                        } else if wordtype == "Intransative" {
                            intransativecount++
                        } else if wordtype == "Ending" {
                            endingscount++
                        } else if wordtype == "iAdj" {
                            iAdjCount++
                        } else if wordtype == "naAdj" {
                            naAdjCount++
                        } else if wordtype == "Adverb" {
                            AdverbCount++
                        } else if wordtype == "Conjunction" {
                            ConjunctionCount++
                        }
                    }
                }

                fmt.Printf("\nNumber of occurences [%s]: %s -> %s\n", userInput, endingscount)
                
                if transativecount > 0 {
                    fmt.Println("Total Transatives: ", transativecount)
                } 
				if intransativecount > 0 {
                    fmt.Println("Total Intransatives: ", intransativecount)
                } 
				if iAdjCount > 0 {
                    fmt.Println("Total iAdj: ", iAdjCount)
                } 
				if naAdjCount > 0 {
                    fmt.Println("Total naAdj: ", naAdjCount)
                } 
				if AdverbCount > 0 {
                    fmt.Println("Total Adverbs: ", AdverbCount)
                } 
				if ConjunctionCount > 0 {
                    fmt.Println("Total Conjunctions: ", ConjunctionCount)
                }

            } else {
                utils.ClearScreen()
                fmt.Println("Enter valid input")
                userInput = "exit"
            }

            fmt.Println("Press Enter to continue...")
            
            scanner.Scan()
            _ = scanner.Text()
        }
    }
}

