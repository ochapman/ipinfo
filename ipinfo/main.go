/*
* File Name:	main.go
* Description:
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2014-10-12
 */
package main

import (
	"fmt"
	"github.com/ochapman/ipinfo"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Usage %s IP\n", os.Args[0])
		return
	}
	ip := os.Args[1]
	res, err := ipinfo.Get(ip)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(res)
}
