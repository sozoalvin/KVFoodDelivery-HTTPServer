package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func index(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req) //getUser function call
	// fmt.Println(&u)

	if req.Method == http.MethodPost {
		sr := req.FormValue("searchtext")

		if u == nil {
			//no username present!
			http.Redirect(w, req, "/login_redirect", http.StatusSeeOther)
			return
		}

		localsearchResult := MyFoodListDB.GetSuggestion(sr, 50) // you will always append a global variable so you pass data this way.

		dbUsers[u.UserName].SearchLogs = &userSearchActivity{SearchResults: localsearchResult}

		http.Redirect(w, req, "/searchresult", http.StatusSeeOther)
	}

	showSessions() // for demonstration purposes
	tpl.ExecuteTemplate(w, "index.gohtml", u)
}

func bar(w http.ResponseWriter, req *http.Request) {
	u := getUser(w, req)
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	if u.Role != "007" {
		http.Error(w, "You must be 007 to enter the bar", http.StatusForbidden)
		return
	}
	showSessions() // for demonstration purposes
	tpl.ExecuteTemplate(w, "bar.gohtml", u)
}

func signup(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	var u user

	var me = make(map[string]string) //make map for error

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("Recovered in f", r)
	// 	}
	// }()

	if req.Method == http.MethodPost {
		// get form values
		un := req.FormValue("username")
		p := req.FormValue("password")
		f := req.FormValue("firstname")
		l := req.FormValue("lastname")
		r := req.FormValue("role")

		boolresult, mapresult := validateInputs(un, p, f, l, me)

		if boolresult == false {
			tpl.ExecuteTemplate(w, "signup.gohtml", mapresult)
			return
		}

		if _, ok := dbUsers[un]; ok {
			http.Error(w, "Username already taken", http.StatusForbidden) //code uses map dbUsers to check for username validity
			return
		}

		// create session
		sID, err := uuid.NewV4()
		//err handling
		if err != nil {
			fmt.Printf("Something went wrong: %s, err")
		}

		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}

		c.MaxAge = sessionLength
		http.SetCookie(w, c)

		dbSessions[c.Value] = session{un, time.Now()} // i wil store your informtion with cookie value UUID
		// store user in dbUsers
		bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		u = user{
			un,
			bs,
			f,
			l,
			r,
			&userSearchActivity{"nil", []string{"nil"}, time.Now(), "nil"},
			[]cartData{},
			[]cartDisplayData{},
			nil,
			nil,
			[]string{},
		}

		dbUsers[un] = &u //storing user information into the username Map
		// redirect

		dbUsers[u.UserName].CartMapData = make(map[string]*cartFullData)            //alvin
		dbUsers[u.UserName].CheckoutMapData = make(map[string]*checkoutMapDataFull) //alvin

		http.Redirect(w, req, "/", http.StatusSeeOther) //once logged in, redirect to where you want the user to be redirected to
		return
	}
	showSessions() // for demonstration purposes
	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func login(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther) //if alreadyLoggedIn == true -> returns them to see what they're supposed to see etc.
		return
	}
	var u user

	var me = make(map[string]string) //make map for error
	rx := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	// process form submission
	if req.Method == http.MethodPost {

		un := req.FormValue("username") //the 'name' of the field
		p := req.FormValue("password")  //the 'name' of the field
		// is there a username?

		if !rx.MatchString(un) {
			me["Username1"] = "Username entered is not a valid email address."
			tpl.ExecuteTemplate(w, "login.gohtml", me)
			return
		}

		u, ok := dbUsers[un] //username = email as login - key is user struct - userName, password etc
		if !ok {
			me["Username1"] = "Invaid Username or Username not found"
			tpl.ExecuteTemplate(w, "login.gohtml", me)
			return
		}

		err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		if err != nil {

			me["Password"] = "Invaid Password Entered. Please try again"
			tpl.ExecuteTemplate(w, "login.gohtml", me)

			return
		}

		// create session

		duplicateLoginCheck(u, w, req)

		sID, err := uuid.NewV4()
		//err handling
		if err != nil {
			fmt.Printf("Something went wrong: %s, err")
		}

		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}

		dbUsers[u.UserName].CartMapData = make(map[string]*cartFullData)            //alvin
		dbUsers[u.UserName].CheckoutMapData = make(map[string]*checkoutMapDataFull) //alvin

		c.MaxAge = sessionLength
		http.SetCookie(w, c)
		dbSessions[c.Value] = session{un, time.Now()}
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	showSessions() // for demonstration purposes
	tpl.ExecuteTemplate(w, "login.gohtml", u)
}

