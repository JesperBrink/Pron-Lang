# Pron-Lang
Pron is a small and simple programming language with an intuitive syntax that everyone can understand and use right away.

## Getting Started with Pron

### Prerequisites
The Pron interpreter is written in The Go Programming Language. Therefore you will have to install Go before you can run Pron. **[Install Go](https://golang.org/dl/)**

### Installing
You just download the git repository and run the main.go file with the folloing command: 
```text
$ go run main.go 
```
This starts Pron in your terminal. If you want to execute a pron file, you just add the filename after the command: 
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

### Operators

#### Code example


### Arrays

#### Code example


### Maps

#### Code example


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