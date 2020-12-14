# Overview

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
nasa 2020-8-6
```
