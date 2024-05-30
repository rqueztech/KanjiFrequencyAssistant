package utils

import (
    "os"
    "os/exec"
    "runtime"
)

// Read the from the os/exec package and create custom clearscreen
func ClearScreen() {
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
