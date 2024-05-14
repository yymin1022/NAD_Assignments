/**
 * ChatClient.c
 * ID : 20194094
 * Name : Yongmin Yoo
 **/

#include <arpa/inet.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>
#include <signal.h>
#include <sys/select.h>
#include <time.h>

#define BUF_SIZE 1024
#define SERVER_NAME "127.0.0.1"
#define SERVER_PORT 14094

int exit_error(char *err_msg);
void handle_sigint(int sig);
void close_connection();
void encode_command(const char *command, char *cmd, char *extra);

int             socket_fd;
int             PING_MODE = 0;
struct timespec PING_START_TIME;

int main(int argc, char *argv[])
{
    char                buffer[BUF_SIZE];
    fd_set              read_fds;
    struct sockaddr_in  server_addr;
    struct timespec     PING_END_TIME;

    if (argc < 2)
    {
        write(2, "Usage: ./ChatClient <nickname>\n", 31);
        exit(EXIT_FAILURE);
    }

    char *client_nickname = argv[1];

    socket_fd = socket(AF_INET, SOCK_STREAM, 0);
    if (socket_fd < 0)
        return (exit_error("Socket Error"));

    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(SERVER_PORT);
    if (inet_pton(AF_INET, SERVER_NAME, &server_addr.sin_addr) <= 0)
        return (exit_error("inet_pton Error"));
    if (connect(socket_fd, (struct sockaddr *)&server_addr, sizeof(server_addr)) < 0)
        return (exit_error("Server Connection Failed"));

    signal(SIGINT, handle_sigint);

    dprintf(socket_fd, "%s\n", client_nickname);
    while (1)
    {
        FD_ZERO(&read_fds);
        FD_SET(socket_fd, &read_fds);
        FD_SET(STDIN_FILENO, &read_fds);

        int max_fd = socket_fd > STDIN_FILENO ? socket_fd : STDIN_FILENO;
        int activity = select(max_fd + 1, &read_fds, NULL, NULL, NULL);

        if (activity < 0)
        {
            exit_error("select error");
            break;
        }

        if (FD_ISSET(socket_fd, &read_fds))
        {
            int n = read(socket_fd, buffer, BUF_SIZE - 1);
            if (n < 0)
            {
                exit_error("Read error");
                break;
            }
            else if (n == 0)
                break;
            buffer[n] = '\0';

            if (PING_MODE)
            {
                PING_MODE = 0;
                clock_gettime(CLOCK_MONOTONIC, &PING_END_TIME);
                long rtt_value = ((PING_END_TIME.tv_sec - PING_START_TIME.tv_sec) * 1000 + (PING_END_TIME.tv_nsec - PING_START_TIME.tv_nsec)) / 1000;
                printf("RTT is %ldms\n", rtt_value);
            }
            else
                write(1, buffer + 1, strlen(buffer) - 1);

            if (buffer[0] == 'K')
            {
                close_connection();
                exit(0);
            }
        }

        if (FD_ISSET(STDIN_FILENO, &read_fds))
        {
            if (fgets(buffer, BUF_SIZE, stdin) != NULL)
            {
                buffer[strcspn(buffer, "\n")] = 0; // Remove newline character
                if (buffer[0] == '\\')
                {
                    char cmd[2], extra[BUF_SIZE];
                    encode_command(buffer, cmd, extra);
                    if (strcmp(cmd, "P") == 0)
                    {
                        PING_MODE = 1;
                        clock_gettime(CLOCK_MONOTONIC, &PING_START_TIME);
                        dprintf(socket_fd, "P ping\n");
                    }
                    else if (strcmp(cmd, "Q") == 0)
                    {
                        close_connection();
                        exit(0);
                    }
                    else if (cmd[0] != '\0')
                        dprintf(socket_fd, "%s %s\n", cmd, extra);
                    else
                        write(1, "Invalid Command.\n", 17);
                }
                else
                    dprintf(socket_fd, "M%s\n", buffer);
            }
        }
    }

    close_connection();
    return 0;
}

void handle_sigint(int sig)
{
    if (sig == SIGINT)
    {
        close_connection();
        exit(0);
    }
}

void close_connection()
{
    write(1, "\rClosing Client Program...\nBye bye~\n", 36);
    if (socket_fd >= 0)
    {
        dprintf(socket_fd, "Q");
        close(socket_fd);
    }
}

int exit_error(char *err_msg)
{
    write(2, "Error: ", 7);
    write(2, err_msg, strlen(err_msg));
    write(2, "\n", 1);
    return -1;
}

void encode_command(const char *command, char *cmd, char *extra)
{
    char *space = strchr(command, ' ');

    if (space != NULL)
    {
        *space = '\0';
        strcpy(extra, space + 1);
    }
    else
        extra[0] = '\0';

    if (strcmp(command, "\\ls") == 0)
        strcpy(cmd, "L");
    else if (strcmp(command, "\\ping") == 0)
        strcpy(cmd, "P");
    else if (strcmp(command, "\\quit") == 0)
        strcpy(cmd, "Q");
    else if (strcmp(command, "\\secret") == 0)
        strcpy(cmd, "S");
    else if (strcmp(command, "\\except") == 0)
        strcpy(cmd, "E");
    else
        cmd[0] = '\0';
}
