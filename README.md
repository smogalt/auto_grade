# auto_grade

automatically grade python scripts. 

### installation
all packages are in the standard library. install with ```go install github.com/smogalt/auto_grade@latest``` then add ```~/go/bin``` to your path if it isn't already there. or manual install with ```git clone``` or download souce then ```go build main```. move binary to somewhere in your path. i reccommend ```/usr/local/bin```

### options
```-a```	path to file containing expected responses. one per line

```-t```	path to file containing tests. one per line

```-pr```	show the pass rates for each test

```-af```	show tests that all programs failed
