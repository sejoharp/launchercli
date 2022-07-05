# launchercli

## About The Project
command line launcher for reoccurring tasks.

## Getting Started
download the binary for your architecture and put it in your path.

## Usage
### create config file
1) create `~/.launcher-config.json`
2) put in the following format:
```
[
    {
        "list": "for directory in $(ls -d /home/user/repos/*); do echo $(basename $directory),$directory;done",
        "command": "code"
    }
]
```
#### list
It should return a line for every directory and display name you want to use. The format is `[display name],[directory]`.  
e.g. `myrepo,/home/user/repos/myrepo`

#### command
It should contain the command you want to execute with the directory from the `list` command.  
e.g. `code /home/user/repos/myrepo`

### execution
Call the binary. It will show a list with all entries. Start typing to narrow down the option. Press enter to execute the command. 
## Build from source
### Prerequisites
1. Install `go`  
    a) macos `brew install go`
    
### Installation
see `Makefile` or do `make` to see all commands, e.g. `make install-linux-amd64`
