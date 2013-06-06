import java.io.DataInputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.net.Socket;
import java.net.ServerSocket;
import java.nio.ByteBuffer;
import java.nio.charset.Charset;

import com.google.gson.JsonArray;
import com.google.gson.JsonObject;

public class Test {
	public static String hostName = "127.0.0.1";
	public static int portNumber = 10000;
	public static ServerSocket server;
	
	public static String res;
	public static Peer peerOne, peerTwo, peerThree, peerFour, peerFive;
	public static void main(String[] args) {
		try {
			server = new ServerSocket(portNumber);
		} catch (Exception e){
			return;
		}
		peerOne = new Peer("127.0.0.1", 10001);
		peerTwo = new Peer("127.0.0.1", 10002);
		peerThree = new Peer("127.0.0.1", 10003);
		peerFour = new Peer("127.0.0.1", 10004);
		peerFive = new Peer("127.0.0.1", 10005);
		
		Peer[] peers = {peerOne, peerTwo, peerThree, peerFour, peerFive};
		Peer[] peersOne = {peerOne, peerThree, peerFour};
		Peer[] peersTwo = {peerTwo, peerThree, peerFive};
		
		testSingleJoin(peerOne);
		testSingleLeave(peerOne);
		
		testNetworkJoin(peers);
		testNetworkLeave(peers);

		testJoinLeaveMultiple(peerOne, peersOne, peersTwo);
		
		testNetworkDoubleJoin(peerOne);
		testNetworkDoubleLeave(peerOne);
		
		testSingleJoin(peerOne);
		testInsert("oi", peers);
		testInsert("mia", peers);
		
		testQuery();
	}

	public static void testSingleJoin(Peer peer) {
		System.out.println("Peer " + peer.port+ " joins the network");
		String joinMsg = createJoinMessage().toString();
		
		res = Message(peer, joinMsg);
		System.out.println(res);
		System.out.println();
	}
	
	public static void testSingleLeave(Peer peer) {
		System.out.println("Peer " + peer.port + " leaves the network");
		String leaveMsg = createLeaveMessage().toString();
		
		res = Message(peer, leaveMsg);
		System.out.println(res);
		System.out.println();
	}
	
	public static void testNetworkJoin(Peer[] peers) {
		System.out.println("All peers join the network");
		for (int i = 0; i < peers.length; i ++) {
			testSingleJoin(peers[i]);
		}
	}
	
	public static void testNetworkLeave(Peer[] peers) {
		System.out.println("All peers join the network");
		for (int i = 0; i < peers.length; i ++) {
			testSingleJoin(peers[i]);
		}
	}
	
	public static void testJoinLeaveMultiple(Peer peer, Peer[] peersOne, Peer[] peersTwo) {	
		System.out.println("Peer 1 and 2 join and leave the network multiple times");
		testNetworkJoin(peersOne);
		testNetworkJoin(peersTwo);
		testNetworkLeave(peersOne);
		testSingleJoin(peer);
		testNetworkLeave(peersTwo);
		testSingleLeave(peer);
	}
	
	public static void testNetworkDoubleJoin(Peer peer) {
		System.out.println("Peer 1 and 2 join the network twice");
		testSingleJoin(peer);
		testSingleJoin(peer);		
	}
	
	public static void testNetworkDoubleLeave(Peer peer) {
		System.out.println("Peer 1 and 2 leave the network twice");
		testSingleLeave(peer);
		testSingleLeave(peer);		
	}
	
	public static void testInsert(String file, Peer[] peers) {
		System.out.println("Inserting a file in peer 1 " + file);
		String insertMsg = createInsertMessage(file).toString();
		testNetworkJoin(peers);		
		
		res = Message(peerOne, insertMsg);
		System.out.println(res);
		System.out.println();

		System.out.println("Inserting the same file in peer 1");
		res = Message(peerOne, insertMsg);
		System.out.println(res);
		System.out.println();
	}
	
	public static void testQuery() {
		System.out.println("Query Peer 1 and Peer 2 status");
		String leaveMsg = createQueryMessage().toString();
		
		res = Message(peerOne, leaveMsg);
		System.out.println(res);
		System.out.println();
	}
	
	public static String Message(Peer peer, String msgString) {
		try {
			Socket socket = new Socket(peer.host, peer.port);
			DataOutputStream os = new DataOutputStream(socket.getOutputStream());
			byte[] msgByte = (msgString.getBytes());

			os.write(msgByte);
			os.close();

			byte[] buf = new byte[100];
			Socket link = server.accept();

			DataInputStream response = new DataInputStream(link.getInputStream());
			try{
				response.readFully(buf);
			} catch (Exception e){
			}

			link.close();
			socket.close();

			return Charset.forName("UTF-8").decode(ByteBuffer.wrap(buf)).toString();
		} catch (IOException e) {
			e.printStackTrace();
			return "";
		}
	}

	public static JsonObject createJoinMessage(){
		JsonObject joinMessage = new JsonObject();
		joinMessage.addProperty("HostName", hostName);
		joinMessage.addProperty("PortNumber", portNumber);
		joinMessage.addProperty("Action", 0);
		return joinMessage;
	}

	public static JsonObject createLeaveMessage(){
		JsonObject leaveMessage = new JsonObject();
		leaveMessage.addProperty("HostName", hostName);
		leaveMessage.addProperty("PortNumber", portNumber);
		leaveMessage.addProperty("Action", 1);
		return leaveMessage;
	}

	public static JsonObject createQueryMessage(){
		JsonObject queryMessage = new JsonObject();
		queryMessage.addProperty("HostName", hostName);
		queryMessage.addProperty("PortNumber", portNumber);
		queryMessage.addProperty("Action", 2);
		return queryMessage;
	}
	
	public static JsonObject createInsertMessage(String f){
		JsonObject insertMessage = new JsonObject();
		insertMessage.addProperty("HostName", hostName);
		insertMessage.addProperty("PortNumber", portNumber);
		insertMessage.addProperty("Action", 3);

		JsonObject file = new JsonObject();
		file.addProperty("fileName", f);
		JsonArray files = new JsonArray();
		files.add(file);
		insertMessage.add("files", files);

		return insertMessage;
	}

}

