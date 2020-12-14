# Overview

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/d32d8df3c31d43969c6aa8876af40e22)](https://app.codacy.com/gh/josecordaz/wndyr?utm_source=github.com&utm_medium=referral&utm_content=josecordaz/wndyr&utm_campaign=Badge_Grade)

Download images gathered by NASA's Curiosity, Opportunity, and Spirit rovers on Mars for specific date.

## Build executable binary
```shell
go build -o nasa && chmod +x nasa
```
### Build Windows executable
```shell
go build -o nasa.exe
```

## Running binary
```shell
./nasa YYYY-M-D
```

## Running on Windows
```shell
nasa.exe YYYY-M-D
```

## Example
### Download images from 2020-Aug-6
```shell
./nasa 2020-8-6
```
### Download images from 2020-Aug-6 on windows
```shell
nasa.exe 2020-8-6
```
