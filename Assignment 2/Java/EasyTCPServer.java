import java.io.*;
import java.net.*;
import java.time.Duration;
import java.time.LocalDateTime;

public class EasyTCPServer {
    static int SERVER_PORT = 14094;

    static int serverResponseCnt;
    static LocalDateTime serverStartTime = null;

    public static void main(String[] args) {
        ServerSocket serverListener = initServer();
        if (serverListener == null) {
            printError("Failed to init server.");
            return;
        }

        Runtime.getRuntime().addShutdownHook(new Thread(() -> closeConnection(serverListener)));

        serverResponseCnt = 0;
        serverStartTime = LocalDateTime.now();

        while(true){
            try {
                Socket serverConnection = serverListener.accept();

                while (true) {
                    BufferedReader reader = new BufferedReader(new InputStreamReader(serverConnection.getInputStream()));																	// 데이터를 읽어옴
                    String requestData = reader.readLine();
                    String requestIP = serverConnection.getInetAddress().toString();
                    int requestPort = serverConnection.getPort();

                    if (requestData.equals("5")) {
                        serverConnection.close();
                        break;
                    }

                    System.out.printf("TCP Connection Request from %s:%d\n", requestIP, requestPort);
                    String responseData = getResponse(requestData.charAt(0),
                            requestData.substring(1),
                            requestIP,
                            requestPort);

                    PrintWriter writer = new PrintWriter(new OutputStreamWriter(serverConnection.getOutputStream()));
                    writer.println(responseData);
                    writer.flush();
                    serverResponseCnt++;
                }
            } catch (IOException e) {
                printError("IO Error.");
                break;
            }
        }
    }

    private static ServerSocket initServer() {
        try {
            System.out.printf("Server is ready to receive on port %s\n", SERVER_PORT);
            return new ServerSocket(SERVER_PORT);
        } catch (IOException e) {
            return null;
        }
    }

    private static void closeConnection(ServerSocket conn) {
        System.out.println("\rClosing Server Program...\nBye bye~");
        try {
            conn.close();
        } catch (IOException _) {}
    }

    private static String getResponse(char cmd, String data, String addr, int port) {
        switch (cmd) {
            case '1':
                return data.toUpperCase();
            case '2':
                LocalDateTime curTime = LocalDateTime.now();
                Duration upTime = Duration.between(serverStartTime, curTime);

                long upTimeSeconds = upTime.getSeconds();
                return String.format("run time = %02d:%02d:%02d",
                        upTimeSeconds / 3600,
                        (upTimeSeconds % 3600) / 60,
                        (upTimeSeconds % 3600) % 60);

            case '3':
                return String.format("client IP = %s, port = %d", addr, port);
            case '4':
                return String.format("requests served = %d", serverResponseCnt);
        }
        return "";
    }

    private static void printError(String msg) {
        System.err.printf("Error: %s\n", msg);
    }
}