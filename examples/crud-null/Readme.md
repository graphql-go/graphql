# Go GraphQL CRUD example

Implementation create, read, update and delete on Go    

To run the program, go to the directory    
`cd examples/crud`

Run the example    
`go run main.go`

## Create
`http://localhost:8080/product?query=mutation+_{create(name:"Inca Kola",info:"Inca Kola is a soft drink that was created in Peru in 1935 by British immigrant Joseph Robinson Lindley using lemon verbena (wiki)",price:1.99){id,name,info,price}}`

## Read
Get single product by id    
`http://localhost:8080/product?query={product(id:1){name,info,price}}`

Get product list    
`http://localhost:8080/product?query={list{id,name,info,price}}`

## Update
`http://localhost:8080/product?query=mutation+_{update(id:1,price:3.95){id,name,info,price}}`

## Delete
`http://localhost:8080/product?query=mutation+_{delete(id:1){id,name,info,price}}`
