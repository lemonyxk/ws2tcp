/**
* @program: ws2tcp
*
* @create: 2024-07-10 13:24
**/

package main

import (
	"fmt"
	"os"
)

func GetArgs(flag ...string) string {
	var args = os.Args[1:]
	for i := 0; i < len(args); i++ {
		for j := 0; j < len(flag); j++ {
			if args[i] == flag[j] {
				if i+1 < len(args) {
					return args[i+1]
				}
			}
		}
	}
	return ""
}

func HasArgs(flag ...string) bool {
	var args = os.Args[1:]
	for i := 0; i < len(args); i++ {
		for j := 0; j < len(flag); j++ {
			if args[i] == flag[j] {
				return true
			}
		}
	}
	return false
}

func Exit(msg any) {
	fmt.Printf("%v\n", msg)
	os.Exit(0)
}
