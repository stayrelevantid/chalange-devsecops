#!/usr/bin/env python3
"""
Security Alert Webhook Receiver — Simulates n8n webhook + IF node logic.

n8n Docker image pull timed out (slow network, ~400MB image).
This Python script achieves the same learning objectives:
  1. Webhook trigger — receives Falco alerts + SAST/SCA/IaC findings via HTTP POST
  2. IF node logic — routes CRITICAL alerts to Slack (Day 48), others to log
  3. AI remediation — LLM auto-summarizes findings and posts to Slack (Day 49)
  4. Audit log — all alerts recorded with timestamp

Endpoints:
  POST /webhook/falco-alert   — Falco runtime alerts (Day 42/48)
  POST /webhook/sast-alert    — Scanner findings: Semgrep/Checkov/Trivy (Day 49)
  GET  /health                — Health check
  GET  /alerts                — Last 20 alerts from log

Run: python3 webhook_receiver.py --port 5678
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
logger = logging.getLogger("sec-webhook")

# Load secrets from .env file in same directory
SLACK_WEBHOOK_URL = None
OPENCODE_ZEN_API_KEY = None
env_path = os.path.join(SCRIPT_DIR, ".env")
if os.path.exists(env_path):
    with open(env_path) as f:
        for line in f:
            line = line.strip()
            if line.startswith("SLACK_WEBHOOK_URL="):
                SLACK_WEBHOOK_URL = line.split("=", 1)[1].strip("\"'")
            elif line.startswith("OPENCODE_ZEN_API_KEY="):
                OPENCODE_ZEN_API_KEY = line.split("=", 1)[1].strip("\"'")

if SLACK_WEBHOOK_URL:
    logger.info(f"Slack webhook configured: {SLACK_WEBHOOK_URL[:40]}...")
else:
    logger.warning("SLACK_WEBHOOK_URL not found — Slack alerts disabled")

if OPENCODE_ZEN_API_KEY:
    logger.info(f"OpenCode Zen API key configured: {OPENCODE_ZEN_API_KEY[:12]}...")
else:
    logger.warning("OPENCODE_ZEN_API_KEY not found — AI remediation disabled")

ZEN_ENDPOINT = "https://opencode.ai/zen/v1/chat/completions"
ZEN_MODEL = "deepseek-v4-flash-free"
SEVERITY_ROUTE_TO_LLM = {"critical", "high", "medium"}


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


def post_to_slack_remediation(scanner, finding, ai_analysis):
    if not SLACK_WEBHOOK_URL:
        logger.warning("Slack disabled — remediation not sent")
        return False

    title = finding.get("title", "unknown")
    severity = finding.get("severity", "unknown")
    file_loc = finding.get("file", "unknown")

    slack_payload = {
        "text": (
            f"🤖 *AI Remediation: {title}*\n"
            f"*Scanner:* {scanner} | *Severity:* {severity} | *File:* {file_loc}\n"
            f"*AI Analysis:*\n{ai_analysis}"
        ),
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
        logger.info(f"Slack POST {resp.status} — remediation '{title}' sent to #security-alerts")
        return True
    except urllib.error.URLError as e:
        logger.error(f"Slack remediation POST failed: {e}")
        return False


def normalize_finding(raw, scanner):
    """Normalize Semgrep/Checkov/Trivy finding to unified format."""
    if scanner == "semgrep":
        check_id = raw.get("check_id", "unknown")
        path = raw.get("path", "unknown")
        line = raw.get("start", {}).get("line", "?")
        extra = raw.get("extra", {})
        refs = extra.get("metadata", {}).get("references", [])
        return {
            "id": check_id,
            "title": check_id.split(".")[-1] if "." in check_id else check_id,
            "severity": extra.get("severity", "unknown"),
            "file": f"{path}:{line}",
            "description": extra.get("message", "no description"),
            "references": refs,
        }
    elif scanner == "checkov":
        check_id = raw.get("check_id", "unknown")
        return {
            "id": check_id,
            "title": raw.get("check_name", check_id),
            "severity": raw.get("severity") or "MEDIUM",
            "file": raw.get("file_path", "unknown"),
            "description": f"Check {check_id} failed on resource {raw.get('resource', 'unknown')}",
            "references": [raw.get("guideline", "")] if raw.get("guideline") else [],
        }
    elif scanner == "trivy":
        vid = raw.get("VulnerabilityID", "unknown")
        refs = raw.get("References", []) or []
        return {
            "id": vid,
            "title": raw.get("Title", vid) or f"{raw.get('PkgName', 'unknown')} vulnerability",
            "severity": raw.get("Severity", "unknown"),
            "file": raw.get("PkgName", "unknown"),
            "description": raw.get("Description", "no description"),
            "references": refs,
        }
    else:
        return {
            "id": raw.get("id", "unknown"),
            "title": raw.get("title", "unknown"),
            "severity": raw.get("severity", "unknown"),
            "file": raw.get("file", "unknown"),
            "description": raw.get("description", "no description"),
            "references": raw.get("references", []),
        }


def call_llm(finding):
    """Call OpenCode Zen LLM to generate remediation advice."""
    if not OPENCODE_ZEN_API_KEY:
        logger.warning("Zen API key not set — skipping LLM call")
        return None

    refs_str = "\n".join(finding.get("references", [])) or "none"
    prompt = (
        f"Berikut temuan keamanan dari scanner:\n"
        f"Judul: {finding['title']}\n"
        f"Severity: {finding['severity']}\n"
        f"File: {finding['file']}\n"
        f"Deskripsi: {finding['description']}\n"
        f"Referensi: {refs_str}\n\n"
        f"Berikan:\n"
        f"1. Ringkasan bahaya (max 2 kalimat)\n"
        f"2. Contoh kode perbaikan (jika ada, singkat)\n"
    )

    payload = {
        "model": ZEN_MODEL,
        "messages": [
            {"role": "system", "content": "Kamu adalah security expert yang memberikan saran perbaikan singkat dan actionable."},
            {"role": "user", "content": prompt},
        ],
        "max_tokens": 512,
        "temperature": 0.7,
    }

    data = json.dumps(payload).encode("utf-8")
    req = urllib.request.Request(
        ZEN_ENDPOINT,
        data=data,
        headers={
            "Content-Type": "application/json",
            "Authorization": f"Bearer {OPENCODE_ZEN_API_KEY}",
            "User-Agent": "SecureBank-WebhookReceiver/1.0",
        },
        method="POST",
    )

    try:
        resp = urllib.request.urlopen(req, timeout=30)
        result = json.loads(resp.read().decode("utf-8"))
        content = result["choices"][0]["message"]["content"]
        logger.info(f"LLM response received ({len(content)} chars)")
        return content
    except urllib.error.URLError as e:
        logger.error(f"LLM call failed: {e}")
        return None
    except (KeyError, IndexError) as e:
        logger.error(f"LLM response parse error: {e}")
        return None


class FalcoWebhookHandler(BaseHTTPRequestHandler):
    def _send_response(self, code, body):
        self.send_response(code)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(json.dumps(body).encode())

    def do_POST(self):
        if self.path == "/webhook/falco-alert":
            self._handle_falco_alert()
        elif self.path == "/webhook/sast-alert":
            self._handle_sast_alert()
        else:
            self._send_response(404, {"error": "not found"})

    def _handle_falco_alert(self):
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

    def _handle_sast_alert(self):
        content_length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(content_length)

        try:
            payload = json.loads(body) if body else {}
        except json.JSONDecodeError:
            logger.warning("Received invalid JSON payload for SAST alert")
            self._send_response(400, {"error": "invalid json"})
            return

        scanner = payload.get("scanner", "unknown")
        raw_findings = payload.get("findings", [])

        # Handle raw scanner reports (auto-detect format)
        if not raw_findings and "results" in payload:
            if isinstance(payload["results"], list):
                raw_findings = payload["results"]
            elif isinstance(payload["results"], dict):
                raw_findings = payload["results"].get("failed_checks", [])

        logger.info(f"SAST ALERT RECEIVED | scanner={scanner} findings={len(raw_findings)}")

        results = []
        for raw in raw_findings:
            finding = normalize_finding(raw, scanner)
            sev = finding["severity"].lower()
            logger.info(f"  -> Finding: {finding['title']} | severity={sev} | file={finding['file']}")

            if sev in SEVERITY_ROUTE_TO_LLM:
                logger.info(f"  -> IF TRUE: {sev} finding — calling LLM + Slack")
                ai_analysis = call_llm(finding)
                if ai_analysis:
                    post_to_slack_remediation(scanner, finding, ai_analysis)
                    results.append({
                        "finding": finding["title"],
                        "severity": sev,
                        "routed_to": "llm+slack",
                        "ai_chars": len(ai_analysis),
                    })
                else:
                    logger.warning(f"  -> LLM call failed for '{finding['title']}'")
                    results.append({
                        "finding": finding["title"],
                        "severity": sev,
                        "routed_to": "llm_failed",
                    })
            else:
                logger.info(f"  -> IF FALSE: {sev} finding — log only (below MEDIUM threshold)")
                results.append({
                    "finding": finding["title"],
                    "severity": sev,
                    "routed_to": "log",
                })

        self._send_response(200, {
            "status": "received",
            "scanner": scanner,
            "findings_processed": len(results),
            "results": results,
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
    logger.info(f"Security webhook receiver listening on {args.host}:{args.port}")
    logger.info(f"  POST /webhook/falco-alert  — Falco runtime alerts (Day 42/48)")
    logger.info(f"  POST /webhook/sast-alert   — Scanner findings: Semgrep/Checkov/Trivy (Day 49)")
    logger.info(f"  GET  /health               — Health check")
    logger.info(f"  GET  /alerts               — Last 20 alerts")
    logger.info(f"Alert log: {LOG_FILE}")
    slack_status = f"Slack: {'connected' if SLACK_WEBHOOK_URL else 'DISABLED'}"
    zen_status = f"Zen LLM: {'connected' if OPENCODE_ZEN_API_KEY else 'DISABLED'}"
    logger.info(f"Falco routing: CRITICAL -> Slack ({slack_status}), others -> log")
    logger.info(f"SAST routing: CRITICAL/HIGH/MEDIUM -> LLM + Slack ({zen_status}), others -> log")

    try:
        server.serve_forever()
    except KeyboardInterrupt:
        logger.info("Shutting down webhook receiver")
        server.shutdown()


if __name__ == "__main__":
    main()
