# Gotron

A small library that allows Golang programs to draw to the screen by sending messages to an Electron (NodeJS) frontend.

![Gotron Screenshot](/examples/screenshot.gif?raw=true | width=624)

Usage:

* Write a Go app that uses the Gotron API; for an example see `basic.go`, `electra.go`, or `swarmz.go`
* Compile your Go app
* Edit config.json to point to the executable
* Run Electron in the directory, e.g. `electron .`
