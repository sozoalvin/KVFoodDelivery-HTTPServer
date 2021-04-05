package main

import (
	"fmt"
	"strconv"
)

var MyFoodListMap = make(map[string]FoodInfo)
var MyFoodListMap2 = make(map[string]FoodInfo)

var MyCheckoutTranID = make(map[string]Checkout)     //map with transaction ID as key
var MyCheckoutIDUsername = make(map[string]Checkout) //map with username as key

var foodNameAddresstoID = make(map[string]string)
var foodNameAddresstoname = make(map[string]string)

var userCartInfo = make(map[string]string)

func CreateFoodListMap(ch chan string) {

	for _, v := range V {

		keyValue := v.FoodName + " - " + v.MerchantName + " - " + v.DetailedLocation
		MyFoodListMap[keyValue] = FoodInfo{v.FoodName, v.MerchantName, v.DetailedLocation, v.PostalCode, v.Price, v.OpeningPeriods}

	}
	ch <- "Food List Map Data Completed"

}

func CreateFoodListMap2(f *FoodInfo) {

	keyValue := f.FoodName + " - " + f.MerchantName + " - " + f.DetailedLocation
	MyFoodListMap2[keyValue] = FoodInfo{f.FoodName, f.MerchantName, f.DetailedLocation, f.PostalCode, f.Price, f.OpeningPeriods}
	fmt.Println(MyFoodListMap2)
}

func FoodMerchantNameAddressProductID() {
	// func CreateFoodList() []string {

	for _, v := range FoodMerchantNameAddress { //range through entire SLICE to populate map

		value1 := "pID" + strconv.Itoa(pid)

		foodNameAddresstoID[v] = value1

		foodNameAddresstoname[value1] = v
		pid++

	}

	// return FoodMerchantNameAddress
}
