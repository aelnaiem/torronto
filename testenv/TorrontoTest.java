import java.io.BufferedReader;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.Socket;
import java.net.ServerSocket;
import com.google.gson.JsonObject;

public class TorrontoTest {
	public static String hostName = "127.0.0.1";
	public static int portNumber = 10000;

	public static void main(String[] args) {	
		peerOne = new Peer("127.0.0.1", 10001);
		peerTwo = new Peer("127.0.0.1", 10002)
		testJoinEmpty(peerOne);
	}
	
	public static void testJoinEmpty(Peer peer)
		try {
			ServerSocket server = new ServerSocket(portNumber);
			Socket socket = new Socket(peer.host, peer.port);
			DataOutputStream os = new DataOutputStream(socket.getOutputStream());

			String msgString = join().toString();
			byte[] msgByte = (msgString.getBytes());

			os.write(msgByte);
			os.close()
			
			server.accept();
			Scanner input = new Scanner(link.getInputStream()); 
			String message = input.nextLine();  
			System.out.println(message);
			
			is.close()
			socket.close();

		} catch (IOException e) {
			e.printStackTrace();
		}
	}

	public static JsonObject join(){
		JsonObject joinMessage = new JsonObject();
		joinMessage.addProperty("hostName", TorrontoTest.hostName);
		joinMessage.addProperty("portNumber", TorrontoTest.portNumber);
		joinMessage.addProperty("action", 0);
		return joinMessage;
	}
	
	public static JsonObject leave(){
		JsonObject leaveMessage = new JsonObject();
		leaveMessage.addProperty("hostName", TorrontoTest.hostName);
		leaveMessage.addProperty("portNumber", TorrontoTest.portNumber);
		leaveMessage.addProperty("action", 1);
		return leaveMessage;
	}
	
	public static JsonObject insert(String f){
		JsonObject insertMessage = new JsonObject();
		insertMessage.addProperty("hostName", TorrontoTest.hostName);
		insertMessage.addProperty("portNumber", TorrontoTest.portNumber);
		insertMessage.addProperty("action", 2);
		
		JsonObject file = new JsonObject();
		file.addProperty("fileName", f);
		JsonArray files = new JsonArray();		
		files.add(file);
		insertMessage.add("files", files);
	
		return insertMessage;
	}

	public static JsonObject query(){
		JsonObject queryMessage = new JsonObject();
		queryMessage.addProperty("hostName", TorrontoTest.hostName);
		queryMessage.addProperty("portNumber", TorrontoTest.portNumber);
		queryMessage.addProperty("action", 3);
		return queryMessage;
	}

}

public class Peer {
	public static String host;
	public static int port;
	
	public Peer(string hostName, portNumber) {
		host = hostName;
		port = portNumber;
	}
}