func logout(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) { //if you are not logged in, there's nothing you need to do. whatever UUID cookie value, belongs to a non-logged in visitor
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	c, _ := req.Cookie("session")
	// delete the session
	delete(dbSessions, c.Value)
	// remove the cookie
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	// clean up dbSessions
	if time.Now().Sub(dbSessionsCleaned) > (time.Second * 30) {
		go cleanSessions()
	}

	http.Redirect(w, req, "/", http.StatusSeeOther) //goes back to home page after logging out
}

func rs2(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "assets/rs2.png")
}

func searchresult(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req) //getUser function call

	var localSlice = []searchResultFormat{}
	localResult := dbUsers[u.UserName].SearchLogs.SearchResults

	for _, v := range localResult { //range through all available data in the slice

		valuepair := foodNameAddresstoID[v]
		localSlice = append(localSlice, searchResultFormat{v, valuepair}) //everytime a new item is added into cart, this gets appended
	}

	if req.Method == http.MethodPost {
		// sr := req.FormValue("s2")
		productid := req.FormValue("pid") //pid is also known as the productID
		quantity1 := req.FormValue("quantity")

		u.CartList = []cartData{{productid, quantity1}} //as soon as you click the add button, whatever data this has, gets reset.

		for _, v := range u.CartList { //for every element you have in your cart, this will be generated

			foodName := foodNameAddresstoname[v.PID]
			quantity := v.Quantity

			value := MyFoodListMap[foodName] //retrieves all food informatin. LIKE FULL DATA.
			unitPrice := value.Price
			quantityInt, _ := strconv.Atoi(quantity)

			totalCost := float64(quantityInt) * unitPrice
			dbUsers[u.UserName].CartMapData[productid] = &cartFullData{productid, foodName, quantity1, unitPrice, totalCost, u.Role}
		}

		// for k, v := range dbUsers[u.UserName].CartMapData {
		// 	fmt.Println("this is key and value pair", k, v)
		// }
		http.Redirect(w, req, "/yourcart", http.StatusSeeOther)
	}

	showSessions() // for demonstration purposes

	tpl.ExecuteTemplate(w, "searchresult.gohtml", localSlice)

}

func yourcart(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req) //getUser function call

	if u == nil {
		//no username present!
		http.Redirect(w, req, "/login_redirect", http.StatusSeeOther)
		return
	}

	u.cartTransID = []string{} //clears cart slice iD

	showSessions() // for demonstration purposes

	if req.Method == http.MethodPost {

		rb := req.FormValue("updatecart")

		if rb == "updatecart" {

			rp := req.FormValue("pid")
			rq := req.FormValue("quantity")

			if rq == "0" { //if quantity is 0, delete key (pid) from map - CartMapData
				delete(dbUsers[u.UserName].CartMapData, rp)
				http.Redirect(w, req, "/yourcart", http.StatusSeeOther)
				return
			}

			irq, _ := strconv.Atoi(rq) //conversion of string rq to integer
			firq := float64(irq)       //conversion of int rq to float64

			dbUsers[u.UserName].CartMapData[rp].Quantity = rq
			dbUsers[u.UserName].CartMapData[rp].TotalCost = dbUsers[u.UserName].CartMapData[rp].UnitPrice * firq
			http.Redirect(w, req, "/yourcart", http.StatusSeeOther)

			return
			//change map cart data by assigning.
		}

		var pi int = 0

		pi, _ = strconv.Atoi(req.FormValue("priorityindex"))

		userCartData := dbUsers[u.UserName].CartMapData //userCartData is a map of values.
		generatedSysQueueID, err := generateSysQueueID()

		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		for _, v := range userCartData {

			generatedID, err := generateTransactionID() //pass user details so we know how to reach dbUsers dataqbase
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			u.cartTransID = append(u.cartTransID, generatedID)

			dbtransIDSystem[generatedID] = &transFullData{u.UserName, v.FoodName, v.Quantity, v.UnitPrice, v.TotalCost, time.Now()}
		}
		dbUsers[u.UserName].CheckoutMapData[generatedSysQueueID] = &checkoutMapDataFull{u.UserName, u.Role, time.Now(), u.cartTransID, generatedSysQueueID, pi}

		dbsysQueueSystem[generatedSysQueueID] = &sysQueueMapDataFull{u.UserName, u.Role, time.Now(), u.cartTransID, generatedSysQueueID, pi, ""} //adds new item into global variable

		SysQueue.Enqueue(generatedSysQueueID, pi) //adds queue to system

		for k, _ := range dbUsers[u.UserName].CartMapData { //delete key instead
			delete(dbUsers[u.UserName].CartMapData, k)
		}
		http.Redirect(w, req, "/checkout_processing", http.StatusSeeOther)
	}

	searchResult2 = []searchResultFormat{} //clear searhResults2

	if u.Role == "Customer Service Officer" || u.Role == "superuser#1" {

		tpl.ExecuteTemplate(w, "yourcart_admin.gohtml", dbUsers[u.UserName].CartMapData)
	} else {
		tpl.ExecuteTemplate(w, "yourcart.gohtml", dbUsers[u.UserName].CartMapData)

	}

}

