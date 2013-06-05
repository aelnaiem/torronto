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
		
//		testSingleJoin();
//		testSingleLeave();
//		
//		testNetworkJoin();
//		testNetworkLeave();
		
		try {
			testJoinLeaveMultiple();
		} catch (InterruptedException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}

	public static void testSingleJoin() {
		System.out.println("Peer 1 joins network");
		String joinMsg = createJoinMessage().toString();
		
		res = Message(peerOne, joinMsg);
		System.out.println(res);
		System.out.println();
	}
	
	public static void testSingleLeave() {
		System.out.println("Peer 1 leaves network");
		String leaveMsg = createLeaveMessage().toString();
		
		res = Message(peerOne, leaveMsg);
		System.out.println(res);
		System.out.println();
	}
	
	public static void testNetworkJoin() {
		System.out.println("Peer 1 and 2 join network");
		String joinMsg = createJoinMessage().toString();
		
		res = Message(peerOne, joinMsg);
		System.out.println(res);
		
		res = Message(peerTwo, joinMsg);
		System.out.println(res);
		System.out.println();
	}
	
	public static void testNetworkLeave() {
		System.out.println("Peer 1 and 2 leave network");
		String leaveMsg = createLeaveMessage().toString();
		
		res = Message(peerOne, leaveMsg);
		System.out.println(res);
		
		res = Message(peerOne, leaveMsg);
		System.out.println(res);
		System.out.println();
	}
	
	public static void testJoinLeaveMultiple() throws InterruptedException {	
		System.out.println("Peer 1 and 2 leave network multiple times");
		testNetworkJoin();
		Thread.sleep(1000);		
		testNetworkLeave();
		Thread.sleep(1000);
		testNetworkJoin();
		Thread.sleep(1000);		
		testNetworkLeave();
		Thread.sleep(1000);
		testSingleJoin();
		Thread.sleep(1000);		
		testSingleLeave();
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

	public static JsonObject createInsertMessage(String f){
		JsonObject insertMessage = new JsonObject();
		insertMessage.addProperty("HostName", hostName);
		insertMessage.addProperty("PortNumber", portNumber);
		insertMessage.addProperty("Action", 2);

		JsonObject file = new JsonObject();
		file.addProperty("fileName", f);
		JsonArray files = new JsonArray();
		files.add(file);
		insertMessage.add("files", files);

		return insertMessage;
	}

	public static JsonObject createQueryMessage(){
		JsonObject queryMessage = new JsonObject();
		queryMessage.addProperty("HostName", hostName);
		queryMessage.addProperty("PortNumber", portNumber);
		queryMessage.addProperty("Action", 3);
		return queryMessage;
	}

}

