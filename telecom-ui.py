# telecom-ui.py
import subprocess
import threading
from nicegui import ui

# ------------------ Functions ------------------ #
def tail_logs(container_name, textarea, color):
    """Continuously fetch logs from a Docker container and append to textarea."""
    process = subprocess.Popen(
        ["docker", "logs", "-f", container_name],
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        text=True,
        bufsize=1
    )
    for line in process.stdout:
        # Prefix with container name and colorize
        textarea.value += f"[{container_name}] {line}"
        textarea.update()

def send_event(device_id, cpu_value):
    """Send test telemetry event using requests to Collector API."""
    import requests, time
    payload = {
        "device_id": device_id,
        "timestamp": int(time.time() * 1000),
        "metrics": {"cpu": cpu_value}
    }
    try:
        resp = requests.post("http://localhost:8080/collect", json=payload)
        ui.notify(f"Sent: {payload} → Status {resp.status_code}")
    except Exception as e:
        ui.notify(f"Failed to send: {e}", color="red")

# ------------------ GUI ------------------ #
ui.label("Telecom Kafka Dashboard").classes("text-2xl font-bold mb-4")

with ui.row().classes("mb-4"):
    device_input = ui.input("Device ID", value="device123").props("outlined")
    cpu_input = ui.input("CPU value", value="75").props("outlined")
    ui.button("Send Event", on_click=lambda: send_event(device_input.value, cpu_input.value))

# ------------------ Logs for each consumer ------------------ #
logs_areas = {}
colors = {
    "consumer-metrics": "text-blue-600",
    "consumer-audit": "text-green-600",
    "consumer-alarm": "text-red-600"
}

for container_name in ["consumer-metrics", "consumer-audit", "consumer-alarm"]:
    ui.label(container_name).classes(f"text-lg font-semibold {colors[container_name]}")
    textarea = ui.textarea().props("outlined readonly")
    textarea.classes("h-64 w-full mb-4 font-mono")
    logs_areas[container_name] = textarea

# ------------------ Threads to tail logs ------------------ #
for container_name, textarea in logs_areas.items():
    threading.Thread(target=tail_logs, args=(container_name, textarea, colors[container_name]), daemon=True).start()

# ------------------ Run GUI ------------------ #
ui.run(title="Telecom Kafka GUI", port=8081)
