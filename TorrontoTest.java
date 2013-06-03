import java.io.BufferedReader;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.Socket;
import com.google.gson.JsonObject;

public class TorrontoTest {
	public static String hostName;
	public static int portNumber;
	public static JsonObject message;
	public static String fileName;

	public static void main(String[ ] args) {
		System.out.println("Hello!");
		hostName = "127.0.0.1";
		portNumber = 10000;
		fileName = "testing123.png"; //include path
		
		joinEmpty();
	}
	
	public static void joinEmpty()
		try {
			Socket socket = new Socket(hostName, portNumber);
			BufferedReader is = new BufferedReader(new InputStreamReader(socket.getInputStream()));
			DataOutputStream os = new DataOutputStream(socket.getOutputStream());

			message = join();
			String msgString = new String();
			msgString = message.toString();

			byte[] msgByte = (msgString.getBytes());

			os.write(msgByte);
			System.out.println(message);
			System.out.println("Join message sent: " + msgByte.toString());
			System.out.println("RECEIVED: " + is.readLine());
			
			is.close()
			os.close()
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
		leaveMessage.addProperty("hostName", hostName);
		leaveMessage.addProperty("portNumber", portNumber);
		leaveMessage.addProperty("action", 1);
		return leaveMessage;
	}
	
	public static JsonObject insert(String f){
		JsonObject insertMessage = new JsonObject();
		insertMessage.addProperty("hostName", hostName);
		insertMessage.addProperty("portNumber", portNumber);
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
		queryMessage.addProperty("hostName", hostName);
		queryMessage.addProperty("portNumber", portNumber);
		queryMessage.addProperty("action", 3);
		return queryMessage;
	}

}