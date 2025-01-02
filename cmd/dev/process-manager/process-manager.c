#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

#define LINE_LEN 256

int main() {
  FILE *fptr = fopen("./cmd/dev/process-manager/process-manager.txt", "r");
  if (fptr == NULL) {
    perror("Error during opening file");
    return 1;
  }

  pid_t *pids = NULL, *pid = NULL;

  char line[LINE_LEN + 1];
  while (fgets(line, sizeof(line), fptr)) {
    // remove last char if it's new line char
    char *last = line + strlen(line) - 1;
    if (*last == '\n')
      *last = '\0';

    // first memory allocation
    if (pids == NULL) {
      pid = pids = (pid_t *)malloc(sizeof(pid_t));

      if (pids == NULL) {
        perror("Error during memory allocation to pids");
        return 1;
      }
    } else {
      // then realloc it
      int ln = pid - pids + 1;
      pids = realloc(pids, sizeof(pid_t) * (ln + 1));
      if (pids == NULL) {
        perror("Error during memory realloc to pids");
        return 1;
      }
      pid = pids + ln;
    }

    // sleep for 1 sec
    sleep(1);

    // create fork
    *pid = fork();

    if (*pid < 0) {
      perror("Error during fork");
      free(pids);
      fclose(fptr);
      return 1;
    }

    if (*pid == 0) {
      execl("/bin/bash", "bash", "-c", line, NULL);

      perror("Execution command was failed");
      exit(1);
    }
  }

  fclose(fptr);

  // wait for all process
  for (pid_t *_pid = pids; _pid < pid + 1; _pid++) {
    waitpid(*_pid, NULL, 0);
  }

  free(pids);
  return 0;
}
