/**
 * ChatServer.c
 * ID : 20194094
 * Name : Yongmin Yoo
 **/

#include <arpa/inet.h>
#include <netinet/in.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/select.h>
#include <unistd.h>

#define BUF_SIZE 1024
#define MAX_CLIENT 8
#define NICK_SIZE 32
#define SERVER_PORT 24094

int     check_message_avail(char *message);
int     check_nick_avail(char *new_nick);
int     find_client_index(int fd);
int     setup_server();
char    *get_ip_port(int fd, int is_client);
char    *str_toupper(char *str);
char    *trim_newline(char *str);
void    broadcast_message(const char *message, int sender_fd);
void    close_server(int sig);
void    exclude_nick(const char *message, const char *exclude_nick, const char *from_nick);
void    exit_error(char *err_msg);
void    handle_new_connection(int server_fd);
void    handle_client_message(int client_fd);
void    remove_client(int client_fd);
void    run_server(int server_fd);
void    send_to_nick(const char *message, const char *target_nick, const char *from_nick);

typedef struct {
    int         fd;
    char        nickname[NICK_SIZE];
} client_t;

int         client_count = 0;
int         client_fd_max;
client_t    clients[MAX_CLIENT];
fd_set      current_sockets;

int main()
{
    int server_fd;

    server_fd = setup_server();
    if (server_fd < 0)
        exit_error("Failed to initialize server");

    signal(SIGINT, close_server);
    run_server(server_fd);
    return 0;
}

int setup_server()
{
    int                 server_fd;
    int                 server_option = 1;
    struct sockaddr_in  server_addr;

    server_fd = socket(AF_INET, SOCK_STREAM, 0);
    if (server_fd == -1)
    {
        perror("Socket creation failed");
        return -1;
    }

    setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR | SO_REUSEPORT, &server_option, sizeof(server_option));
    server_addr.sin_family = AF_INET;
    server_addr.sin_addr.s_addr = INADDR_ANY;
    server_addr.sin_port = htons(SERVER_PORT);

    if (bind(server_fd, (struct sockaddr *)&server_addr, sizeof(server_addr)) < 0)
    {
        perror("Bind failed");
        return -1;
    }

    if (listen(server_fd, MAX_CLIENT) < 0)
    {
        perror("Listen failed");
        return -1;
    }

    client_fd_max = server_fd;
    FD_ZERO(&current_sockets);
    FD_SET(server_fd, &current_sockets);

    printf("Server is ready to receive on port %d\n", SERVER_PORT);
    return server_fd;
}

void    run_server(int server_fd)
{
    fd_set  client_fds;

    while (1)
    {
        client_fds = current_sockets;
        if (select(client_fd_max + 1, &client_fds, NULL, NULL, NULL) < 0)
        {
            perror("Select error");
            return;
        }

        for (int i = 0; i <= client_fd_max; i++)
        {
            if (FD_ISSET(i, &client_fds))
            {
                if (i == server_fd)
                    handle_new_connection(server_fd);
                else
                    handle_client_message(i);
            }
        }
    }
}

void    handle_new_connection(int server_fd)
{

    char                client_nick[BUF_SIZE];
    char                client_welcome_msg[BUF_SIZE];
    int                 client_nick_read;
    int                 client_fd;
    struct sockaddr_in  client_addr;
    client_t            new_client;
    socklen_t           addr_size;

    addr_size = sizeof(client_addr);
    client_fd = accept(server_fd, (struct sockaddr *)&client_addr, &addr_size);

    if (client_fd < 0)
    {
        perror("Accept failed");
        return;
    }

    if (client_count >= MAX_CLIENT)
    {
        char *message = "KChatting Room is Full. Cannot connect\n";
        send(client_fd, message, strlen(message), 0);
        close(client_fd);
        return;
    }

    client_nick_read = read(client_fd, client_nick, BUF_SIZE);
    if (client_nick_read <= 0)
    {
        char *message = "KSocket Error\n";
        send(client_fd, message, strlen(message), 0);
        remove_client(client_fd);
        return;
    }

    client_nick[strlen(client_nick) - 1] = '\0';

    if (check_nick_avail(client_nick) == 0)
    {
        char *message = "KNickname is already used by another user. Cannot connect\n";
        send(client_fd, message, strlen(message), 0);
        remove_client(client_fd);
        return;
    }

    new_client.fd = client_fd;
    strcpy(new_client.nickname, client_nick);
    clients[client_count] = new_client;

    FD_SET(client_fd, &current_sockets);
    if (client_fd > client_fd_max)
        client_fd_max = client_fd;

    char *client_ip_port = get_ip_port(client_fd, 1);
    char *server_ip_port = get_ip_port(server_fd, 0);
    sprintf(client_welcome_msg, "M[Welcome %s to CAU Net-Class Chat Room at %s.]\n", client_nick, server_ip_port);
    send(client_fd, client_welcome_msg, strlen(client_welcome_msg), 0);
    sprintf(client_welcome_msg, "M[There are %d users in the room]\n", client_count);
    send(client_fd, client_welcome_msg, strlen(client_welcome_msg), 0);

    client_count++;
    printf("%s Joined from %s. There are %d users in the room.\n", client_nick, client_ip_port, client_count);
    free(client_ip_port);
    free(server_ip_port);
}

