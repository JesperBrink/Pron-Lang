# Pron-Lang
Pron is a small and simple programming language with an intuitive syntax that is easy to understand and use right away.

## Getting Started with Pron

### Prerequisites
The Pron interpreter is written in The Go Programming Language. Therefore you will have to install Go before you can run Pron. **[Install Go](https://golang.org/dl/)**

### Running Pron
You can download the executable file called 'pron' from the project by downloading the 'pron.zip'. You can then write the following in your terminal:
```text
$ ./pron
```
This starts Pron in your terminal. You can stop it again by typing `quit`. Its also possible to download the whole project from github if you want and then run the main.go file with the following command: 
```text
$ go run main.go 
```
If you want to execute a pron file, you just add the filename after the command: 
```text
$ go run main.go filename.pron
```
or you can type this equivalent command if you are using the executable file:
```text
$ ./pron filename.pron
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
In Pron there are two types of For Loops. The first increments or decrements a local variables by one each iteration and the other runs through every element in an array.
```go
// Increment
for (i from 0 to 5) {
    print(i)
}
// Prints:
// 0
// 1
// 2
// 3 
// 4

// Decrement
for (i from 5 to 0) {
    print(i)
}
// Prints:
// 5
// 4
// 3
// 2 
// 1

// Iterate through array
var arr = ["Hey", true, 42, 3.14159265359]
for (elem in arr) {
    print(elem)
}
// Prints:
// "Hey"
// true
// 42
// 3.14159265359
```

### If/Else statements
If/Else statements in Pron is as follows:
```go
if (10 > 5) {
    print("10 > 5")
} elif (10 < 5) {
    print("10 < 5")
} else {
    print("10 == 5")
}
```
You can add as many `elif` statements as you want.

### Functions
In Pron you define a function in one of the following two ways:
```go
// Define add
func add(x, y) {
    return x + y
}

// Define sub
var sub = func(x, y) {
    return x - y
}

// Call functions
add(5,2)
sub(5,2)
```
Both ways does the same. The second way just shows that it is possible to save functions in variables in Pron. Its therefore also possible to save functions inside an array or map.

### Classes
In Pron you define a class as follows:
```go
class Person{
    var name

    Init(this.name) {
        /* Do something when an Person object is being initialized*/
    }

    func GetName() {
        return name
    }

    func SetName(name) {
        this.name = name
    }
    
    func somePrivateMethod() {
        /* Do something useful */
    }    
}
```
The Init function in Pron is the constructor. The `this.name` is a short way of taking an argument `name` and then writing `this.name = name`. Pron automatically knows that you want `name` initialized to this argument.
To indicate that a class method is public, make the first letter in the method upper case. Otherwise it will be a local method.

### Builtin Functions

* `print(content)` - prints the content you give as an argument to the terminal

#### Arrays
* `len(array)` - returns the number of elements in the array
* `first(array)` - return the first element in the array
* `last(array)` - return the last element in the array
* `rest(array)` - return a copy of array without the first element
* `add(array, elem)` - returns a copy of array with elem added to it
* `remove(array, idx)` - returns a copy of array without the element at index idx

#### Maps
* `add(map, key, value)` - returns a copy of map with the key and value added to it
* `remove(map, key)` - returns a copy of array without the key/value pair associated with the key argument given

### Comments
```go
/* This is a comment in Pron */
```

## License
This project is licensed under the MIT License - see the LICENSE.md file for details

## Special Thanks
Pron could not have been written without the help of the fantastic book, [Writing an Interpeter in Go](https://interpreterbook.com) by Thorsten Ball. If you are interested in learning the process of writing an interpreter, I highly recommend you to read his book. 
