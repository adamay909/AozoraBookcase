/*
aozoraBookcase serves books from Aozora Bunko. The options are

	-c	force re-downloading of all data from Aozora Bunko

	-children
	  	start a kid's library

	-d string
	  	directory containing server files. Defaults to $HOME/aozorabunko.

	-i string
	  	network interface

	-p string
	  	network interface (default "3333")

	-refresh string
	  	interval between library refreshes. (default "24h")

	-strict
	  	set library to show only public domain texts (default true)

	-v	show log on screen
*/
package main