func viewall(w http.ResponseWriter, req *http.Request) {

	// showSessions()                         // for demonstration purpose
	// searchResult2 = []searchResultFormat{} //clear searhResults2

	// SysQueue.PrintAllNodes()
	display := SysQueue.parsesystemqueuedata() //should get a slice of structs

	tpl.ExecuteTemplate(w, "viewall.gohtml", display)

}

// func productdisplay(w http.ResponseWriter, req *http.Request) {

// 	showSessions() // for demonstration purposes
// 	tpl.ExecuteTemplate(w, "productdisplay.gohtml", selectedProduct)

// }

func checkout_processing(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req) //getUser function call

	showSessions() // for demonstration purposes

	tpl.ExecuteTemplate(w, "checkout_processing.gohtml", dbUsers[u.UserName])

}

func checkout(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req) //getUser function call

	showSessions() // for demonstration purposes

	m := parseDataforCheckout(u)

	if req.Method == http.MethodPost {
		//no username present!
		r := req.FormValue("homebutton")
		if r == "home" {
			http.Redirect(w, req, "/", http.StatusSeeOther)
		}

	}
	tpl.ExecuteTemplate(w, "checkout.gohtml", m)
}

func parseDataforCheckout(u *user) map[string]*checkoutDisplay {

	newMap := make(map[string]*checkoutDisplay)

	var newKey string

	for k, v := range dbUsers[u.UserName].CheckoutMapData {

		for _, m := range v.TransID { //v.TransID is a slice

			newKey = k + "-" + m
			newMap[newKey] = &checkoutDisplay{
				dbtransIDSystem[m].FoodName,
				dbtransIDSystem[m].Quantity,
				dbtransIDSystem[m].UnitPrice,
				dbtransIDSystem[m].TotalCost,
			}
		}
	}
	return newMap
}

func allsystemorders(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req) //getUser function call

	showSessions() // for demonstration purposes

	if u == nil {
		//no username present!
		http.Redirect(w, req, "/allthefoodisgone", http.StatusSeeOther)
		return
	}

	if u.Role == "Customer Service Officer" || u.Role == "superuser#1" {

		tpl.ExecuteTemplate(w, "allsystemorders.gohtml", dbsysQueueSystem)

	} else {

		tpl.ExecuteTemplate(w, "allthefoodisgone.gohtml", nil)
	}
}

func alltransactions(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req) //getUser function call

	showSessions() // for demonstration purposes

	if u == nil {
		//no username present!
		http.Redirect(w, req, "/allthefoodisgone", http.StatusSeeOther)
		return
	}

	// searchResult2 = []searchResultFormat{} //clear searhResults2

	// tpl.ExecuteTemplate(w, "alltransactions.gohtml", dbtransIDSystem)

	if u.Role == "Customer Service Officer" || u.Role == "superuser#1" {

		tpl.ExecuteTemplate(w, "alltransactions.gohtml", dbtransIDSystem)

	} else {

		tpl.ExecuteTemplate(w, "allthefoodisgone.gohtml", nil)
	}

}

func login_redirect(w http.ResponseWriter, req *http.Request) {

	// u := getUser(w, req) //getUser function call

	showSessions() // for demonstration purposes

	// searchResult2 = []searchResultFormat{} //clear searhResults2

	tpl.ExecuteTemplate(w, "login_redirect.gohtml", nil)

}

