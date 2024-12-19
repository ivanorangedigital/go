const ws = new WebSocket("ws://localhost:4000");
ws.addEventListener("message", () => {
  location.reload();
});
