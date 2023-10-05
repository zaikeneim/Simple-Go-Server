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
	CONN_PORT   = "9999"
	//CONN_TYPE   = "udp"
	CONN_TYPE   = "tcp"
	DB_HOST     = "localhost"
	DB_PORT     = 5432
	DB_USERNAME = "test"
	DB_PASSWORD = "test5999"
	DB_NAME     = "test"
	CARRIAGE_RETURN = 13
)

/*
*	This function will create a new server instance listening on the
	CONN_PORT with a connection type of CONN_TYPE [tcp or udp]
*/
func StartServer(writeToDb bool) {
	//l, err := net.ListenPacket(CONN_TYPE, ":"+CONN_PORT)
	l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)
	if err != nil {
		fmt.Println(err)
		//tolto il return per poter continuare anche in caso di errore
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
		//c.Write([]byte(string("RESP: OK\n")))
	}
	c.Close()
}


/*
func handleConnection(l net.PacketConn, writeToDb bool) {

	//fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {

		buffer := make([]byte, 1024)
		//netData, err := bufio.NewReader(l).ReadString(CARRIAGE_RETURN)
		n, addr, err := l.ReadFrom(buffer)
		if err != nil {
			fmt.Println(err)
			continue
			//return
		}
		fmt.Println(addr)
		fmt.Println(n)

		//temp := strings.TrimSpace(string(netData))
		/*if temp == "STOP" {
			go logToFile(temp)
			break
		} else {*/
			//go logToFile(temp)
			//err = testDbConnection
			/*if err != nil {
				panic(err)
			}else{
				go writeToDB(temp)
			}*/
			/*if writeToDb == true {
				//testDbConnection()
				go writeToDB(temp)
			}*/
		//}

		//result := strconv.Itoa(random()) + "\n"
		//c.Write([]byte(string("RESP: OK\n")))
/*	}
	//c.Close()
}
*/

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
	Testing the DB CONNECTION to prevent it collapsing 
	to our faces
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