func viewqueue(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req) //getUser function call
	showSessions()       // for demonstration purposes

	if u == nil {
		//no username present!
		http.Redirect(w, req, "/allthefoodisgone", http.StatusSeeOther)
		return
	}

	display := SysQueue.parsesystemqueuedata()
	// tpl.ExecuteTemplate(w, "viewqueue.gohtml", display)

	if u.Role == "Customer Service Officer" || u.Role == "superuser#1" {

		tpl.ExecuteTemplate(w, "viewqueue.gohtml", display)
	} else {
		tpl.ExecuteTemplate(w, "allthefoodisgone.gohtml", nil)
	}

}

func clearcart(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req) //getUser function call

	showSessions() // for demonstration purposes

	if u == nil {
		//no username present!
		http.Redirect(w, req, "/login_redirect", http.StatusSeeOther)
		return
	}

	for k, _ := range dbUsers[u.UserName].CartMapData { //delete key instead
		delete(dbUsers[u.UserName].CartMapData, k)
	}
	// searchResult2 = []searchResultFormat{} //clear searhResults2

	tpl.ExecuteTemplate(w, "yourcart.gohtml", nil)

}

func validateInputs(un string, p string, f string, l string, me map[string]string) (bool, map[string]string) {

	rx := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if len(un) == 0 {
		me["Username"] = "Username is not valid. Please Enter again"
	} else if !rx.MatchString(un) {
		me["Username"] = "Username is not valid email address. Please Enter again"
	}
	if len(p) == 0 {
		me["Password"] = "Password is not valid. Please Enter again"
	}
	if len(f) == 0 {
		me["FirstName"] = "First Name is not valid. Please Enter again"
	}
	if len(l) == 0 {
		me["LastName"] = "Last Name is not valid. Please Enter again"
	}
	if len(un) != 0 && len(p) != 0 && len(un) != 0 && len(l) != 0 && rx.MatchString(un) {
		return true, me
	}
	return false, me
}

func dispatchdriver(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req) //getUser function call
	if u == nil {
		http.Redirect(w, req, "/allthefoodisgone", http.StatusSeeOther)
		return
	}

	showSessions() // for demonstration purposes

	display := SysQueue.parsesystemqueuedata()

	if req.Method == http.MethodPost {
		rb := req.FormValue("updatedriver")
		if rb == "updatedriver" {

			rsq := req.FormValue("systemqueuenumber") //request system queue number
			rdn := req.FormValue("drivername")        //request assigned driver name

			dbsysQueueSystem[rsq].DriverName = rdn

			display = SysQueue.parsesystemqueuedata() //call function again to get updated slice (which is reierated for map. otherwise driver values WILL still display old ones)
			tpl.ExecuteTemplate(w, "dispatchdriver.gohtml", display)

			return
			//change map cart data by assigning.
		}

	}

	if u.Role == "Dispatch Supervisor" || u.Role == "superuser#1" {

		tpl.ExecuteTemplate(w, "dispatchdriver.gohtml", display)
	} else {
		tpl.ExecuteTemplate(w, "allthefoodisgone.gohtml", nil)

	}

}

func allthefoodisgone(w http.ResponseWriter, req *http.Request) {

	showSessions() // for demonstration purposes
	// searchResult2 = []searchResultFormat{} //clear searhResults2

	tpl.ExecuteTemplate(w, "allthefoodisgone.gohtml", nil)

}

func dispatchqueue(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req) //getUser function call

	showSessions() // for demonstration purposes

	if u == nil {
		http.Redirect(w, req, "/allthefoodisgone", http.StatusSeeOther)
		return
	}

	display := SysQueue.parsesystemqueuedata()

	if req.Method == http.MethodPost {

		rb := req.FormValue("dispatchqueue")
		if rb == "dispatchqueue" {

			// rsq := req.FormValue("systemqueuenumber") //request system queue number
			// rdn := req.FormValue("drivername")        //request assigned driver name

			// SysQueue.parsesystemqueuedata() //call function again to get updated slice (which is reierated for map. otherwise driver values WILL still display old ones)
			SysQueue.Dequeue()
			display1 := SysQueue.parsesystemqueuedata()
			fmt.Println(SysQueue)

			tpl.ExecuteTemplate(w, "dispatchqueue.gohtml", display1)

			return
			//change map cart data by assigning.
		}

	}

	// searchResult2 = []searchResultFormat{} //clear searhResults2

	tpl.ExecuteTemplate(w, "dispatchqueue.gohtml", display)

}