void    handle_client_message(int client_fd)
{
    char    client_msg[BUF_SIZE];
    char    client_msg_copy[BUF_SIZE];
    char    *client_command;
    char    *client_command_extra;
    char    *client_command_message;
    char    *client_command_target;
    int     client_msg_size;

    client_msg_size = read(client_fd, client_msg, BUF_SIZE);
    if (client_msg_size <= 0)
    {
        char client_nick[NICK_SIZE];
        char send_buf[BUF_SIZE];
        char *message = "KSocket Error\n";

        send(client_fd, message, strlen(message), 0);
        strcpy(client_nick, clients[find_client_index(client_fd)].nickname);
        remove_client(client_fd);

        snprintf(send_buf, sizeof(send_buf), "M%s has left the chat.", client_nick);
        broadcast_message(send_buf, client_fd);
        printf("%s left the room. There are %d users in the room.\n", client_nick, client_count);

        return;
    }

    client_msg[client_msg_size] = '\0';
    trim_newline(client_msg);

    strcpy(client_msg_copy, client_msg);
    client_command = strtok(client_msg_copy, " ");
    client_command_extra = strtok(NULL, "");
    if (strcmp(client_command, "L") == 0)
    {
        for (int i = 0; i < client_count; i++)
        {
            char line[BUF_SIZE];
            char *client_ip_port = get_ip_port(clients[i].fd, 1);

            snprintf(line, sizeof(line), "I%s - %s\n", clients[i].nickname, client_ip_port);
            send(client_fd, line, strlen(line), 0);
            free(client_ip_port);
        }
    }
    else if (strcmp(client_command, "P") == 0)
        send(client_fd, "P\n", 2, 0);
    else if (strcmp(client_command, "Q") == 0)
    {
        char client_nick[NICK_SIZE];
        char send_buf[BUF_SIZE];

        strcpy(client_nick, clients[find_client_index(client_fd)].nickname);
        remove_client(client_fd);

        snprintf(send_buf, sizeof(send_buf), "M%s has left the chat.", client_nick);
        broadcast_message(send_buf, client_fd);
        printf("%s left the room. There are %d users in the room.\n", client_nick, client_count);
    }
    else if (strcmp(client_command, "S") == 0 && client_command_extra != NULL)
    {
        client_command_target = strtok(client_command_extra, " ");
        client_command_message = strtok(NULL, "");

        for (int i = 0; i < client_count; i++)
        {
            if (clients[i].fd == client_fd)
                send_to_nick(client_command_message, client_command_target, clients[i].nickname);
        }

    }
    else if (strcmp(client_command, "E") == 0 && client_command_extra != NULL)
    {
        client_command_target = strtok(client_command_extra, " ");
        client_command_message = strtok(NULL, "");

        for (int i = 0; i < client_count; i++)
        {
            if (clients[i].fd == client_fd)
                exclude_nick(client_command_message, client_command_target, clients[i].nickname);
        }
    }
    else
    {
        strcpy(client_msg_copy, client_msg);
        if (check_message_avail(client_msg_copy) == 0)
        {
            char client_nick[NICK_SIZE];
            char send_buf[BUF_SIZE];

            strcpy(client_nick, clients[find_client_index(client_fd)].nickname);
            sprintf(send_buf, "KBanned Keyword.\n");
            send(client_fd, send_buf, strlen(send_buf), 0);
            remove_client(client_fd);

            sprintf(send_buf, "M[%s is disconnected. There are %d users in the room.]", client_nick, client_count);
            printf("[%s is disconnected. There are %d users in the room.]\n", client_nick, client_count);
            broadcast_message(send_buf, client_fd);
        }
        else
            broadcast_message(client_msg, client_fd);
    }
}

