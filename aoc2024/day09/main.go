package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/adventofcode/pkg/straid"
)

func main() {
	log := logging.DefaultLogger()
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	line := scanner.Text()

	log.Debugw("line", "line", line)

	diskBuilder := make([]string, 0)
	fileID := 0
	space := false
	for _, c := range line {
		i, err := strconv.Atoi(string(c))
		if err != nil {
			panic(err)
		}
		if space {
			for j := 0; j < i; j++ {
				diskBuilder = append(diskBuilder, ".")
			}
			space = false
		} else {
			for j := 0; j < i; j++ {
				diskBuilder = append(diskBuilder, fmt.Sprintf("%d", fileID))
			}
			space = true
			fileID++
		}
	}
	log.Debugw("disk", "disk", strings.Join(diskBuilder, ""))

	disk2 := make([]string, len(diskBuilder))
	copy(disk2, diskBuilder)

	{
		front := 0
		back := len(diskBuilder) - 1
		for front < back {
			for diskBuilder[front] != "." {
				front++
			}
			for diskBuilder[back] == "." {
				back--
			}
			if front >= back {
				break
			}

			log.Debugw("moving", "front", front, "val", diskBuilder[front], "back", back, "backVal", diskBuilder[back])

			diskBuilder[front] = diskBuilder[back]
			front++
			diskBuilder[back] = "."
			back--
		}
		var checksum int64
		for i, val := range diskBuilder {
			if val != "." {
				checksum += straid.AsInt(val) * int64(i)
			}
		}

		log.Debugw("defraged", "disk", strings.Join(diskBuilder, ""))
		log.Infow("checksum", "part1", checksum)
	}

	{
		for toMove := fileID - 1; toMove > 0; toMove-- {
			fileIDStr := fmt.Sprintf("%d", toMove)
			index := 0
			length := 0
			for i, val := range disk2 {
				if index == 0 && val == fileIDStr {
					index = i
					length = 1
					continue
				} else if index > 0 && val == fileIDStr {
					length++
				} else if index > 0 && val != fileIDStr {
					break
				}
			}
			log.Debugw("defragging", "toMove", toMove, "index", index, "length", length)

			// find an empty space to put this file in
			moveTo := 0
			moveBlock := false
			for ; moveTo < len(disk2)-length; moveTo++ {
				if disk2[moveTo] == "." {
					found := true
					for i := 0; i < length; i++ {
						if disk2[moveTo+i] == "." {
							continue
						} else {
							found = false
							break
						}
					}
					if found {
						moveBlock = true
						break
					}
				}
			}
			if moveBlock {
				if moveTo < index {
					log.Debugw("moving", "toMove", toMove, "from", index, "to", moveTo)
					for i := 0; i < length; i++ {
						disk2[moveTo+i] = fileIDStr
						disk2[index+i] = "."
					}
				}
			}

		}
		var checksum int64
		for i, val := range disk2 {
			if val != "." {
				checksum += straid.AsInt(val) * int64(i)
			}
		}

		log.Debugw("defraged", "disk", strings.Join(disk2, ""))
		log.Infow("checksum", "part2", checksum)
	}
}
