# GWV

### Group XYZ

| Name | Matrikel nummer |
| :---: | :-------------: |
|Patrick Zierahn | 7065799 |
|Simon ... | ... |
|Pascal ... | ... |

### Sources

* [A*_search_algorithm](https://en.wikipedia.org/wiki/A*_search_algorithm)
* [a-search-algorithm](https://www.geeksforgeeks.org/a-search-algorithm/)

### How to use

Algorithms: ```best-first```, ```aStar```, ```breadth-first```, ```depth-first```

Simple run ```go run ./go SEARCH_Algorithms PATH_TO_ENVIOMENT_TXT```

or

Use binary ```./bin/find_mac SEARCH_ALGORITHMS PATH_TO_ENVIOMENT_TXT```

or

Docker ```docker-compose run --rm gwv a* environment/blatt3_environment_portal.txt```

### How to compile

```
env GOOS=darwin GOARCH=amd64 go build -o bin/search_darwin ./go
env GOOS=linux GOARCH=amd64 go build -o bin/search_linux ./go
env GOOS=windows GOARCH=amd64 go build -o bin/search_windows.exe ./go
```
