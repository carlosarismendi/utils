# Database

This package provides four main types:

- `DBConfig`: contains the different parameters to be configured.
- `DBHolder`: creates the connection to the database and runs SQL migrations.
- `TestDBHolder`: is a wrapper of `DBHolder` that provides a `Reset()` method that cleans the database and runs again all migrations.
- `DBrepository`: is built on top of [GORM](https://gorm.io/) to provide easier transaction management as well as methods like `Save` or `Find`.

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
dbHolder := NewDHolder(dbConfig)
// Run SQL migrations found in the folder specified by DBConfig.MigrationsDir
dbHolder.RunMigrations()
```

### DBrepository

```Go
// This map will be used in the method Find(context.Context, url.values) to use the filters and sorters provided in the url.values parameter. In case the url.values contains a filter that it is not in the filters map, it will return an error.
filters := map[string]filters.Filter{
    "id":            filters.TextField("id"),
    "name":          filters.TextField("name"),
    "random_number": filters.NumField("random_number"),
    "sort":          filters.Sorter(),
}
repository := NewDBRepository(dbHolder, filters)
```

#### Transactions

```Go
func DoSomething(ctx context.Context) (rErr error) {}
    // ctx must be of type context.Context
    // repository must implement the Transactional interface
    // type Transactional interface {
    //     Begin(ctx context.Context) (context.Context, error)
    //     Commit(ctx context.Context) error
    //     Rollback(ctx context.Context) error
    // }
    // DBrepository already implements it.
    // The returned ctx has the transaction within it.
    ctx, err := BeginTx(ctx, repository)
    // If there is an error, the transaction could not be created.

    // Expects an *error as last parameter since it will
    // automatically perform a rollback if the function
    // finishes with an error Or a commit in case there is not error.
    defer EndTx(ctx, repository, &rErr)
    // do stuff

    type Resource struct {
        ID     string
        Name   string
        Random int
    }
    err = repository.Save(ctx, &Resource{
        ID: "an_ID",
        Name: "Name"
        Random: 4,
    })
    // Check err

    // do more stuff
    return err // or nil
}
```

#### Find

```Go
var obj Resource
var list []*Resource

// It will return the resource found in the variable obj.
// Notice the &.
err = repository.FindByID(ctx, "an_ID", &obj)

// Filtering by id
v := url.values{}
v.Add("id", "an_ID")
// It is necessary to pass the list parameter so
// internally can infer the type and table to use to
// request the data.
// resourcePage is of type:
// type ResourcePage struct {
// 	   Total  int64 `json:"total"`
// 	   Limit  int64 `json:"limit"`
// 	   Offset int64 `json:"offset"`
//
//     // Resource will be a pointer to the type pased as
//     // dst parameter in Find method. In this example,
//     // *[]*Resource.
//     Resources interface{} `json:"resources"`
// }
resourcePage, err = repository.Find(ctx, v, list)

// Filtering by name
v2 := url.values{}
v2.Add("name", "the name to filter")
resourcePage, err = repository.Find(ctx, v, list)

// Filtering by a number
v2 := url.values{}
v2.Add("random", "4")
resourcePage, err = repository.Find(ctx, v, list)

// Sorting by name field in ascending order
v2 := url.values{}
v2.Add("sort", "name")
resourcePage, err = repository.Find(ctx, v, list)

// Sorting by name field in descending order
v2 := url.values{}
v2.Add("sort", "-name")
resourcePage, err = repository.Find(ctx, v, list)
```
