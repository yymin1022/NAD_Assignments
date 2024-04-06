import java.io.*;
import java.net.InetSocketAddress;
import java.net.ServerSocket;
import java.net.Socket;

public class EasyTCPServer {
    static int SERVER_PORT = 14094;

    public static void main(String[] args) throws IOException{
        ServerSocket serverSocket = new ServerSocket(SERVER_PORT);
        while(true) {
            System.out.println("Waiting for Connection...");
            Socket socket = serverSocket.accept();

            InetSocketAddress isa = (InetSocketAddress) socket.getRemoteSocketAddress();
            System.out.println("Connection from " + isa.getHostName());

            BufferedReader reader = new BufferedReader(new InputStreamReader(socket.getInputStream()));																	// 데이터를 읽어옴
            String line = reader.readLine();
            PrintWriter writer = new PrintWriter(new OutputStreamWriter(socket.getOutputStream()));
            writer.println(line);
            writer.flush();
        }
    }
}