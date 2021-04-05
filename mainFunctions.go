package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var FoodMerchantNameAddress []string
var FoodMerchantNameAddress1 []string
var KVorderforQueue []KVorder

func checkUsernameStructure(s string) bool {

	for _, v := range s {
		switch string(v) {

		case "0":
			fallthrough
		case "1":
			fallthrough
		case "2":
			fallthrough
		case "3":
			fallthrough
		case "4":
			fallthrough
		case "5":
			fallthrough
		case "6":
			fallthrough
		case "7":
			fallthrough
		case "8":
			fallthrough
		case "9":
			fmt.Println("\nYour username should not contain any numbers. Please try again")
			return false
		default:
			continue
		}
	}
	return true
}

func displayMainMenu() (int, error) {

	var usrInpt int
	fmt.Println("\nPlease select from the following menu\n")
	fmt.Println("1. Access Food Items")
	fmt.Println("2. Search and Add Items to Cart")
	fmt.Println("3. Check Merchant Postal Code Validity")
	fmt.Println("4. Check on Current Order Queue")
	fmt.Println("5. Dispatch Order")
	fmt.Println("6. Display Databases")
	fmt.Println("7. End Program")
	fmt.Scanln(&usrInpt)

	if usrInpt <= 0 || usrInpt > 7 {
		return -1, errors.New("Input cannot be negative or more than the number of options provided")
	}
	return usrInpt, nil
}

func concatenateFoodList() { //return you a slice of strings
	fmt.Println(V)
}

func CreateFoodList(ch chan string) {
	// func CreateFoodList() []string {

	for _, v := range V {
		FoodMerchantNameAddress = append(FoodMerchantNameAddress, v.FoodName+" - "+v.MerchantName+" - "+v.DetailedLocation)
	}
	sort.Strings(FoodMerchantNameAddress)
	ch <- "Mandatory - Food List Data Generated"
	// return FoodMerchantNameAddress
}

func CreateFoodList1() []string { //function for template testing.

	for _, v := range V_2 {
		FoodMerchantNameAddress1 = append(FoodMerchantNameAddress1, v.FoodName+" - "+v.MerchantName+" - "+v.DetailedLocation)
	}
	sort.Strings(FoodMerchantNameAddress1)
	return FoodMerchantNameAddress1
}

func PrintSliceinLines(s []string) {
	sum := 0

	fmt.Println("You are now browsing all available items regardless of your postal code")

	for i, v := range s {
		fmt.Printf("%d. %s\n", i+1, v)
		sum++
	}
	fmt.Printf("\nThere are a total of %d food items available for order\n", sum)
	fmt.Println("You can order them in the add to cart option")
}

func Case1DisplayAllFoodItems(s []string) string {
	var usrInpt string
	fmt.Println("You are now browsing all available items regardless of your postal code")
	fmt.Println("Please use the search function on the previous menu if you need to search for something\n")
	PrintSliceinLinesFoodListGeneral(s)
	fmt.Println("\nEnter the letter Q to exit to previous menu")
	fmt.Scanln(&usrInpt)
	return usrInpt
}

func PrintSliceinLinesFoodListGeneral(s []string) {

	for i, v := range s {
		fmt.Printf("%d. %s\n", i+1, v)
	}
}

func PrintSliceinLinesGeneral(s []string) {

	for i, v := range s {
		fmt.Printf("%d. %s\n", i+1, v)
	}

}

func PrintSliceinLinesGeneralSearch(s []string) int {

	for i, v := range s {
		fmt.Printf("%d. %s\n", i+1, v)
	}

	return len(s)
}

func ToQuit(s string) bool { //use this function when u want to check for input = q . usually used to return to main menu

	r := strings.ToLower(s)
	if r == "q" {
		return true
	}
	return false
}

func Case2DisplayAllSearchAndATC() (string, int, error) {

	var usrInpt string
	var usrInpt2 int

	fmt.Println("\nPlease Search for an Item(lower case alphabets, including numbers 0-9)\n")
	fmt.Scanln(&usrInpt)
	fmt.Println("\nPlease Indicate the Number of Similar Results You Want to Display\n")
	fmt.Scanln(&usrInpt2)
	if usrInpt2 <= 0 {
		return "", -1, errors.New("Number of Search Results to Display Should be a positive number. i.e. 10 ")
	}
	return usrInpt, usrInpt2, nil

}

func AddToCart(s []string, un string) { // s is the search result
	var usrInptChoice string
	var usrInptQty string
	fmt.Println("\nPlease Indicate the item you would like to place an order for:")
	for {
		fmt.Scanln(&usrInptChoice)
		x, _ := strconv.Atoi(usrInptChoice)
		if x <= 0 || x > len(s) {
			fmt.Println("\nInvalid Input. Please try again. Please enter your choice again")
			// return
		} else {
			break
		}
	}
	fmt.Println("\nPlease Indicate the quantity of the item you would like to place an order for:")
	fmt.Scanln(&usrInptQty)
	MatchUsrInptToSlice(s, usrInptChoice, usrInptQty, un) //item Properties should return you item FULL NAME and price with reference to map values

}

