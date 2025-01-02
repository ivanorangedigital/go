import { WebSocketServer, WebSocket } from "ws";
import { createServer } from "http";

const server = createServer();
const wss = new WebSocketServer({ server });
wss.on("connection", (ws) => {
  ws.on("error", console.error);
  ws.on("message", () => {
    wss.clients.forEach((client) => {
      if (client.readyState === WebSocket.OPEN) {
        client.send();
      }
    });
  });
});

// start server
server.listen(4000, "localhost");
