package torronto

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func main() {
	// connect to socket
	// service := get service port, maybe it's provided when
	// the program is run

	// tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	// checkError(err)
	// listener, err := net.ListenTCP("tcp", tcpAddr)
	// checkError(err)

	// read for input from stdin and do the appropriate action
	// e.g. query for file, insert a file
	// should probably be run as a goroutine...
	// go listenForInput()?

	// listen for messages?
	// for {
	//   conn, err := listener.Accept()
	//   if err != nil {
	//     continue
	//   }

	//   goroutine for handling connection
	//   go peerhandler(conn)?
	// }
	//
	// instantiate peer and join network

}

// Listen for user input to stdin
// func listenForInput() {

//}

// handle a peer connection
// func peerhandler(conn net.Conn) {
//  check message:

//  if (new peer) {
  //   create peer object
  //   if (first) {
  //     peer.newConnection(newPeer)
  //   }
  //   peer.Peers.addPeer(newPeer)
  //   passAlongMessage(message)
  // }

if (leaving) {
  peer.updateConnect(oldPeer)
  peer.Peers.removePeer(oldPeer)
  passAlongMessage(message)
}
  if leaving network, send out message to everyone

// }

// Possibly use a check error function to handle the many different
// communication errors that can happen
// func checkError(err error) {
//   if err != nil {
//     // error handling
//   }
// }
