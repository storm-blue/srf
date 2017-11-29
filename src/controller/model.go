package controller

type Book struct {
    Name  string
    Price float32
}

type Fucker struct {
    Name string
    Age  int
}

type Response struct {
    Code    string
    Message string
    Data    interface{}
}