func MatchUsrInptToSlice(s []string, s1 string, s2 string, un string) (string, float64) {

	// noOfChoices := len(s) //if search results only have 2 elements; this returns 2.
	s1n, _ := strconv.Atoi(s1)
	s1nn := s1n - 1
	// fmt.Println(s1nn)
	s2n, _ := strconv.ParseFloat(s2, 32)

	// var usrInpt int

	fmt.Println("\n===================================================\n")
	fmt.Printf("\nThis Order Will be Tagged to User: %s. Please exit program to relogin if your login ID is incorrect\n", un)
	fmt.Printf("\nYour Option: %s Has Been Selected\n\n", s[s1nn])
	fmt.Printf("The Name of the Selected Food Item is: %s\n", MyFoodListMap[s[s1nn]].FoodName)
	fmt.Printf("\nThe Total Quantity Ordered: %0.0f\n", s2n)
	fmt.Printf("\nThe Total Cost for this is: $%0.2f\n", s2n*MyFoodListMap[s[s1nn]].Price)
	fmt.Println("\nPlease select save shopping cart if you want to checkout the above items. ")

	fmt.Println("\n===================================================")

	fmt.Println("Please choose from the following options\n")
	SearchSaveCheckOut(un, s[s1nn], s2n, s2n*MyFoodListMap[s[s1nn]].Price, MyFoodListMap[s[s1nn]].MerchantName, MyFoodListMap[s[s1nn]].FoodName) //passes username, the complete food dish, the total quantity, the price, the merchant name, the food name

	return "", 0
}

func SearchSaveCheckOut(un string, s string, d float64, f float64, s2 string, s3 string) { //passes username, the complete food dish, the total quantity, the price, the merchant name, the food name

	var usrInpt int

	fmt.Println("1. Search for new item to add to cart")
	fmt.Println("2. Update Shopping Cart")
	fmt.Println("3. Clear Shopping Cart")
	fmt.Println("4. Checkout Shopping Cart")
	fmt.Println("5. Return to Previous Menu")
	fmt.Scanln(&usrInpt)

	switch usrInpt {
	case 1:
		case2rseultstring, case2rseultint, errorSearchandATC := Case2DisplayAllSearchAndATC()
		if errorSearchandATC != nil {
			fmt.Print("System Error:", errorSearchandATC)
		} else {
			case2SearchResults := MyFoodListDB.GetSuggestion(case2rseultstring, case2rseultint)
			PrintKeywordSearchResults(case2SearchResults, un)
			AddToCart(case2SearchResults, un)
		}
		break

	case 2:
		AddOrdertoCart(un, s, d, f, s2, s3) //append item to cart list
		SearchSaveCheckOut(un, s, d, f, s2, s3)
		break

	case 3:
		ClearShoppingCartAndCheckoutinformation()
		break

	case 4:
		CheckoutConfirm(un) //transform MyShoppingCart into checkout with transactionID
		break

	case 5:
		return

	}
	return
}

func CheckoutConfirm(s string) { //s in this case is your username

	fmt.Println("\n========================================================")
	fmt.Println("\n\t\tCheckout Confirmed")
	fmt.Println("\n========================================================\n")

	var PriorityIndex = 1
	_ = PriorityIndex //declared but should be removed if program is in written in terminal
	var usrInpt int
	totalTransPerSession := []string{} //everytime this function is freshly called, total transaction per session resets to 0 which is correct.

	for _, v := range MyShoppingCart {
		TransID++                         //everytime there is a new shopping cart; we need to have a new transactioID for EACH registered mercahnt
		ns := strconv.Itoa(TransID)       //converts transaction ID to string
		transactionID := "MC" + ns + "KV" //generates transaction ID
		MyCheckoutTranID[transactionID] = Checkout{v.FoodName, v.MerchantName, v.Quantity, v.Price, transactionID, s}
		MyCheckoutIDUsername[s] = Checkout{v.FoodName, v.MerchantName, v.Quantity, v.Price, transactionID, s}
		totalTransPerSession = append(totalTransPerSession, transactionID) //if there 5 different transactions for 5 different merchats, all will be appended to this variable
	}
	QueueID++
	ns2 := strconv.Itoa(QueueID)       //converts queue ID to string
	systemQueueID := "OS" + ns2 + "KV" //generates system queue ID. NOT queue for merchant but on overall system level

	checkUserAdminResult := checkUserAdmin(s) //returns T or F
	if checkUserAdminResult {                 //this allows an admin user to push a queue to priority
		fmt.Printf("\nHi user: %s,! You are authorised as a customer service officer!\n\n", s)
		fmt.Println("Please enter priority number more than 0 if an order is an order is meant for service recovery")
		fmt.Println("Enter 0 for default if is not for service recovery")
		fmt.Scanln(&usrInpt)
		PriorityIndex = usrInpt
	}
	KVorderforQueue = append(KVorderforQueue, KVorder{totalTransPerSession, s, systemQueueID})
	// SysQueue.Enqueue(KVorderforQueue, PriorityIndex)
	fmt.Println("\nOrder successfully processed. What would you like to do next?\n")
	ClearShoppingCartAndCheckoutinformation()
	return
}

