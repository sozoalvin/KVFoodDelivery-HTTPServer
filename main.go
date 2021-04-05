// go run main.go database.go databaseMaps.go mainFunctions.go messageTemplates.go rawData.go queueManagement.go
/*Prototyping app for Go Assignment 2*/

package main

import (
	"fmt"
	"net/http"
	"sync"
	"text/template"
	"time"
)

var MyFoodListDB = InitMyFoodList()
var TransID int = 0
var QueueID int = 500
var pid int = 0
var SysQueue = InitSysQueue()
var searchResult []string
var selectedProduct string //selected product is always changed.

var searchResult2 []searchResultFormat
var mutex sync.Mutex

type searchResultFormat struct {
	FoodName string
	PID      string
}

type cartWeb struct {
	PID      string
	Quantity string
}

type cartDisplay struct {
	FoodName   string
	Quantity   string
	UnitPrice  float64
	TotalPrice float64
}

var cartDisplayList []cartDisplay //add new items = append slice

var cartListWeb []cartWeb // everytime you add an order, you add it in cartWeb format that is appended into cartListWeb

type CustomerData struct { //every customer data must be populated this way.
	Username   string  //their username
	Password   string  //their password
	Email      string  //their email
	TotalSpend float64 //their total spend with the kay cafe
}

type FoodInfo struct { //struct type because we need to hold values

	FoodName         string
	MerchantName     string
	DetailedLocation string
	PostalCode       int
	Price            float64
	OpeningPeriods   OpeningPeriods
}

type FoodInfo2 struct { //struct type because we need to hold values

	FoodName         string
	MerchantName     string
	DetailedLocation string
	PostalCode       int
	Price            float64
	OpeningPeriods   OpeningPeriods
}

type ShoppingCartOrder struct {
	FoodName     string
	MerchantName string
	Quantity     float64
	Price        float64
	// OrderID      int
	Username string
}

type Checkout struct {
	FoodName     string
	MerchantName string
	Quantity     float64
	Price        float64
	OrderID      string
	Username     string
}

type KVorder struct {
	transID       []string
	username      string
	systemQueueID string
	// priority      int
}

type OpeningPeriods map[string][]string

type user struct {
	UserName        string
	Password        []byte
	First           string
	Last            string
	Role            string //change use role to struct down the road since we need to check for admin/rider/supervisor/dipsatch features
	SearchLogs      *userSearchActivity
	CartList        []cartData //this should be updated based on the user's cart items. ITEMS ONLY. does not include other things
	CartDisplay     []cartDisplayData
	CartMapData     map[string]*cartFullData        //cart is only for 'staging'. key is PID. gets deleted as soon as checkout is initiatied
	CheckoutMapData map[string]*checkoutMapDataFull //systemQUEUEid is key
	cartTransID     []string
	// SearchResult {}string
}

type cartData struct {
	PID      string
	Quantity string
}

type cartDisplayData struct {
	FoodName  string
	Quantity  string
	UnitPrice float64
	TotalCost float64
}

type cartFullData struct {
	PID       string
	FoodName  string
	Quantity  string
	UnitPrice float64
	TotalCost float64
	UserRole  string
}

type checkoutMapDataFull struct {
	UserName      string
	UserRole      string
	Time          time.Time
	TransID       []string //every checkout can have a few accompanying transaction IDs
	CheckoutID    string   //every checkout will have a UNIQUE ID
	PriorityIndex int
}

type sysQueueMapDataFull struct { //checkoutMapDataFull
	UserName      string
	UserRole      string
	Time          time.Time
	TransID       []string //every checkout can have a few accompanying transaction IDs
	CheckoutID    string   //every checkout will have a UNIQUE ID
	PriorityIndex int
	DriverName    string
}

var dbtransIDSystem = map[string]*transFullData{}        //declares dbtransIDsystem as global variable. MADE MAP but with no values inside
var dbsysQueueSystem = map[string]*sysQueueMapDataFull{} //declares dbsysQueueSystem as global variable. MADE MAP but with no values inside

// var transID int = 500

type transFullData struct {
	// PID string
	UserName  string
	FoodName  string
	Quantity  string
	UnitPrice float64
	TotalCost float64
	Time      time.Time
}

type checkoutDisplay struct {
	// PID string
	// UserName  string
	FoodName  string
	Quantity  string
	UnitPrice float64
	TotalCost float64
}

type userSearchActivity struct {
	SearchKeyword     string // everytime user searches for something, information should be updated in their profile DB
	SearchResults     []string
	Time              time.Time //everytime user seaches for sometime, time stamp has to be included so we can market to them since differemt marketing agencies use different x no. of days to retarget users
	LastViewedProduct string
}

type session struct {
	un           string
	lastActivity time.Time //consider logging duration of each session tied to each user. This provides a KPI to measure customer sit engagement
}

var tpl *template.Template
var dbUsers = map[string]*user{}      // user ID, user of custom type
var dbSessions = map[string]session{} // session ID, session of custom type
var dbSessionsCleaned time.Time       //

const sessionLength int = 6000 //package level variable

func init() {
	tpl = template.Must(template.ParseGlob("templates/*")) //parses the template folder the .gohtml files that are inside
	dbSessionsCleaned = time.Now()                         //ensure that the time stamp is logged when the server is started
	// cartListWeb = make([]cartWeb{})
}

