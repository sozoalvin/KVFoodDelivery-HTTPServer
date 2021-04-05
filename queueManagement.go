package main

// import (
// 	"errors"
// 	"fmt"
// )

// type QNode struct {
// 	priority int
// 	item     []KVorder
// 	next     *QNode
// }

// type Queue struct {
// 	front *QNode
// 	back  *QNode
// 	size  int
// }

// func InitSysQueue() *Queue {
// 	result := &Queue{nil, nil, 0}
// 	return result
// }

// func (p *Queue) Enqueue(k []KVorder, prty int) error {
// 	newNode := &QNode{
// 		priority: prty,
// 		item:     k,
// 		next:     nil,
// 	}
// 	if p.front == nil {
// 		p.front = newNode
// 	} else {
// 		if p.front.priority < prty {
// 			newNode.next = p.front
// 			p.front = newNode
// 		} else {
// 			currentNode := p.front
// 			for currentNode.next != nil && currentNode.next.priority >= prty {
// 				currentNode = currentNode.next
// 			}
// 			newNode.next = currentNode.next
// 			currentNode.next = newNode
// 		}
// 	}
// 	p.size++
// 	return nil
// }

// func (p *Queue) PrepreEnqueue(ch chan string) {

// 	p.PreEnqueue(FirstQueueValue, 3)
// 	p.PreEnqueue(ThirdQueueValue, 7)
// 	p.PreEnqueue(SecondQueueValue, 5)

// 	ch <- "Queues Loaded into System"

// }

// func (p *Queue) PreEnqueue(k []KVorder, prty int) error {
// 	newNode := &QNode{
// 		priority: prty,
// 		item:     k,
// 		next:     nil,
// 	}
// 	if p.front == nil {
// 		p.front = newNode
// 	} else {
// 		if p.front.priority < prty {
// 			newNode.next = p.front
// 			p.front = newNode
// 		} else {
// 			currentNode := p.front
// 			for currentNode.next != nil && currentNode.next.priority >= prty {
// 				currentNode = currentNode.next
// 			}
// 			newNode.next = currentNode.next
// 			currentNode.next = newNode
// 		}
// 	}
// 	p.size++
// 	return nil
// }

// func (p *Queue) Dequeue() error {

// 	if p.front == nil {
// 		return errors.New("empty queue")
// 	}
// 	if p.size == 1 {
// 		p.front = nil
// 		p.back = nil
// 	} else {
// 		p.front = p.front.next
// 	}
// 	p.size--
// 	return nil
// }

// func (p *Queue) PrintAllNodes() error {
// 	currentNode := p.front
// 	if currentNode == nil {
// 		fmt.Println("Queue is empty")
// 		return nil
// 	}

// 	fmt.Printf("\nAuthorized Priority Queue Number: %d tagged to queue.\n\n", currentNode.priority)

// 	for _, v := range currentNode.item {

// 		for i, v := range v.transID {
// 			fmt.Printf("Associated Merchant Transaction id %d:\t\t%s\n", i+1, v)
// 		}
// 		fmt.Printf("Order System Queue ID:\t\t\t\t%s\n", v.systemQueueID)
// 		fmt.Printf("Order tagged to username:\t\t\t%s\n", v.username)
// 	}

// 	fmt.Println("\n==========================================================")

// 	for currentNode.next != nil {
// 		currentNode = currentNode.next
// 		fmt.Printf("\nAuthorized Priority Queue Number: %d tagged to queue.\n\n", currentNode.priority)
// 		for _, v := range currentNode.item {

// 			for i, v := range v.transID {
// 				fmt.Printf("%d.Associated Merchant Transaction id: \t\t%s\n", i+1, v)
// 			}
// 			fmt.Printf("Order System Queue ID:\t\t\t\t%s\n", v.systemQueueID)
// 			fmt.Printf("Order tagged to username:\t\t\t%s\n", v.username)
// 		}

// 		fmt.Println("\n==========================================================")

// 	}
// 	return nil

// }
