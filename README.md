my-ls

my-ls is a lightweight file listing utility written in Go. Inspired by the Unix ls command, it provides a rich set of features to display directory contents in various formats with options for recursion, sorting, filtering, and even colorized output.
Features

    Long Listing Format (-l):
    Displays detailed file information (permissions, number of links, owner, group, size, and modification time).
    (See [display.go] and [utils.go].)

    Recursive Listing (-R):
    Recursively traverses directories to list subdirectories and their contents.
    (Implemented in [recursive.go].)

    Show Hidden Files (-a):
    Includes files and directories starting with a dot (.).
    (Implemented in [filter.go].)

    Time Sorting (-t):
    Sorts files by modification time, showing the most recently modified files     first.
    (See [sort.go].)

    Reverse Sorting (-r):
    Reverses the sort order.
    (See [sort.go].)

    Output Capture (-c):
    Optionally capture the output into a file (myls_output.txt) as well as display it on the console.
    (Handled in [output.go].)

    Colorized Output:
    Applies ANSI colors to differentiate file types such as directories, executables, symlinks, devices, and sockets.
    (See [colorize.go].)

Installation

    git clone https://github.com/Br0m0r/My-ls.git
    go build -o myls 

Usage

    ./myls [options] [path...]

Options

    -l: Use long listing format.
    -R: List subdirectories recursively.
    -a: Include directory entries whose names begin with a dot (.).
    -t: Sort by modification time, newest first.
    -r: Reverse order while sorting.
    -c: Capture output to a file (myls_output.txt).
    -h: Display help and exit.

Examples:

List files in long format for the current directory:

    ./myls -l .

Recursively list all files (including hidden ones) in a directory:

    ./myls -Ra /path/to/directory

Project Structure

    main.go
    Entry point that delegates execution to the core listing logic.
    (See [main.go].)

    ls.go
    Contains the main logic for processing paths, handling errors, and coordinating the listing process.
    (See [ls.go].)

    recursive.go
    Implements recursive directory traversal for the -R flag.
    (See [recursive.go].)

    flags.go
    Parses command-line arguments and sets options accordingly.
    (See [flags.go].)

    display.go
    Handles the output formatting for both standard and long listing formats.
    (See [display.go].)

    colorize.go
    Provides functions for adding ANSI color codes to file names based on type.
    (See [colorize.go].)

    sort.go
    Implements file sorting based on alphabetical order, modification time, and reverse order.
    (See [sort.go].)

    filter.go
    Filters files to include or exclude hidden files, and adds pseudo entries for . and .. when needed.
    (See [filter.go].)

    utils.go
    Contains utility functions for fetching file permissions, owner, and group information.
    (See [utils.go].)

    logger.go
    Sets up logging to help with error tracking and debugging.
    (See [logger.go].)

    output.go
    Manages output streams, allowing output to be directed to both the console and a file.
    (See [output.go].)

