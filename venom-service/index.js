import venom from "venom-bot";
import http from "http";

let client = null;

venom
  .create({
    session: "autoflow",
    multidevice: true,
    logQR: true,
    headless: false,
    useChrome: true,
    disableInstanceLock: true,
    updatesLog: true,

    browserArgs: [
      "--no-sandbox",
      "--disable-setuid-sandbox",
      "--disable-dev-shm-usage",
      "--disable-background-timer-throttling",
      "--disable-backgrounding-occluded-windows",
      "--disable-renderer-backgrounding",
      "--disable-gpu",
    ]
  })
  .then((c) => {
    client = c;
    console.log("Venom client connected & ready!");

    // Listen for messages
    client.onMessage(async (msg) => {
      console.log("Incoming:", msg.from, msg.body);

      // Forward to GO (localhost)
      sendToGoBackend({
        From: msg.from,
        Body: msg.body,
      });
    });
  })
  .catch((err) => console.error("Venom FAILED:", err));


// --- SEND MESSAGE TO GO BACKEND ---
function sendToGoBackend(json) {
  const body = JSON.stringify(json);

  const req = http.request(
    {
      hostname: "localhost",
      port: 8080,
      path: "/webhook/whatsapp",
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Content-Length": Buffer.byteLength(body),
      }
    },
    (res) => {
      res.on("data", () => {});
    }
  );

  req.on("error", (err) =>
    console.log("Go backend request failed:", err.message)
  );

  req.write(body);
  req.end();
}


// --- SEND MESSAGE API ----
export async function sendMessage(to, message) {
  if (!client) throw new Error("Client not ready");

  const number = to.replace(/\D/g, "") + "@c.us";
  await client.sendText(number, message);
  return true;
}

console.log("⚡ Venom service booting...");
