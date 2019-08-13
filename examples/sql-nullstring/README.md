# Go GraphQL SQL null string example

<a target="_blank" rel="noopener noreferrer" href="https://golang.org/pkg/database/sql/#NullString">database/sql Nullstring</a> implementation, with JSON marshalling interfaces.

To run the program, go to the directory  
`cd examples/sql-nullstring`

Run the example  
`go run main.go`

## sql.NullString

On occasion you will encounter sql fields that are nullable, as in

```sql
CREATE TABLE persons (
    id INT PRIMARY KEY,
    name TEXT NOT NULL,
    favorite_dog TEXT -- this field can have a NULL value
)
```

For the struct

```golang
import "database/sql"

type Person struct {
    ID          int             `json:"id" sql:"id"`
    Name        string          `json:"name" sql:"name"`
    FavoriteDog sql.NullString  `json:"favorite_dog" sql:"favorite_dog"`
}
```

But `graphql` would render said field as an object `{{ false}}` or `{{Bulldog true}}`, depending on their validity.

With this implementation, `graphql` would render the null items as an empty string (`""`), but would be saved in the database as `NULL`, appropriately.

The pattern can be extended to include other `database/sql` null types.
