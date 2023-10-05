# A Simple Go Server connected to a PostgreSQL server
I’ve always found interest in the ascending Google language GO. The main reasons are well documented in the following links:

https://medium.com/@kevalpatel2106/why-should-you-learn-go-f607681fad65

https://www.freecodecamp.org/news/here-are-some-amazing-advantages-of-go-that-you-dont-hear-much-about-1af99de3b23a/ (Kudos for the Rick and Morty reference!)

https://www.quora.com/What-is-golang-good-for (I myself am as Vyas a Java veteran & Scala lover, so check out his answer)

I’m going to write how to build a simple Go Server which will listen on a specified port and it will write everything it catches on a connected PostgresSQL database. So, it’ll be a basic TCP listener but thanks to the goroutines we don’t have to worry about thread synchronization (java) or forking process (C\C++). It’ll be astonishing simple.

It’s still a raw outline and many improvements could be made (I’ll make them in following posts) but it’s a really a good start.

###First things first

Install GO Lang or use a docker go container to compile your code. Since there are tons of good guides out there that will teach you how to do both of these task, I’m linking you a couple of them but feel free to google “golang installation guide” or “compile go program in a go docker container”:

https://golang.org/doc/install

https://medium.com/@rrgarciach/bootstrapping-a-go-application-with-docker-47f1d9071a2a

We are going to use the standard Go folder structure:

https://github.com/golang-standards/project-layout

First we create the main project folder then we add internal and cmd folders

Create a folder server inside the internal folder.

Create a file app.go or main.go file inside the cmd folder. The file will be used as entrypoint for the application.

Create a server.go inside the internal/server folder.

The go package system

If you are familiar with Java, C or C++ (but also Python, Javascript, nodeJS) packaging system all you have to know is that Golang use a similar mechanism. More infos could be found at: https://medium.com/rungo/everything-you-need-to-know-about-packages-in-go-b8bac62b74cc-.

Since we are using PostgreSQL connector we need to install it:

https://golang.org/cmd/go/#hdr-Download_and_install_packages_and_dependencies

run the following command:

go get github.com/lib/pq

Now using your favorite editor (I’m using Visual Studio Code but you could use any text editor you like.) let’s write some code!

--- /internal/server/server.go ---

//the main package is mandatory for the entrypoint

package server

 

import (

    "bufio"

    "database/sql"

    "fmt"

    "log"

    "net"

    "os"

    "strings"

 

    "../logging"

    _ "github.com/lib/pq"

)

 

const (

    CONN_PORT  = "9999"

    CONN_TYPE  = "tcp"

    DB_HOST    = "localhost"

    DB_PORT    = 5432

    DB_USERNAME = "username"

    DB_PASSWORD = "password"

    DB_NAME    = "dbname"

)

 

/*

*   This function will create a new server instance listening on the

    CONN_PORT with a connection type of CONN_TYPE [tcp or udp]

*/

func StartServer(writeToDb bool) {

    l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)

    if err != nil {

        fmt.Println(err)

      //comment out the return statement: the server will continue in case of

// error. Be sure to understand the implications

     //

        //return

    }

    defer l.Close()

    //rand.Seed(time.Now().Unix())

 

    for {

        c, err := l.Accept()

        if err != nil {

            fmt.Println(err)

            return

        }

        go handleConnection(c, writeToDb)

    }

}

 

// Handles incoming requests.

/*

    TO keep the function from halting the main server thread/process

    it's advisable to call as a go routine

    go handleConnection( the connection, boolean representing if it has to write to db or not the message received)

*/

func handleConnection(c net.Conn, writeToDb bool) {

    

    fmt.Printf("Serving %s\n", c.RemoteAddr().String())

    for {

        netData, err := bufio.NewReader(c).ReadString('\n')

        if err != nil {

            fmt.Println(err)

            //return

        }

 

        temp := strings.TrimSpace(string(netData))

        if temp == "STOP" {

            go logToFile(temp)

            break

        } else {

            go logToFile(temp)

            //err = testDbConnection

            /*if err != nil {

                panic(err)

            }else{

                go writeToDB(temp)

            }*/

            if writeToDb == true {

                //testDbConnection()

                go writeToDB(temp)

            }

        }

 

        //result := strconv.Itoa(random()) + "\n"

        c.Write([]byte(string("RESP: OK\n")))

    }

    c.Close()

}

 

/*

    Handle the logging to file process

 

    TO keep the function from halting the main server thread/process

    it's advisable to call as a go routine

 

    logToFIle(a string UTF-8 representation of the message received)

*/

func logToFile(message string) {

    standardLogger := logging.NewLogger()

    

    f, err := os.OpenFile("server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

    if err != nil {

        log.Fatalf("error opening file: %v", err)

    }

    defer f.Close()

 

    standardLogger.SetOutput(f)

    standardLogger.MessageReceived(message)

}

/*

    Testing the DB CONNECTION

*/

func testDbConnection() {

    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+

        "password=%s dbname=%s sslmode=disable",

        DB_HOST, DB_PORT, DB_USERNAME, DB_PASSWORD, DB_NAME)

 

    db, err := sql.Open("postgres", psqlInfo)

    if err != nil {

        log.Fatal(err)

    }

    err = db.Ping()

    if err != nil {

        log.Fatal(err)

    }

}

 

/*

    Handle the write to database process

 

    TO keep the function from halting the main server thread/process

    it's advisable to call as a go routine

 

    writeToDB(a string UTF-8 representation of the message received)

    WARNING: many db have to be configured to accept UTF-8 4-byte ecoded ["/uffff"]

    while computers do this netively, so keep extra care when checking this

*/

func writeToDB(message string) {

    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+

        "password=%s dbname=%s sslmode=disable",

        DB_HOST, DB_PORT, DB_USERNAME, DB_PASSWORD, DB_NAME)

 

    db, err := sql.Open("postgres", psqlInfo)

    if err != nil {

        panic(err)

    }

 

    err = db.Ping()

    if err != nil {

        panic(err)

    }

    sqlStatement := `INSERT INTO log(data)

    VALUES ($1) RETURNING id`

    id := 0

    err = db.QueryRow(sqlStatement, message).Scan(&id)

    if err != nil {

        panic(err)

    }

    fmt.Println("New record ID is:", id)

}

 

 

-- /cmd/main.go –

package main

 

import "../internal/server"

 

func main() {

    //server.StartServer(bool writeToDB)

    server.StartServer(false)

}

 

Now to build the executable:

go build -o .\bin\simple_go_server.exe .\cmd

 

And it’s all done!

Let's try it.

Open a telnet terminal:
No alt text provided for this image

And see that the line is recorded to DB:
No alt text provided for this image


Simple. Fast. Efficent and portable!
