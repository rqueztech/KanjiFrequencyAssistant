package keigo

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "regexp"
)

type KeigoReadings struct {
    AlreadyRead bool
    KeigoMap map[string][]rune
    KeigoEnglishSlice []string
    KeigoRomajiSlice []string
    KeigoJapaneseSlice []string
    regex * regexp.Regexp
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