void    broadcast_message(const char *message, int sender_fd)
{
    char sender_nick[NICK_SIZE];
    char send_buf[BUF_SIZE];

    strcpy(sender_nick, clients[find_client_index(sender_fd)].nickname);
    snprintf(send_buf, sizeof(send_buf), "M%s> %s\n", sender_nick, message + 1);
    for (int i = 0; i < client_count; i++)
    {
        if (clients[i].fd != sender_fd)
            send(clients[i].fd, send_buf, strlen(send_buf), 0);
    }
}

int check_message_avail(char *message)
{
    if (strstr(str_toupper(message), "I HATE PROFESSOR") != 0)
        return 0;
    return 1;
}

int check_nick_avail(char *new_nick)
{
    for (int i = 0; i < client_count; i++)
    {
        if (strcmp(clients[i].nickname, new_nick) == 0)
            return 0;
    }
    return 1;
}

int find_client_index(int fd)
{
    for (int i = 0; i < client_count; i++)
    {
        if (clients[i].fd == fd)
            return i;
    }
    return -1;
}

void    send_to_nick(const char *message, const char *target_nick, const char *from_nick)
{
    char send_buf[BUF_SIZE];

    for (int i = 0; i < client_count; i++)
    {
        if (strcmp(clients[i].nickname, target_nick) == 0)
        {
            snprintf(send_buf, sizeof(send_buf), "M%s> %s\n", from_nick, message);
            send(clients[i].fd, send_buf, strlen(send_buf), 0);
            break;
        }
    }
}

void    exclude_nick(const char *message, const char *exclude_nick, const char *from_nick)
{
    char send_buf[BUF_SIZE];

    for (int i = 0; i < client_count; i++)
    {
        if (strcmp(clients[i].nickname, exclude_nick) != 0)
        {
            snprintf(send_buf, sizeof(send_buf), "M%s> %s\n", from_nick, message);
            send(clients[i].fd, send_buf, strlen(send_buf), 0);
        }
    }
}

void    remove_client(int client_fd)
{
    for (int i = 0; i < client_count; i++)
    {
        if (clients[i].fd == client_fd)
        {
            for (int j = i; j < client_count - 1; j++)
            {
                clients[j].fd = clients[j + 1].fd;
                strcpy(clients[j].nickname, clients[j + 1].nickname);
            }
            client_count--;

            close(client_fd);
            FD_CLR(client_fd, &current_sockets);
            break;
        }
    }
}

char    *get_ip_port(int fd, int is_client)
{
    char                res[30];
    socklen_t           addr_size;
    struct sockaddr_in  addr_info;

    addr_size = sizeof(struct sockaddr_in);
    if (is_client)
        getpeername(fd, (struct sockaddr *)&addr_info, &addr_size);
    else
        getsockname(fd, (struct sockaddr *)&addr_info, &addr_size);
    sprintf(res, "%s:%d", inet_ntoa(addr_info.sin_addr), ntohs(addr_info.sin_port));
    return (strdup(res));
}

char    *str_toupper(char *str)
{
    size_t  i;

    i = 0;
    while(str[i])
    {
        if (str[i] >= 'a' && str[i] <= 'z')
            str[i] = str[i] - 'a' + 'A';
        i++;
    }
    return (str);
}

char    *trim_newline(char *str)
{
    char *pos;
    if ((pos = strchr(str, '\n')) != NULL)
        *pos = '\0';
    return str;
}

void    close_server(int sig)
{
    if (sig == SIGINT)
    {
        printf("\rClosing Server Program...\nBye bye~\n");
        for (int i = 0; i < client_count; i++)
            close(clients[i].fd);
        exit(0);
    }
}

void    exit_error(char *err_msg)
{
    write(2, "Error: ", 7);
    write(2, err_msg, strlen(err_msg));
    write(2, "\n", 1);
    exit(1);
}
