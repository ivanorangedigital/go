#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main() {
  pid_t watchPrebuildPID = fork();
  if (watchPrebuildPID < 0) {
    return 1;
  }

  if (watchPrebuildPID == 0) {
    char *args[] = {"make", "watch-prebuild", NULL};
    execvp(args[0], args);

    perror("exec 'make watch-prebuild' failed");
    return 1;
  }

  pid_t watchRunPID = fork();
  if (watchRunPID < 0) {
    return 1;
  }

  if (watchRunPID == 0) {
    sleep(1);

    char *args[] = {"make", "watch-run", NULL};
    execvp(args[0], args);

    perror("exec 'make watch-run' failed");
    return 1;
  }

  pid_t watchWebScoketPID = fork();
  if (watchWebScoketPID < 0) {
    return 1;
  }

  if (watchWebScoketPID == 0) {
    sleep(1);

    char *args[] = {"make", "start-websocket-server", NULL};
    execvp(args[0], args);

    perror("exec 'make start-websocket-server' failed");
    return 1;
  } else {
    sleep(2);

    char *args[] = {"make", "watch-websocket-ping", NULL};
    execvp(args[0], args);

    perror("exec 'make watch-websocket-ping' failed");
    return 1;
  }

  waitpid(watchPrebuildPID, NULL, 0);
  waitpid(watchRunPID, NULL, 0);
  waitpid(watchWebScoketPID, NULL, 0);

  return 0;
}
