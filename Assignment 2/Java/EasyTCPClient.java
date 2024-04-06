import java.io.*;
import java.net.Socket;

public class EasyTCPClient {
    static String SERVER_NAME = "localhost";
    static int SERVER_PORT = 14094;

    public static void main(String[] args) throws IOException {
        Socket serverConnection = new Socket(SERVER_NAME, SERVER_PORT);

        String msg= "Hello, World!";
        PrintWriter writer = new PrintWriter(new OutputStreamWriter(serverConnection.getOutputStream()));
        writer.println(msg);
        writer.flush();

        BufferedReader reader = new BufferedReader(new InputStreamReader(serverConnection.getInputStream()));

        writer.close();
        reader.close();
    }
}