func main() {

	ch := make(chan string) //create a channel called c
	// var usrnameInpt string
	// var count int = 0
	// var postalCodeInpt int
	currentTime := time.Now()

	// usernameDS := InitUsernameTrie() //inits Trie data for username DS
	go CreateFoodList(ch) //newResult is a slice that is being returned byCreateFoodList function
	fmt.Println("\nSystem Message :", <-ch)
	go CreateFoodListMap(ch)
	go MyFoodListDB.PreInsertTrie(FoodMerchantNameAddress, ch) //populates Trie Data for Food LIst
	fmt.Println("System Message :", <-ch)
	fmt.Println("System Message :", <-ch)
	myPostalCodesDB := InitPostalCode()   //creates PostalCode BST DB
	myPostalCodesDB.PreInsertPostalCode() //preinset POSTAL Code DB
	FoodMerchantNameAddressProductID()
	fmt.Println("System Message : System is Ready", currentTime.Format("2006-01-02 15:04:05"))

	http.HandleFunc("/", index)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/assets/rs2.png", rs2)
	http.HandleFunc("/searchresult", searchresult)
	http.HandleFunc("/yourcart", yourcart)
	http.HandleFunc("/viewall", viewall)
	http.HandleFunc("/checkout_processing", checkout_processing)
	http.HandleFunc("/checkout", checkout)
	http.HandleFunc("/allsystemorders", allsystemorders)
	http.HandleFunc("/alltransactions", alltransactions)
	http.HandleFunc("/login_redirect", login_redirect)
	http.HandleFunc("/viewqueue", viewqueue)
	http.HandleFunc("/dispatchdriver", dispatchdriver)
	http.HandleFunc("/clearcart", clearcart)
	http.HandleFunc("/allthefoodisgone", allthefoodisgone)
	http.HandleFunc("/dispatchqueue", dispatchqueue)

	// http.HandleFunc("/dispatch", dispatch)
	// http.HandleFunc("/editOrder", editOrder)
	// http.HandleFunc("/productdisplay", productdisplay) //delete

	http.Handle("/favicon.ico", http.NotFoundHandler()) //NotFoundHandler returns a simple request handler that replies to each request with a “404 page not found” reply.
	http.ListenAndServe(":5221", nil)                   //launches HTTP server

	/*
		Other functins include:

		searchResults
		browseAll
		dispatch
		editOrder
	*/

	// PrintWelcomeMessage() //prints welcome message

	// fmt.Println()
	// for {
	// 	fmt.Println("\nPlease enter username to proceed\n")
	// 	fmt.Scanln(&usrnameInpt) //username checked against trie
	// 	checkUsername := usernameDS.UserSearch(usrnameInpt)
	// 	if checkUsername {
	// 		PrintUserValidated(usrnameInpt)
	// 		break
	// 	} else {
	// 		count++
	// 		PrintUserNotValidated(usrnameInpt)
	// 		if count >= 5 {
	// 			PrintNoOfTriesExceeded()
	// 			break //tested. breaks the main looping for loop
	// 		}
	// 	}
	// }
	// for { //display main meu and loop until you die.
	// 	switchRS, err := displayMainMenu()
	// 	if err != nil {
	// 		fmt.Println("System Error:", err)
	// 	} else {
	// 		switch switchRS {
	// 		case 1:
	// 			for {
	// 				fmt.Println("case 1 Access all food items")
	// 				case1result := Case1DisplayAllFoodItems(FoodMerchantNameAddress)
	// 				if ToQuit(case1result) {
	// 					fmt.Println("Returned to previous menu!")
	// 					break
	// 				}
	// 			}
	// 		case 2:
	// 			// fmt.Println("case 2 search and add item to cart")

	// 			case2rseultstring, case2rseultint, errorSearchandATC := Case2DisplayAllSearchAndATC() //call sub menu function// case2result is a string
	// 			// fmt.Println("recover to main")
	// 			if errorSearchandATC != nil {
	// 				fmt.Print("System Error:", errorSearchandATC)
	// 			} else {

	// 				case2SearchResults := MyFoodListDB.GetSuggestion(case2rseultstring, case2rseultint)
	// 				PrintKeywordSearchResults(case2SearchResults, usrnameInpt)
	// 				// fmt.Println("i got called")
	// 				AddToCart(case2SearchResults, usrnameInpt)
	// 			}
	// 			break
	// 		case 3:
	// 			// fmt.Println("case 3 check merchant postal code validity")
	// 			fmt.Println("\nEnter Postal Code to check\n")
	// 			fmt.Scanln(&postalCodeInpt)
	// 			checkPC, errPC := myPostalCodesDB.Search(postalCodeInpt)
	// 			if errPC != nil {
	// 				fmt.Println("System Error:", errPC)
	// 			} else {
	// 				if checkPC != nil {
	// 					fmt.Println("We have merchants registered at this postal code. ")
	// 				} else {
	// 					fmt.Println("We have no merchants registered at this postal code. Please advise the sales team accordingly")
	// 				}
	// 			}
	// 			break
	// 		case 4:
	// 			fmt.Println("case 4 check on current order queue")
	// 			SysQueue.PrintAllNodes()
	// 			break

	// 		case 5:
	// 			fmt.Println("Order Successfully dispatched. Please check current order queue again to check latest queues.")
	// 			SysQueue.Dequeue()
	// 			break

	// 		case 6:
	// 			fmt.Println("Display Databases")
	// 			DisplayAllDatabase()
	// 			break

	// 		case 7:
	// 			os.Exit(1)

	// 		default:
	// 			break
	// 		}
	// 	} // end switch else statement for error handling
	// }

} // close main functioin
