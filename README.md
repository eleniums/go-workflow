# go-workflow
In-memory tool for managing program flow. Allows chaining functions and making decisions. Useful for shorter lived workflows and is not intended to be durable.

NOTE: This was created as an experiment and is most likely overkill for any real-world application.

## Usage
Start by creating some functions to do some work. A function can accept a single input and returns a single output with an error. If an error occurs, no further functions are executed and the error is returned along with the output from the function that returned the error.

Here's some sample basic functions:
```go
add1 := func(in int) (int, error) {
    return in + 1, nil
}

add2 := func(in int) (int, error) {
    return in + 2, nil
}

add3 := func(in int) (int, error) {
    return in + 3, nil
}
```

String functions together to create a workflow that compiles down to a single function in the end:
```go
action := Sequential(
    Do(add1),
    Parallel(sum,
        Do(add1),
        Do(add2),
        Sequential(
            Do(add1),
            Do(add2),
            If(isOdd,
                Do(add2),
                Do(add3),
            ),
        )),
    If(isOdd,
        NoOp(),
        Do(add2),
    ),
)
```

Calling the function will execute the entire workflow:
```go
result, err := action(1)
```

For further examples, look at the unit tests.

## Functions
- `Do`: Perform an action. Takes a function and wraps it to the Action type.
- `Sequential`: Perform some actions in sequence.
- `Parallel`: Perform some actions in parallel.
- `If`: Conditionally perform one action or another.
- `NoOp`: Does nothing. Useful as a dead end.

## TODO
- [ ] Add Catch function to handle an error.
- [ ] Add Retry function to add retry functionality. Perhaps use Catch function internally.
