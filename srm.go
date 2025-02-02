package main

import (
	"fmt"
	"os"
	"strings"
)

// Checklist
//     rm [-f | -i] [-dIRrvWx] file ...
//
// [X] -d      Attempt to remove directories as well as other types of files.
// [X] -f      Attempt to remove the files without prompting for confirmation, regardless of the file's permissions.  If the file does not exist, do not display a diagnostic message or modify the exit status
//             to reflect an error.  The -f option overrides any previous -i options.
// [X] -i      Request confirmation before attempting to remove each file, regardless of the file's permissions, or whether or not the standard input device is a terminal.  The -i option overrides any
//             previous -f options.
// [ ] -I      Request confirmation once if more than three files are being removed or if a directory is being recursively removed.  This is a far less intrusive option than -i yet provides almost the same
//             level of protection against mistakes.
// [X] -P      This flag has no effect.  It is kept only for backwards compatibility with 4.4BSD-Lite2.
// [X] -R      Attempt to remove the file hierarchy rooted in each file argument.  The -R option implies the -d option.  If the -i option is specified, the user is prompted for confirmation before each
//             directory's contents are processed (as well as before the attempt is made to remove the directory).  If the user does not respond affirmatively, the file hierarchy rooted in that directory is
//             skipped.
// [X] -r      Equivalent to -R.
// [ ] -v      Be verbose when deleting files, showing them as they are removed.
// [ ] -W      Attempt to undelete the named files.  Currently, this option can only be used to recover files covered by whiteouts in a union file system (see undelete(2)).
// [ ] -x      When removing a hierarchy, do not cross mount points.
// [ ] --      Makes all args after the double dash filenames (would be required to delete a file literally named "-i" for example)
// [ ] rename file if it already exists in destination

var VALIDARGS = []string{
	"-h",
	"--help",
	"-P",
	"-f",
	"-i",
	"-r",
	"-R",
	"-d",
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("    srm [-f | -i] [-dRr] <filepath> <...>")
	fmt.Println("Note:")
	fmt.Println("    Intended to replace `rm` via a shell alias")
}

func parseArgs() ([]string, []string) {
	// TODO support --
	// srm -- -f would remove a file named -f instead of being parsed as the "force flag"
	args := os.Args[1:]

	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	flags := []string{}
	files := []string{}
	seenDoubleDash := false

	for _, arg := range args {
		if arg == "--" {
			seenDoubleDash = true
			continue
		}

		// flags/params
		if In(arg, VALIDARGS) && !seenDoubleDash {
			flags = append(flags, arg)
			continue
		}

		// files
		files = append(files, arg)
	}

	return flags, files
}

// Get target dir for safely removed files
func getTargetRmDir() string {
	// First check if ~/.Trash exists (macOS)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Could not get users home dir")
		os.Exit(1)
	}

	path := homeDir + "/.Trash"
	if _, err := os.Stat(path); err == nil {
		// ~/.Trash
		return path
	}

	// Otherwise just use /tmp
	return "/tmp"
}

func main() {
	targetDir := getTargetRmDir()
	flags, files := parseArgs()

	// help
	helpFlag := In("-h", flags) || In("--help", flags)
	if helpFlag {
		usage()
		os.Exit(0)
	}

	// Force
	forceFlag := In("-f", flags)

	// interactive
	interactiveFlag := In("-i", flags)

	// recursive
	recursiveFlag := In("-r", flags) || In("-R", flags)

	// allow directories to be deleted
	directoryFlag := In("-d", flags)

	//fmt.Println("Flags: ", flags)
	//fmt.Println("Files: ", files)

	for _, filepath := range files {
		// directory and -r check
		if isDir, err := IsDir(filepath); isDir {
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if !recursiveFlag && !directoryFlag {
				fmt.Printf("srm: %s: is a directory\n", filepath)
				os.Exit(1)
			}
		}

		splitFilePath := strings.Split(filepath, "/")
		filename := splitFilePath[len(splitFilePath)-1]
		dest := targetDir + "/" + filename

		// -i
		if interactiveFlag {
			var interactiveResponse string
			fmt.Printf("remove %s?\n", filepath)
			fmt.Scanln(&interactiveResponse)
			interactiveResponse = strings.ToLower(interactiveResponse)
			if !In(interactiveResponse, []string{"y", "yes", "yea", "yeah", "da", "si", "letsgo"}) {
				continue
			}
		}

		// fmt.Printf("attempting to move %s to %s\n", filepath, dest)

		// if it ends with a / strip it
		if filepath[len(filepath)-1:] == "/" {
			filepath = strings.TrimRight(filepath, "/")
		}

		// check file isn't RO
		fileIsReadOnly, err := IsReadOnly(filepath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if fileIsReadOnly && !forceFlag {
			fmt.Println("File is read-only")
			os.Exit(1)
		}

		os.Rename(filepath, dest)
	}
}
