package main

import (
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

func getUser(w http.ResponseWriter, req *http.Request) *user {
	// get cookie
	c, err := req.Cookie("session") //requesting to check if the client has a cookie, named session

	if err != nil { //fired if no cookie named, session is present

		sID, err := uuid.NewV4() //if no cookie present, we give it one.
		//err handling
		if err != nil {
			fmt.Printf("Something went wrong: %s, err") //prints if there is error when have error generating UUID
		}

		c = &http.Cookie{
			Name:  "session",    //nametype of cookie
			Value: sID.String(), // Returns canonical string representation of UUID:
		}
	} // if cookie not present; all above codes will run.

	c.MaxAge = sessionLength //tells cookie how short or long lived it is.

	http.SetCookie(w, c) // http.SetCookie is required to 'set'
	// if the user exists already, get user
	var u *user

	//dbSessionMaps keyvalue[c.values] are only added INTO the map during LOGIN AND SIGN UP
	//this function returns FALSE if user is not logged in.

	if s, ok := dbSessions[c.Value]; ok { //c.value == UUID generated.
		fmt.Printf("%+v", s)
		s.lastActivity = time.Now()
		dbSessions[c.Value] = s
		u = dbUsers[s.un]
	}
	return u //if not logged in, this is a nil address
}

func alreadyLoggedIn(w http.ResponseWriter, req *http.Request) bool { //checks if a user is logged or by returning equivalent bool value
	c, err := req.Cookie("session")
	if err != nil {
		return false
	}
	s, ok := dbSessions[c.Value] //if it is an existing sesion if ok, that means the user is already signed up or logged in. because dbSessions only has key values populated when a user logs in / signed in
	if ok {
		s.lastActivity = time.Now() //these updates the s.lastAcitivity time for the instance of the user struct
		dbSessions[c.Value] = s
	}
	_, ok = dbUsers[s.un]
	// refresh session
	c.MaxAge = sessionLength
	http.SetCookie(w, c)
	return ok //ok returns true when the user is already logged in. false when the user is not logged in.
}

func cleanSessions() {
	fmt.Println("(before) db session cleaned") // for demonstration purposes
	showSessions()                             // for demonstration purposes
	for k, v := range dbSessions {
		if time.Now().Sub(v.lastActivity) > (time.Second * 30) {
			delete(dbSessions, k)
		}
	}
	dbSessionsCleaned = time.Now()
	fmt.Println("(after) db session cleaned") // for demonstration purposes
	showSessions()                            // for demonstration purposes
}

func showSessions() {
	fmt.Println("********")
	// fmt.Printf("%s link accessed", s)
	for k, v := range dbSessions {
		fmt.Println(k, v.un)
	}
	fmt.Println("")
}

func duplicateLoginCheck(u *user, w http.ResponseWriter, req *http.Request) {

	// c, _ := req.Cookie("session")
	for k, v := range dbSessions {
		if u.UserName == v.un { //means there is an existing sesion
			delete(dbSessions, k)
			fmt.Println("duplicate session detected - logged out")
			fmt.Println("your other session has been logged out")
			fmt.Println("pls resume on this one")
			break //so we don't waste resources checking everything else
		}

	}

} //end func loop

// c, _ := req.Cookie("session")
// // delete the session
// delete(dbSessions, c.Value)
// // remove the cookie
// c = &http.Cookie{
// 	Name:   "session",
// 	Value:  "",
// 	MaxAge: -1,
// }
// http.SetCookie(w, c)
