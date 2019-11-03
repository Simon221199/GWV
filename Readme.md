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

Simple run ```go run ./go PATH_TO_ENVIOMENT_TXT```

or

Use binary ```./bin/find_mac PATH_TO_ENVIOMENT_TXT```

or

Docker ```docker-compose run --rm gwv environment/test_env_3.txt```

### Output

Output for ```go run ./go ./environment/blatt3_environment.txt```

```
sourcing blatt3_environment.txt
######## Finding path form (4, 4) to (14, 7)
cell: (4, 4) --> -10
cell: (5, 4) --> -9
cell: (6, 5) --> -8
cell: (7, 5) --> -7
cell: (8, 4) --> -7
cell: (9, 5) --> -5
cell: (9, 3) --> -6
cell: (9, 4) --> -6
cell: (7, 3) --> -8
cell: (7, 4) --> -8
cell: (7, 2) --> -9
cell: (8, 1) --> -8
cell: (9, 1) --> -8
cell: (10, 1) --> -7
cell: (11, 2) --> -6
cell: (12, 2) --> -5
cell: (13, 2) --> -5
cell: (14, 2) --> -5
cell: (15, 3) --> -4
cell: (14, 4) --> -3
cell: (15, 4) --> -3
cell: (13, 4) --> -3
cell: (12, 5) --> -3
cell: (13, 6) --> -1
cell: (14, 7) --> 0
######## Path form (4, 4) to (14, 7)
Steps: 21
Coordinates:
(4, 4)
(5, 4)
(6, 5)
(7, 5)
(8, 4)
(7, 3)
(7, 2)
(8, 1)
(9, 1)
(10, 1)
(11, 2)
(12, 2)
(13, 2)
(14, 2)
(15, 3)
(14, 4)
(13, 4)
(12, 5)
(13, 6)
(14, 7)
xxxxxxxxxxxxxxxxxxxx
x                  x
x       xxx        x
x       x xxxxx    x
x   s     x        x
x       x x  xxxxxxx
x  xx xxxxx        x
x      x      g    x
x      x           x
xxxxxxxxxxxxxxxxxxxx
xxxxxxxxxxxxxxxxxxxx
x       ***        x
x      *xxx****    x
x      *x xxxxx*   x
x   **  * x  **    x
x     **x x *xxxxxxx
x  xx xxxxx  *     x
x      x      *    x
x      x           x
xxxxxxxxxxxxxxxxxxxx
```