import java.io.*;
import java.net.*;
import java.nio.charset.StandardCharsets;

public class EasyTCPClient {
    static String SERVER_NAME = "localhost";
    static int SERVER_PORT = 14094;

    public static void main(String[] args) {
        Socket serverConnection = makeConnection();
        if (serverConnection == null) {
            printError("Failed to init client.");
            return;
        }

        Runtime.getRuntime().addShutdownHook(new Thread(() -> closeConnection(serverConnection)));

        try {
            serverConnection.setSoTimeout(5000);
            OutputStream requestStream = serverConnection.getOutputStream();
            InputStream responseStream = serverConnection.getInputStream();

            while (true) {
                printMenu();

                int cmd = readCommand();
                String text = "";
                if (cmd == 0) {
                    continue;
                } else if (cmd == 5) {
                    break;
                } else if (cmd == 1) {
                    System.out.print("Input lowercase sentence: ");
                    BufferedReader reader = new BufferedReader(new InputStreamReader(System.in));
                    text = reader.readLine();

                    if (text.length() >= 1024) {
                        printError("Text too long.");
                        continue;
                    }
                }

                long requestTime = System.nanoTime();
                byte[] requestBuffer = String.format("%d%s", cmd, text).getBytes(StandardCharsets.UTF_8);
                requestStream.write(requestBuffer);
                requestStream.flush();

                byte[] responseBuffer = new byte[1024];
                int responseSize = responseStream.read(responseBuffer);
                long responseTime = System.nanoTime();

                System.out.printf("\nReply from server: %s\n", new String(responseBuffer, 0, responseSize).trim());
                System.out.printf("RTT = %.3fms\n", (responseTime - requestTime) / 1000000f);
            }

            requestStream.close();
            responseStream.close();
        } catch (SocketTimeoutException e) {
            printError("Server Timeout.");
        } catch (IOException e) {
            printError(e.toString());
        }
    }

    private static Socket makeConnection() {
        try {
            return new Socket(SERVER_NAME, SERVER_PORT);
        } catch (IOException e) {
            return null;
        }
    }

    private static void closeConnection(Socket conn) {
        System.out.println("\rClosing Client Program...\nBye bye~");
        try {
            OutputStream requestStream = conn.getOutputStream();
            requestStream.write("5".getBytes(StandardCharsets.UTF_8));
            conn.close();
        } catch (IOException e) {}
    }

    private static void printMenu() {
        System.out.println();
        System.out.println("< Select Menu. >");
        System.out.println("1) Convert Text to UPPER-case Letters");
        System.out.println("2) Get Server Uptime");
        System.out.println("3) Get Client IP / Port");
        System.out.println("4) Get Count of Requests Server Got");
        System.out.println("5) Exit Client");
    }

    private static int readCommand() {
        int cmd;
        String input;

        try {
            System.out.print("Input Command: ");
            BufferedReader reader = new BufferedReader(new InputStreamReader(System.in));
            input = reader.readLine();
            cmd = Integer.parseInt(input);
        } catch (NumberFormatException e) {
            printError("Invalid Command.");
            return 0;
        } catch (IOException e) {
            printError("System Error.");
            return 0;
        }

        if (cmd < 1 || cmd > 5) {
            printError("Invalid Command.");
            return 0;
        }

        return cmd;
    }

    private static void printError(String msg) {
        System.err.printf("Error: %s\n", msg);
    }
}