# systemMgt-go
Go/Golang library for file and directory management primarily for Linux systems. More updates and additions to follow.

 1. **DiskUsage**: Get the size and number of files of each directory specified in the command line.
 
 *Extract DiskUsage/diskusage and run it*

```
$ ./diskusage [-v]  /usr [/home] [...]
```
 
 2. **ImageConv**: Convert format of every image file specified in the command line. Supported formats are JPEG, PNG, and GIF.
 
 *Extract and run ImageConv/imgconv*

```
$ ./imgconv -type jpeg [-dest] [~/Pictures] pic1.png [pic2.png] [...]
```
