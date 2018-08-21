# Pron-Lang
Pron is a small and simple programming language with an intuitive syntax that everyone can understand and use right away.

## Getting Started with Pron

### Prerequisites
The Pron interpreter is written in The Go Programming Language. Therefore you will have to install Go before you can run Pron. **[Install Go](https://golang.org/dl/)**

### Installing
You just download the git repository and run the main.go file with the following command: 
```text
$ go run main.go 
```
This starts Pron in your terminal. You can stop it again by typing `quit`. If you want to execute a pron file, you just add the filename after the command: 
```text
$ go run main.go filename.pron
```
You can find some code examples in the main package of the project called 'testfile.pron' and 'TestClass.pron'.


## Documentation

### Datatypes
* string - `"Hello World!"`
* boolean - `true, false`
* int (64 bit) - `1, 2, 3, 42`
* real (64 bit) - `1.0, 2.54, 3.14159265359`

### Variables
In Pron you don't specify the datatype of your variable. You just declare it with the keyword: `var`.
```go
var helloStr = "Hello World!"
var myInt = 42
var isPronAwesome = true
var pi = 3.14159265359
var thisIsInitializedToNull
```

### Arrays
Like the variables you don't specify the type of the array. This means that you can combine anything in an array in Pron. 
```go
var emptyArr = []
var strings = ["Hello World!", "Cool", "myStr"]
var combined = [42, "Hello World!", true, 3.14159265359]

// Extract data from array
var str = strings[1] // str becomes "Cool"

// Add to array - use add(array, elemToBeAdded)
// biggerArr becomes: ["Hello World!", "Cool", "myStr", "anotherString"] while strings hasn't changed.
var biggerArr = add(strings, "anotherString")

// Remove from array - use remove(array, indexToRemoveFrom)
// smallerArr becomes: ["Hello World!", "myStr"] while strings hasn't changed
var smallerArr = remove(strings, 1)
```

### Maps
Like the variables and the arrays you don't specify the type. This means that you can have anything as values in you map. On the other hand is it only possible to use string, int and bool as you keys.
```go
var emptyMap = {}
var myMap = {"Hello": "World!", 4: 2, "T": true, 3: ".14159265359"}

// Extract data from map
var restOfPi = myMap[3] // restOfPy becomes ".14159265359"
var myBool = myMap["T"] // myBool becomes true

// Add to map - use add(map, key, value)
// biggerMap becomes: {"Hello": "World!", 4: 2, "T": true, 3: ".14159265359", "1": 1} while myMap hasn't changed.
var biggerMap = add(myMap, "1", 1)

// Remove from map - use remove(map, keyToRemove)
// smallerMap becomes: {"Hello": "World!", 4: 2, "T": true} while myMap hasn't changed
var smallerMap = remove(myMap, 3)
```

### For Loops

#### Increment

#### Decrement

#### Run through an array

#### Code example


### If/Else statements

#### Code example


### Functions
(to typer)
Ingen specifik returtype

#### Code example


### Builtin Functions

#### Code example

### Classes

#### Code example

### Comments


## License
This project is licensed under the MIT License - see the LICENSE.md file for details

## Special Thanks
Pron could not have been written without the help of the fantastic book, [Writing an Interpeter in Go](https://interpreterbook.com) by Thorsten Ball. If you are interested in learning the process of writing an interpreter, I highly recommend you to read his book. 