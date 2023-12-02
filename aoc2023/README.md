# AoC 2023

To run a day, you can use the program in this folder. 

Run day 1, example 1

`go run . -d 1 -e1`

OR Run day 2 over the full input

`go run . -d 2`

Add debug logging

`go run . -d 2 --debug`

To run in a specific sub folder, you have to deal with piping input
files and setting env vars for log levels. So, I just use the launcher.
