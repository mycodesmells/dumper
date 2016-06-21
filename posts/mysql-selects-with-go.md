# MySQL selects with Go

Nowadays every application has to deal with data, which means your application has to be able to persist and read some information. You can do this using various data sources, such as files, SQL databases or NoSQL. This post deals with a thing I needed to do recently, that is fetching data from MySQL database using Go language.

### Idea

Recently I've decided to fetch all blog posts I created so far, and as a _lazy_ developer I thought that logging into the machine or into administration panel and making a dump is just too much. What a waste of five minutes of my time, right? Instead, I decided to spend a few minutes more to create a tiny tool that would create such dump, so that I could repeat the process any time I need.

### Necessary imports

In order to connect to MySQL database we need to import two packages into our app:

    // main.go
    ...
    import (
        "database/sql"
        _ "github.com/go-sql-driver/mysql"
    )
    
If you are just starting with Go, just as you got comfortable with imports, you might feel a bit lost with the second one. What does that underscore mean? In fact, it's not difficult to grasp at all. You should already know, that whenever you import something, you **have to use it**, as Go compiler demands you to. You can trick it, however, with an underscore - this way the package is imported, but it doesn't have to be used anywhere. This doesn't make any sense right? Let's look into the `go-sql-driver/mysql` package then:

    // go-sql-driver/mysql/driver.go
    ...
    func init() {
	    sql.Register("mysql", &MySQLDriver{})
    }
    
As you can see, there something important happens inside `init()` function, which is executed when the package is being loaded. This line of code adds MySQLDriver to the list of supported ones inside `sql` package. This way we don't need to use `mysql` package explicitly, but have to import it at some point.
    
### Connecting to the database

This is fairly simple, as long as you remember to put all the necessery parts into the conection string. The string should look like this:

    username:password@protocol(hostname:port)/dbname?param1=value1&param2=value2...
    
In order to actually connect to the database, we just need to invoke one simple function from `sql` package:

    connectionString := "localhost:3306/sample"
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
        ...
	}
	defer db.Close()
    
**Remember** about closing the connection using `defer`, and make sure you handle any errors that may occur there. Now we are ready to proceed and fetch some data.

### Fetching data into structs

The main thing we want to accomplish when selecting data from our database into Go application is to reuse it in an easy way. We would probably like to describe our table using some struct, and then store each row as its instance. So firts, let's create one:

    type Post struct {
        ID       int
        Title    string
        Slug     string
        Abstract string
        Content  string
        Dates    struct {
            Create      time.Time
            Publication time.Time
        }
    }

As you can see, despite having everything stored as a flat structure in an SQL table, we can have nested structures within our application. One thing you need to remember is to have all the struct fields exported (that is starting with a capital letter), as it will be necessary later on.

Now we need to perform some `SELECT` statements on the database. If you have some experience with Java, Python, or basically any high-level language, you may be used to using some 3rd party library for that access. Unfortunately (or is it?), Go sticks to its _do it yourself_ rule here as well, so we need to write an SQL statement ourselves:

    res, err := db.Query("SELECT id, title, slug, abstract, content, create_date, publication_date FROM posts")
	if err != nil {
        ...
	}
	defer res.Close()

Once again, we need to remember to handle errors and close data stream with `defer`. The last step we need to perform is to map all the results into struct's instances:

	for res.Next() {
        var p Post
		res.Scan(&p.ID, &p.Title, &p.Slug, &p.Abstract, &p.Content, &p.Dates.Create, &p.Dates.Publication)
        
        fmt.Printf("Saved another post: %+v", p)
	}

First of all, we need to loop over the results object until there is no more data, that's what `res.Next()` is for. Then, inside the loop we define a varaiable of given type and map incoming data into its attributes. Finally, for the sake of the example, we can print `Post` object to confirm that we fetched data successfully.

You can see this code in action [on Github](https://github.com/mycodesmells/dumper).
