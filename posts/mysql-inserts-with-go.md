# MySQL inserts with Go

See [MySQL selects with Go](http://mycodesmells.com/post/mysql-selects-with-go/).

In the previous post we saw how to make queries to MySQL database using a small Go application. But your app will rarely be limited to read-only mode, so we obviously need to know how to save stuff to the database as well. Don't worry, it's even easier than selects!

### It's just that simple

After having your application set up according to my previous post, you will probably find making inserts into database almost trivial. The only thing you need to know is an SQL syntax and what data you want to persist.

We start, as last time, with connecting to the database:

    // main.go
    ...
    import (
        "database/sql"
        _ "github.com/go-sql-driver/mysql"
    )
    ...
    connectionString := "localhost:3306/sample"
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
        ...
	}
	defer db.Close()
    ...
    
Now, let's say that we have a table in our database to store information about people: first name, last name, age and information whether they are cool (true/false). For easier data handling on our app's side, we should create a struct:

    // main.go
    ...
    type Person struct {
        FirstName string
        LastName  string
        Age       int
        Cool      bool
    }
    ...

Now, all we need to do is, once we have some `Person` struct instance, insert its fields values into an appropriate SQL statement:

	var p = Person{FirstName: "Jerry", LastName: "West", Age: 78, Cool: true}
	res, err := db.Query("INSERT INTO People(first_name, last_name, age, cool) VALUES (?, ?, ?, ?)", p.FirstName, p.LastName, p.Age, p.Cool)
	if err != nil {
		panic(fmt.Errorf("failed to persist data: %+v\n", err))
	}
	defer res.Close()

	fmt.Println("Data persisted successfully!")
    
As you can see, there is no explicit output from our `INSERT` command, but we need do check for errors to make sure that the operation has finished wthout any problems.

Now let's try to break this code. How about adding another field to our struct, for example a flag that informs whether a person is smart or not (`Smart`). We obviously need to have a matching column in our SQL table. But what happens if there is no such column? Trying to insert such data would result in the following console output:

    $ go run main.go
    panic: failed to persist data: Error 1054: Unknown column 'smart' in 'field list'

As you can see, inside `err` variable we receive a human-readable information about an SQL error, which we need to fix. Thanks to that, saving data using Go applications is very easy, isn't it?
