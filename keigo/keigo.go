package keigo

import (
    "KanjiFrequencyHelper/utils"
    "bufio"
    "fmt"
    "os"
    "strings"
    "sort"
)

type KeigoReadings struct {
    AlreadyRead bool
    KeigoMap map[string][]rune
    KeigoEnglishSlice []string
    KeigoRomajiSlice []string
    KeigoJapaneseSlice []string

}

func (keigoOps* KeigoReadings) ReadKeigo(userinput string) {
    utils.ClearScreen()
    fmt.Println("KEIGO ASSISTANT: Enter english word to get all keigo Readings ('exit' to quit)")

    if keigoOps.AlreadyRead == false {
        for key, _ := range keigoOps.KeigoMap {
            if utils.GetPatternCleaning().IsIllegalCharactersPattern(key) {
                break
            }

            if utils.GetPatternCleaning().IsEnglishPattern(key) {
                keigoOps.KeigoEnglishSlice = append(keigoOps.KeigoEnglishSlice, key[:len(key)-1])
            } else if utils.GetPatternCleaning().IsRomajiPattern(key) {
                keigoOps.KeigoRomajiSlice = append(keigoOps.KeigoRomajiSlice, key)
            } else {
                keigoOps.KeigoJapaneseSlice = append(keigoOps.KeigoJapaneseSlice, key)
            }
        }
    }

    sort.Strings(keigoOps.KeigoEnglishSlice)
    sort.Strings(keigoOps.KeigoRomajiSlice)

    keigoOps.AlreadyRead = true

    fmt.Println(keigoOps.KeigoEnglishSlice)
    fmt.Println(keigoOps.KeigoRomajiSlice)
    fmt.Println(keigoOps.KeigoJapaneseSlice)
}

func (keigoOps* KeigoReadings) PrintMapKeigo(userinput string) {
    fmt.Println(userinput)
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Println(scanner.Text())
    userintputdashed := userinput + "-"

    if keigoOps.KeigoMap[userinput] != nil {
        keigostring := string(keigoOps.KeigoMap[userinput])
        keigostring = strings.ReplaceAll(keigostring, "*", "\n")
        fmt.Printf("\n%s: %s\n", keigostring, userinput)
    } else if keigoOps.KeigoMap[userintputdashed] != nil {
        keigostring := string(keigoOps.KeigoMap[userintputdashed])
        keigostring = strings.ReplaceAll(keigostring, "*", "\n")
        fmt.Printf("\n%s: %s\n", keigostring, userinput)
    } else {
        fmt.Printf("\n%s: DOES NOT EXIST\n", userinput)
    }
}
