# About ui2go #

ui2go is a toolkit for the creation of graphical user interfaces. It is a prototype that features some interesting concepts, especially in the areas of event dispatching and widget layout management.

![http://www.projectory.de/ui2go/paint-small.png](http://www.projectory.de/ui2go/paint-small.png)

# Documentation #

  * [White Paper](http://www.projectory.de/ui2go/whitepaper/index.html): General ideas and concepts.
  * [Design Decisions](http://www.projectory.de/ui2go/designdecisions/designdecisions.html): Thoughts about the design in short notation.

# Installation #

For high performance and to keep the code small ui2go uses C libraries and relies upon
[go-cairo](https://github.com/ungerik/go-cairo) for drawing operations.

## Debian/Ubuntu based distributions ##

  * go get code.google.com/p/ui2go
  * sudo apt-get install libxcb1-dev
  * sudo apt-get install libcairo2-dev

## Other distributions ##

I don't know the corresponding package names in other distributions. So I am afraid you have to figure out yourself. [go-cairo](https://github.com/ungerik/go-cairo) provides additional information about the go-cairo installation.

## Windows and Mac ##

Not implemented yet.
