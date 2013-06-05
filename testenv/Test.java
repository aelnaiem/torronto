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
	public static Peer peerOne;
	public static Peer peerTwo;
	public static void main(String[] args) {
		try {
			server = new ServerSocket(portNumber);
		} catch (Exception e){
			return;
		}
		peerOne = new Peer("127.0.0.1", 10001);
		peerTwo = new Peer("127.0.0.1", 10002);
		
		testSingleJoin();
		testSingleLeave();
		
		testNetworkJoin();
		testNetworkLeave();

		testJoinLeaveMultiple();
		
		testNetworkDoubleJoin();
		testNetworkDoubleLeave();
	
		testInsert("oi");
		testInsert("oi");
		
		testQuery();
	}

	public static void testSingleJoin() {
		System.out.println("Peer 1 joins the network");
		String joinMsg = createJoinMessage().toString();
		
		res = Message(peerOne, joinMsg);
		System.out.println(res);
		System.out.println();
	}
	
	public static void testSingleLeave() {
		System.out.println("Peer 1 leaves the network");
		String leaveMsg = createLeaveMessage().toString();
		
		res = Message(peerOne, leaveMsg);
		System.out.println(res);
		System.out.println();
	}
	
	public static void testNetworkJoin() {
		System.out.println("Peer 1 and 2 join the network");
		String joinMsg = createJoinMessage().toString();
		
		res = Message(peerOne, joinMsg);
		System.out.println(res);
		
		res = Message(peerTwo, joinMsg);
		System.out.println(res);
		System.out.println();
	}
	
	public static void testNetworkLeave() {
		System.out.println("Peer 1 and 2 leave the network");
		String leaveMsg = createLeaveMessage().toString();
		
		res = Message(peerOne, leaveMsg);
		System.out.println(res);
		
		res = Message(peerTwo, leaveMsg);
		System.out.println(res);
		System.out.println();
	}
	
	public static void testJoinLeaveMultiple() {	
		System.out.println("Peer 1 and 2 join and leave the network multiple times");
		testNetworkJoin();
		testNetworkLeave();
		testNetworkJoin();	
		testNetworkLeave();
		testSingleJoin();
		testSingleLeave();
	}
	
	public static void testNetworkDoubleJoin() {
		System.out.println("Peer 1 and 2 join the network twice");
		testNetworkJoin();
		testNetworkJoin();		
	}
	
	public static void testNetworkDoubleLeave() {
		System.out.println("Peer 1 and 2 leave the network twice");
		testNetworkLeave();
		testNetworkLeave();		
	}
	
	public static void testInsert(String file) {
		System.out.println("Inserting a file in peer 1");
		String insertMsg = createInsertMessage(file).toString();
		testNetworkJoin();		
		
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

