import java.io.*;
import java.net.*;
import java.nio.file.*;
import java.util.concurrent.*;

public class SplitFileServer {
    private static int serverPort;
    private static String filenameSuffix;

    public static void main(String[] args) {
        if (args.length != 1) {
            exitError("Usage: java SplitFileServer <port>");
        }

        try {
            serverPort = Integer.parseInt(args[0]);
        } catch (NumberFormatException e) {
            exitError("Invalid Argument");
        }

        if (serverPort / 10000 == 4) {
            filenameSuffix = "-part1";
        } else {
            filenameSuffix = "-part2";
        }

        ServerSocket serverSocket = initServer();

        Runtime.getRuntime().addShutdownHook(new Thread(() -> closeServer(serverSocket)));

        ExecutorService executorService = Executors.newCachedThreadPool();
        while (true) {
            try {
                Socket socket = serverSocket.accept();
                executorService.execute(() -> handleConnection(socket));
            } catch (IOException e) {
                e.printStackTrace();
            }
        }
    }

    private static void closeServer(ServerSocket serverSocket) {
        System.out.println("\rClosing Server Program...\nBye bye~");
        try {
            if (serverSocket != null && !serverSocket.isClosed()) {
                serverSocket.close();
            }
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    private static ServerSocket initServer() {
        try {
            ServerSocket serverSocket = new ServerSocket(serverPort);
            System.out.println("Server listening on port " + serverPort);
            return serverSocket;
        } catch (IOException e) {
            exitError(e.getMessage());
            return null;
        }
    }

    private static void handleConnection(Socket socket) {
        try (BufferedReader reader = new BufferedReader(new InputStreamReader(socket.getInputStream()));
             PrintWriter writer = new PrintWriter(socket.getOutputStream(), true)) {

            char[] readBuffer = new char[1024];
            int readLength = reader.read(readBuffer);
            if (readLength == -1) {
                return;
            }

            String readData = new String(readBuffer, 0, readLength);
            String[] readDataParts = readData.split(":", 2);
            if (readDataParts.length != 2) {
                System.out.println("Invalid message format");
                return;
            }

            String cmd = readDataParts[0];
            String filename = readDataParts[1];
            if ("PUT".equals(cmd)) {
                writer.println("READY");
                saveHalfFile(socket, filename);
            } else if ("GET".equals(cmd)) {
                sendHalfFile(socket, filename);
            } else {
                System.out.println("Unknown command");
            }
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    private static void saveHalfFile(Socket socket, String filename) {
        String partFilename = getPartedFilename(filename);
        try (FileOutputStream partFile = new FileOutputStream(partFilename);
             InputStream socketInput = socket.getInputStream()) {

            byte[] partFileBuffer = new byte[1024];
            int partFileLength;
            while ((partFileLength = socketInput.read(partFileBuffer)) != -1) {
                partFile.write(partFileBuffer, 0, partFileLength);
            }
        } catch (IOException e) {
            System.out.println("Error handling file: " + e.getMessage());
        }
    }

    private static void sendHalfFile(Socket socket, String filename) {
        String partFilename = getPartedFilename(filename);
        try (FileInputStream partFile = new FileInputStream(partFilename);
             OutputStream socketOutput = socket.getOutputStream()) {

            byte[] partFileBuffer = new byte[1024];
            int partFileLength;
            while ((partFileLength = partFile.read(partFileBuffer)) != -1) {
                byte[] prefixedBuffer = new byte[partFileLength + 1];
                prefixedBuffer[0] = 'N';
                System.arraycopy(partFileBuffer, 0, prefixedBuffer, 1, partFileLength);
                socketOutput.write(prefixedBuffer, 0, prefixedBuffer.length);
            }
            socketOutput.write("EOF".getBytes());
        } catch (FileNotFoundException e) {
            System.out.println("Error opening file: " + e.getMessage());
        } catch (IOException e) {
            System.out.println("Error handling file: " + e.getMessage());
        }
    }

    private static String getPartedFilename(String filename) {
        String ext = filename.substring(filename.lastIndexOf('.'));
        String base = filename.substring(0, filename.lastIndexOf('.'));
        return base + filenameSuffix + ext;
    }

    private static void exitError(String msg) {
        System.err.println("Error: " + msg);
        System.exit(1);
    }
}
