# Database

This package provides four main types:

- `DBConfig`: contains the different parameters to be configured.
- `DBHolder`: creates the connection to the database and runs SQL migrations.
- `TestDBHolder`: is a wrapper of `DBHolder` that provides a `Reset()` method that cleans the database and runs again all migrations.
- `DBrepository`: is built on top of [sqlx](https://github.com/jmoiron/sqlx) to provide easier transaction management as well as methods like `Save` or `Find`.

## Usage

### DBConfig

```Go
// Returns a *DBConfig initialized by env variables
// Host                 string `env:"POSTGRES_HOST" envDefault:"localhost"`
// Port                 string `env:"POSTGRES_PORT" envDefault:"5432"`
// User                 string `env:"POSTGRES_USER" envDefault:"postgres"`
// Password             string `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
// DatabaseName         string `env:"POSTGRES_DATABASE" envDefault:"postgres"`
// SchemaName           string `env:"POSTGRES_SCHEMA" envDefault:"public"`
// MigrationsDir        string `env:"POSTGRES_MIGRATIONS_DIR" envDefault:"./migrations"`
// RunMigrationsOnReset bool   `env:"POSTGRES_RUN_MIGRATIONS" envDefault:"false"`
dbConfig := NewDBConfigFromEnv()
```

### DBHolder

```Go
// Returns a *DBHolder initialized with the provided config.
// In case the *DBConfig object has zero values, those will
// be filled with default values.
dbHolder := NewDBHolder(dbConfig)
// Run SQL migrations found in the folder specified by DBConfig.MigrationsDir
dbHolder.RunMigrations()
```

### DBrepository

```Go
type Resource struct {
    ID     string
    Name   string
    Random int
}
// This map will be used in the method Find(context.Context, url.values) to use the filters
// and sorters provided in the url.values parameter. In case the url.values contains a filter
// that it is not in the filters map, it will return an error.
filtersMap := map[string]usqlFilters.Filter{
    "id":            usqlFilters.TextField("id"),
    "name":          usqlFilters.TextField("name"),
    "random":        usqlFilters.NumField("random"),
}

sortersMap := map[string]usqlFilters.Sorter{
    "sort": usqlFilters.Sort("name", "random"),
}

r := NewDBRepository[*Resource](dbHolder, filtersMap, sortersMap)
```

#### Transactions

```Go
func DoSomething(ctx context.Context) (rErr error) {}
    // ctx must be of type context.Context
    // repository must implement the Transactional interface
    // type Transactional interface {
    //     Begin(ctx context.Context) (context.Context, error)
    //     Commit(ctx context.Context) error
    //     Rollback(ctx context.Context)
    // }
    // DBrepository already implements it.
    // The returned ctx has the transaction within it.
    ctx, err := BeginTx(ctx, repository)
    // If there is an error, the transaction could not be created.

    // Expects an *error as last parameter since it will
    // automatically perform a rollback if the function
    // finishes with an error or a commit in case there is not error.
    defer EndTx(ctx, repository, &rErr)
    // do stuff
	
    res := Resource{}
    // It's important to use the GetTransaction method instead of GetDBInstance if we want the following operations to be part of the transaction.
    tx := repository.GetTransaction(ctx)
    err = tx.NamedExecContext(ctx,
        "INSERT INTO resources (id, name, random) VALUES (:id, :name, :random);",
        &res,
    )
    // Check err

    // do more stuff
    return err // or nil
}
```

#### Search

```Go
var obj Resource
var list []*Resource

// It will return the resource found in the variable obj.
// Notice the &.
dbInstance := repository.GetDBInstance()
query := "SELECT id, name, random FROM resources"
err := repository.GetContext(ctx, dbInstance, &obj, query, url.Values{})

// Filtering by id
v := url.values{}
v.Add("id", "an_ID")
// type ResourcePage[T any] struct {
//     Total  int64 `json:"total"`
//     Limit  int64 `json:"limit"`
//     Offset int64 `json:"offset"`
//
//     In this example, *[]*Resource.
//     Resources []T `json:"resources"`
// }
// In this example, repository.SelectContext return type is ResourePage[*Resource].
resourcePage, err := repository.SelectContext(ctx, dbInstance, query, url.Values{})

// Filtering by name
v = url.values{}
v.Add("name", "the name to filter")
resourcePage, err = repository.SelectContext(ctx, dbInstance, query, v)

// Filtering by a number
v = url.values{}
v.Add("random", "4")
resourcePage, err = repository.SelectContext(ctx, dbInstance, query, v)

// Sorting by name field in ascending order
v = url.values{}
v.Add("sort", "name")
resourcePage, err = repository.SelectContext(ctx, dbInstance, query, v)

// Sorting by name field in descending order
v = url.values{}
v.Add("sort", "-name")
resourcePage, err = repository.SelectContext(ctx, dbInstance, query, v)
```

#### Build SQL

```Go
v := url.values{}
v.Add("id", "an_ID")
v.Add("limit", "10")
v.Add("sort", "name")
v.Add("sort", "-random")
// output:
//     query: "SELECT id, name, random FROM resources WHERE id = ? ORDER BY name, random desc LIMIT 10"
//     args: []interface{}{"an_ID"}
//     limit: 10
//     offset: 0
//     err: nil
query, args, limit, offset, err := repository.ApplyFilters(dbInstance, query, v)
```