// func checkUserAdmin(s string) bool {
// 	for _, v := range UsernameList2 {
// 		if v.UserName == s {
// 			if v.isAdmin == true {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

func checkUserAdmin(s string) bool {
	// for _, v := range dbUsers {
	// 	if v.UserName == s {
	// 		if v.isAdmin == true {
	// 			return true
	// 		}
	// 	}
	// }
	// return false

	if v, ok := dbUsers[s]; ok {
		if v.Role == "admin" {
			return true
		}

	}
	return false
}

func AddOrdertoCart(un string, s string, d float64, f float64, s2 string, s3 string) { //passes username, the complete food dish, the total quantity, the price, the merchant name.

	ok := RepeatedShoppingCartCheck(s2, s3) // check for 100% && match with FOOD NAME as well as merchant name
	if ok {
		MyShoppingCart = append(MyShoppingCart, ShoppingCartOrder{s3, s2, d, f, un})
		fmt.Println("Order Successfully Added to Cart!")
	} else {
		fmt.Println("\nCart Already contains similar items. Please consider adding something else into cart")
	}
	PrintShoppingCartItems(un, s, d, f, s2, s3)
}

func RepeatedShoppingCartCheck(s2 string, s3 string) bool {

	for _, v := range MyShoppingCart {
		if v.FoodName == s3 && v.MerchantName == s2 {
			return false
		}
	}
	return true
}

func PrintShoppingCartItems(un string, s string, d float64, f float64, s2 string, s3 string) { //passes username, the complete food dish, the total quantity, the price, the merchant name.
	var totalQty float64
	var totalSum float64

	if len(MyShoppingCart) == 0 {
		//do nothing
	} else {
		fmt.Println("\n========================================================")
		fmt.Println("\n\tYour Current Shopping Basket Items")
		// fmt.Println("\n================================================")
		for i, v := range MyShoppingCart {

			fmt.Println("\n========================================================")
			fmt.Printf("\n%d.\n", i+1)
			fmt.Printf("\tMerchant Name: \t\t\t%s", v.MerchantName)
			fmt.Printf("\n\n\tDish Name: \t\t\t%s", v.FoodName)
			fmt.Printf("\n\n\tQuantity: \t\t\t%0.0f", v.Quantity)
			fmt.Printf("\n\n\tUnit Price: \t\t\t%0.2f", v.Price)
			fmt.Printf("\n\n\tTotal Cost: \t\t\t%0.2f\n", v.Quantity*v.Price)
			// fmt.Printf("\n\n", s)
			// fmt.Printf("\n\nMerchant Name: %s", 2)
			totalQty += v.Quantity
			totalSum += v.Quantity * v.Price

		}
		fmt.Println("\n========================================================")
		fmt.Printf("\nTotal Number of Food Dishes Ordered: \t%0.0f\n", totalQty)
		fmt.Printf("Total Cost of all Food Dishes Ordered: \t$%0.2f\n", totalSum)
		fmt.Println("\n========================================================")

	}
}

func PrintKeywordSearchResults(ss []string, s string) {

	var usrInpt string
	if len(ss) != 0 {
		fmt.Println("\nHere are the search results:\n")
		PrintSliceinLinesGeneral(ss)

	} else { //if no serach results; we need to do something else
		// fmt.Println(ss)
		fmt.Println("There are no registered merchants that are selling the items with your search terms. \n\nPlease consider searching for another keyword. Enter S to launch the search menu")

		for {
			fmt.Scanln(&usrInpt)
			//loop until function call exit
			if strings.ToLower(usrInpt) == "s" {
				case2rseultstring, case2rseultint, errorSearchandATC := Case2DisplayAllSearchAndATC()
				if errorSearchandATC != nil {
					fmt.Print("System Error:", errorSearchandATC)
					fmt.Println("Please Enter s to Search Again")
					// return
				} else {
					case2SearchResults := MyFoodListDB.GetSuggestion(case2rseultstring, case2rseultint)
					PrintKeywordSearchResults(case2SearchResults, s)
					// fmt.Println("am i called?")
					AddToCart(case2SearchResults, s)
				}
				continue
				// break
			} else {
				fmt.Println("Invalid input. Enter S to search again")
			}
		}
	}
}

