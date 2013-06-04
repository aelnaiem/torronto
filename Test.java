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

	public static void main(String[] args) {
		String res;
		try {
			server = new ServerSocket(portNumber);
		} catch (Exception e){
			return;
		}
		Peer peerOne = new Peer("127.0.0.1", 10001);
		Peer peerTwo = new Peer("127.0.0.1", 10002);

		res = Join(peerOne);
		System.out.println(res);

		res = Join(peerTwo);
		System.out.println(res);
	}

	public static String Join(Peer peer) {
		try {
			Socket socket = new Socket(peer.host, peer.port);
			DataOutputStream os = new DataOutputStream(socket.getOutputStream());

			String msgString = createJoinMessage().toString();
			byte[] msgByte = (msgString.getBytes());

			os.write(msgByte);
			os.close();

			byte[] buf = new byte[100];
			Socket link = server.accept();
			System.out.println("hello");
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

