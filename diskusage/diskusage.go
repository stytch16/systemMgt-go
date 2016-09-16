/* Package du reports disk usage of one or more directories specified on command line */
package filesystem

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

/* GetDiskUsage computes the total disk usage along with number of files in each directory name in 'roots'.
If 'verbose' is true, GetDiskUsage prints a verbose progress report every 100 ms while it runs. */
func GetDiskUsage(roots []string, verbose bool) (nfiles, nbytes []int64) {
	var tick <-chan time.Time
	if verbose {
		tick = time.Tick(100 * time.Millisecond)
	}

	/* Output variables */
	nfiles = make([]int64, len(roots)) /* Number of files */
	nbytes = make([]int64, len(roots)) /* Size of directory (bytes) */

	/* Initiate a channel to transmit file sizes for each specified directory */
	var fileSizes []chan int64
	for range roots {
		fileSizes = append(fileSizes, make(chan int64))
	}

	/* For each directory, start a goroutine to walk thru it and compute its size. */
	for i, root := range roots {
		var wg sync.WaitGroup
		wg.Add(1) /* Add 1 for every goroutine that calls WalkDir */
		go walkDir(root, &wg, fileSizes[i])
		go func() {
			wg.Wait()
			close(fileSizes[i])
		}()
	loop:
		/* Receive data from fileSizes channel. When it has been closed, break from loop and continue to next directory. */
		for {
			select {
			case size, ok := <-fileSizes[i]:
				if !ok {
					break loop
				}
				/* record data */
				nfiles[i]++
				nbytes[i] += size
			case <-tick:
				/* print progress report */
				printProgress(nfiles[i], nbytes[i], root)
			}
		}
	}
	return nfiles, nbytes
}

/* print statement for -v flag */
func printProgress(nfiles, nbytes int64, dir string) {
	fmt.Printf("%-20s%10d files\t%.3f GB\n", dir, nfiles, float64(nbytes)/1e9)
}

/* walkDir recursively traverses the file tree rooted at dir and sends size of each found file on filesizes channel */
func walkDir(dir string, wg *sync.WaitGroup, filesizes chan<- int64) {
	defer wg.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			wg.Add(1)
			subdir := filepath.Join(dir, entry.Name()) /* path/to/directory/entry */
			go walkDir(subdir, wg, filesizes)
		} else {
			filesizes <- entry.Size()
		}
	}
}

/* Counting semaphore: limit number of open files to 20 */
var sema = make(chan struct{}, 100)

/* dirents returns entries of directory dir */
func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}
	defer func() {
		<-sema
	}()

	entries, err := ioutil.ReadDir(dir) /* returns list of directory entries of dir sorted by filename */
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}
