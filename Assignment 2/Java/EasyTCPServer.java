import java.io.*;
import java.net.*;
import java.nio.charset.StandardCharsets;
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
                InputStream requestStream = serverConnection.getInputStream();
                OutputStream responseStream = serverConnection.getOutputStream();

                while (true) {
                    byte[] requestBuffer = new byte[1024];
                    int requestSize = requestStream.read(requestBuffer);

                    if(requestSize >= 0){
                        String requestData = new String(requestBuffer, 0, requestSize).trim();
                        String requestIP = serverConnection.getInetAddress().toString();
                        int requestPort = serverConnection.getPort();

                        if (requestData.equals("5") || requestData.isEmpty()) {
                            serverConnection.close();
                            break;
                        }

                        System.out.printf("TCP Connection Request from %s:%d\n", requestIP, requestPort);
                        System.out.printf("Command %c\n", requestData.charAt(0));
                        String responseData = getResponse(requestData.charAt(0),
                                requestData.substring(1),
                                requestIP,
                                requestPort);

                        byte[] responseBuffer = responseData.getBytes(StandardCharsets.UTF_8);
                        responseStream.write(responseBuffer);
                        responseStream.flush();

                        serverResponseCnt++;
                    }else{
                        serverConnection.close();
                        break;
                    }
                }

                requestStream.close();
                responseStream.close();
            } catch (IOException e) {
                printError(e.toString());
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
                return data.toUpperCase().trim();
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