#!/usr/bin/env python3
"""
Falco Alert Webhook Receiver — Simulates n8n webhook + IF node logic.

n8n Docker image pull timed out (slow network, ~400MB image).
This Python script achieves the same learning objectives:
  1. Webhook trigger — receives Falco alerts via HTTP POST
  2. IF node logic — routes CRITICAL alerts to Slack (Day 48), others to log
  3. Audit log — all alerts recorded with timestamp

Run: python3 webhook_receiver.py --port 5678
Test: curl -X POST http://localhost:5678/webhook/falco-alert -d '{"rule":"test","priority":"Critical"}'
"""
import argparse
import json
import logging
import os
from datetime import datetime, timezone
from http.server import HTTPServer, BaseHTTPRequestHandler

LOG_DIR = os.path.join(os.path.dirname(os.path.abspath(__file__)), "logs")
os.makedirs(LOG_DIR, exist_ok=True)
LOG_FILE = os.path.join(LOG_DIR, "falco-alerts.log")

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
    handlers=[
        logging.FileHandler(LOG_FILE),
        logging.StreamHandler(),
    ],
)
logger = logging.getLogger("falco-webhook")


class FalcoWebhookHandler(BaseHTTPRequestHandler):
    def _send_response(self, code, body):
        self.send_response(code)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(json.dumps(body).encode())

    def do_POST(self):
        if self.path != "/webhook/falco-alert":
            self._send_response(404, {"error": "not found"})
            return

        content_length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(content_length)

        try:
            alert = json.loads(body) if body else {}
        except json.JSONDecodeError:
            logger.warning("Received invalid JSON payload")
            self._send_response(400, {"error": "invalid json"})
            return

        rule = alert.get("rule", "unknown")
        priority = alert.get("priority", "unknown")
        output = alert.get("output", "")
        ts = alert.get("time", datetime.now(timezone.utc).isoformat())
        pod = alert.get("output_fields", {}).get("k8s.pod.name", "unknown")
        ns = alert.get("output_fields", {}).get("k8s.ns.name", "unknown")
        proc = alert.get("output_fields", {}).get("proc.name", "unknown")

        logger.info(f"ALERT RECEIVED | priority={priority} rule={rule} pod={pod} ns={ns} proc={proc}")

        # IF node logic — CRITICAL routes to Slack (Day 48), others just log
        if priority.lower() == "critical":
            logger.info(f"  -> IF TRUE: CRITICAL alert — route to Slack #security-alerts (Day 48)")
            logger.info(f"  -> Slack payload prepared: '{rule}' in pod '{pod}' (ns={ns})")
            # Day 48 will add: requests.post(slack_webhook_url, json=slack_payload)
        else:
            logger.info(f"  -> IF FALSE: {priority} alert — log only (below CRITICAL threshold)")

        # Acknowledge receipt
        self._send_response(200, {
            "status": "received",
            "rule": rule,
            "priority": priority,
            "routed_to": "slack" if priority.lower() == "critical" else "log",
            "timestamp": ts,
        })

    def do_GET(self):
        if self.path == "/health":
            self._send_response(200, {"status": "healthy"})
        elif self.path == "/alerts":
            # Return last 20 alerts from log
            alerts = []
            if os.path.exists(LOG_FILE):
                with open(LOG_FILE, "r") as f:
                    alerts = [line.strip() for line in f.readlines()[-20:]]
            self._send_response(200, {"alerts": alerts, "count": len(alerts)})
        else:
            self._send_response(404, {"error": "not found"})

    def log_message(self, format, *args):
        pass  # suppress default HTTP logging


def main():
    parser = argparse.ArgumentParser(description="Falco Alert Webhook Receiver")
    parser.add_argument("--port", type=int, default=5678, help="Listen port (default: 5678)")
    parser.add_argument("--host", type=str, default="0.0.0.0", help="Listen host (default: 0.0.0.0)")
    args = parser.parse_args()

    server = HTTPServer((args.host, args.port), FalcoWebhookHandler)
    logger.info(f"Falco webhook receiver listening on {args.host}:{args.port}")
    logger.info(f"Webhook endpoint: POST http://<host>:{args.port}/webhook/falco-alert")
    logger.info(f"Health check: GET http://<host>:{args.port}/health")
    logger.info(f"Alert log: {LOG_FILE}")
    logger.info("IF node logic: CRITICAL -> Slack (Day 48), others -> log only")

    try:
        server.serve_forever()
    except KeyboardInterrupt:
        logger.info("Shutting down webhook receiver")
        server.shutdown()


if __name__ == "__main__":
    main()
