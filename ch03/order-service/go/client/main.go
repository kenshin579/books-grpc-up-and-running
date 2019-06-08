package main

import (
	"context"
	"io"
	"log"
	"time"

	wrapper "github.com/golang/protobuf/ptypes/wrappers"
	pb "github.com/grpc-up-and-running/samples/ch03/order-service/go/order_service"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Setting up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewOrderManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Add Order
	order1 := pb.Order{Id: "101", Items:[]string{"iPhone XS", "Mac Book Pro"}, Destination:"San Jose, CA", Price:2300.00}
	res, _ := c.AddOrder(ctx, &order1)
	log.Print("AddOrder Response -> ", res.Value)



	// Get Order
	retrievedOrder , err := c.GetOrder(ctx, &wrapper.StringValue{Value: "106"})
	log.Print("GetOrder Response -> : ", retrievedOrder)


	// Search Order
	searchStream, _ := c.SearchOrders(ctx, &wrapper.StringValue{Value: "Google"})
	for {
		searchOrder, err := searchStream.Recv()
		if err == io.EOF {
			break
		}
		log.Print("Search Result : ", searchOrder)
	}


	// Update Orders

	updOrder1 := pb.Order{Id: "102", Items:[]string{"Google Pixel 3A", "Google Pixel Book"}, Destination:"Mountain View, CA", Price:1100.00}
	updOrder2 := pb.Order{Id: "103", Items:[]string{"Apple Watch S4", "Mac Book Pro", "iPad Pro"}, Destination:"San Jose, CA", Price:2800.00}
	updOrder3 := pb.Order{Id: "104", Items:[]string{"Google Home Mini", "Google Nest Hub", "iPad Mini"}, Destination:"Mountain View, CA", Price:2200.00}

	updateStream, _ := c.UpdateOrders(ctx)
	_ = updateStream.Send(&updOrder1)
	_ = updateStream.Send(&updOrder2)
	_ = updateStream.Send(&updOrder3)

	updateRes, _ := updateStream.CloseAndRecv()
	log.Printf("Update Orders Res : ", updateRes)


	//// Process Order
	//order4 := pb.Order{Id: "100600", Name: "Order for prc"}
	//streamProc, _ := c.ProcessOrders(ctx)
	//streamProc.Send(&order4)
	//
	//combinedShipment, _ := streamProc.Recv()
	//
	//log.Printf("Combined shipment ", combinedShipment.OrderIDList)

}