/*
The program computes total disk usage for each directory name specified in command line and stores it in a variable.
A -v flag renders a progress report for user while program runs. Otherwise, a spinner is shown.
If no directory name is specified in the command line, the default is the current directory.
This program derives from the du utility command. */

package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/stytch16/systemMgt-go/DiskUsage/filesys"
)

var verbose = flag.Bool("v", false, "show verbose progress reports")

func main() {
	/* Parse flag arguments and non-flag arguments */
	flag.Parse()
	roots := flag.Args()
	/* Specify default directory as current directory if none specified */
	if len(roots) == 0 {
		roots = []string{"."}
	}

	/* Start a goroutine to process the spinner if -v flag not set. */
	if !*verbose {
		go func() {
			go spinner(100 * time.Millisecond)
		}()
	}

	/* Get and print disk usage */
	nfiles, nbytes := filesystem.GetDiskUsage(roots, *verbose)
	printDiskUsageReport(nfiles, nbytes, roots)
}

/* spinner to display progress */
func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

/* printDiskUsage displays the disk usage of each directory specified in the command line in a organized format */
func printDiskUsageReport(nfiles, nbytes []int64, dirs []string) {
	fmt.Print("\r")
	for i, dir := range dirs {
		fmt.Printf("%-30s%10d files\t%.3f GB\n", dir, nfiles[i], float64(nbytes[i])/1e9)
	}
}