func ClearShoppingCartAndCheckoutinformation() {
	MyShoppingCart = []ShoppingCartOrder{} //clears shopping cart
	KVorderforQueue = []KVorder{}
	// return
}

func DisplayAllDatabase() {

	var usinpt int
	i := 1
	j := 1

	fmt.Println("Select to see all data of merchants using transaction ID")
	fmt.Println("\n========================================================")

	fmt.Println("1. View all data related to checkout transaction ID(s)")
	fmt.Println("2. Export all Data regarding Transaction ID(s)")
	fmt.Println("3. View all Data related to usernames")
	fmt.Println("4. Export all Data related to usernames")
	fmt.Scanln(&usinpt)
	fmt.Println("\n========================================================\n")

	switch usinpt {

	case 1:

		// MyCheckoutTranID[transactionID]
		fmt.Println("All data related to Transaction IDs\n")

		for k, _ := range MyCheckoutTranID {

			fmt.Printf("\n%d. Order ID: %s \t\t\t\n", i, k)
			fmt.Printf("\nOrdered Dish Name:\t\t\t%s\n", MyCheckoutTranID[k].FoodName)
			fmt.Printf("\nFulfiled by Merchant: \t\t\t%s\n", MyCheckoutTranID[k].MerchantName)
			fmt.Printf("\nOrdered Quantity:%0.0f\t\t\tUnit Price %0.2f\n", MyCheckoutTranID[k].Quantity, MyCheckoutTranID[k].Price)
			fmt.Printf("\nOrder tagged to username: \t\t%s\n\n", MyCheckoutTranID[k].Username)
			fmt.Printf("Total Order Value: \t\t\t%0.2f\n", MyCheckoutTranID[k].Quantity*MyCheckoutTranID[k].Price)
			fmt.Println("\n============================================================")
			i++
		}

	case 2:
		fmt.Println("Export Data regarding Transaction ID(s)\n")

	case 3:
		// MyCheckoutIDUsername[s]
		fmt.Println("All data related to usernames\n")

		for k, _ := range MyCheckoutIDUsername {
			fmt.Printf("\n\n%d. Order ID: %s \t\t\t\n", j, k)
			fmt.Printf("\nOrdered Dish Name:\t\t\t%s\n", MyCheckoutIDUsername[k].FoodName)
			fmt.Printf("\nFulfiled by Merchant: \t\t\t%s\n", MyCheckoutIDUsername[k].MerchantName)
			fmt.Printf("\nOrdered Quantity:%0.0f\t\t\tUnit Price %0.2f\n", MyCheckoutIDUsername[k].Quantity, MyCheckoutIDUsername[k].Price)
			fmt.Printf("\nOrder tagged to username: \t\t%s\n\n", MyCheckoutIDUsername[k].Username)
			fmt.Printf("Total Order Value: \t\t\t%0.2f\n", MyCheckoutIDUsername[k].Quantity*MyCheckoutIDUsername[k].Price)
			fmt.Println("\n============================================================")
			j++
		}

	case 4:
		fmt.Println("Export Data regarding Usernames\n")

	}

}

func addToCartWeb(s []string, un string) { // s is the search result
	var usrInptChoice string
	var usrInptQty string
	fmt.Println("\nPlease Indicate the item you would like to place an order for:")
	for {
		fmt.Scanln(&usrInptChoice)
		x, _ := strconv.Atoi(usrInptChoice)
		if x <= 0 || x > len(s) {
			fmt.Println("\nInvalid Input. Please try again. Please enter your choice again")
			// return
		} else {
			break
		}
	}
	fmt.Println("\nPlease Indicate the quantity of the item you would like to place an order for:")
	fmt.Scanln(&usrInptQty)
	MatchUsrInptToSlice(s, usrInptChoice, usrInptQty, un) //item Properties should return you item FULL NAME and price with reference to map values

}

func generateTransactionID() (string, error) {

	var generatedID string
	//for every element wihch is stored based on PID as key, should have unique generated ID

	mutex.Lock()
	{
		sTransID := strconv.Itoa(TransID)

		generatedID = "MC" + sTransID + "KV" // generated ID will always be unique

		// transIDSystem[generatedID] = &transFullData{u.UserName, v.FoodName, v.Quantity, v.UnitPrice, v.TotalCost, time.Now()}

		TransID++
	}
	mutex.Unlock()

	return generatedID, nil

}

func generateSysQueueID() (string, error) {

	var generatedSysQueueID string

	mutex.Lock()
	{
		SQueueID := strconv.Itoa(QueueID)

		generatedSysQueueID = "OS" + SQueueID + "KV" // generated ID will always be unique

		QueueID++
	}
	mutex.Unlock()

	return generatedSysQueueID, nil
}
