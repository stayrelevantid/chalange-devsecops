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
import urllib.error
import urllib.request
from datetime import datetime, timezone
from http.server import HTTPServer, BaseHTTPRequestHandler

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
LOG_DIR = os.path.join(SCRIPT_DIR, "logs")
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

# Load Slack webhook URL from .env file in same directory
SLACK_WEBHOOK_URL = None
env_path = os.path.join(SCRIPT_DIR, ".env")
if os.path.exists(env_path):
    with open(env_path) as f:
        for line in f:
            line = line.strip()
            if line.startswith("SLACK_WEBHOOK_URL="):
                SLACK_WEBHOOK_URL = line.split("=", 1)[1].strip("\"'")
                break

if SLACK_WEBHOOK_URL:
    logger.info(f"Slack webhook configured: {SLACK_WEBHOOK_URL[:40]}...")
else:
    logger.warning("SLACK_WEBHOOK_URL not found — Slack alerts disabled")


def post_to_slack(alert):
    if not SLACK_WEBHOOK_URL:
        logger.warning("Slack disabled — alert not sent")
        return False

    rule = alert.get("rule", "unknown")
    priority = alert.get("priority", "unknown")
    pod = alert.get("output_fields", {}).get("k8s.pod.name", "unknown")
    ns = alert.get("output_fields", {}).get("k8s.ns.name", "unknown")
    proc = alert.get("output_fields", {}).get("proc.name", "unknown")
    output = alert.get("output", "no output")

    slack_payload = {
        "text": f"🚨 *Falco CRITICAL Alert*\n*Rule:* {rule}\n*Pod:* {pod} (ns: {ns})\n*Process:* {proc}\n*Output:* `{output[:200]}`",
    }

    data = json.dumps(slack_payload).encode("utf-8")
    req = urllib.request.Request(
        SLACK_WEBHOOK_URL,
        data=data,
        headers={"Content-Type": "application/json"},
        method="POST",
    )

    try:
        resp = urllib.request.urlopen(req, timeout=10)
        logger.info(f"Slack POST {resp.status} — alert '{rule}' sent to #security-alerts")
        return True
    except urllib.error.URLError as e:
        logger.error(f"Slack POST failed: {e}")
        return False


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

        # IF node logic — CRITICAL routes to Slack, others just log
        if priority.lower() == "critical":
            logger.info(f"  -> IF TRUE: CRITICAL alert — posting to Slack #security-alerts")
            success = post_to_slack(alert)
            if success:
                logger.info(f"  -> Slack notification sent for rule '{rule}'")
            else:
                logger.warning(f"  -> Slack notification FAILED for rule '{rule}'")
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
    slack_status = f"Slack: {'connected' if SLACK_WEBHOOK_URL else 'DISABLED'}"
    logger.info(f"IF node logic: CRITICAL -> Slack #security-alerts ({slack_status}), others -> log only")

    try:
        server.serve_forever()
    except KeyboardInterrupt:
        logger.info("Shutting down webhook receiver")
        server.shutdown()


if __name__ == "__main__":
    main()
