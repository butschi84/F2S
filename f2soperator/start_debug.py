#!/usr/bin/env python3
import subprocess
import os
import random

# Helper function to generate a random port
def generate_random_port():
    return random.randint(1025, 60000)

# Use local Kubernetes config
os.environ["KUBECONFIG"] = os.path.expanduser("~/.kube/config")

# Terminate all previous port forwardings
print("Killing previous port-forwarding threads")
pids = subprocess.run(["pgrep", "-f", "kubectl port-forward"], capture_output=True, text=True).stdout.splitlines()
if not pids:
    print("No kubectl port-forward processes found.")
else:
    print("Killing kubectl port-forward processes...")
    for pid in pids:
        subprocess.run(["kill", "-9", pid])
    print("All kubectl port-forward processes killed.")

# Forward ports
print("Forward ports from Kubernetes services")
subprocess.Popen(["kubectl", "port-forward", "-n", "f2s", "service/prometheus-service", "9090:9090"], stdout=subprocess.DEVNULL)
subprocess.Popen(["kubectl", "port-forward", "-n", "f2s", "service/f2s-memberlist", "7079:7079"], stdout=subprocess.DEVNULL)
subprocess.Popen(["kubectl", "port-forward", "-n", "f2s", "service/f2s-api", "8080:8080"], stdout=subprocess.DEVNULL)

# Use forwarded services
print("Set environment variables for f2s debugging")
os.environ["Prometheus_URL"] = "localhost:9090"
os.environ["MEMBERLIST_ADDRESS"] = "localhost:7079"

# Forward port of each f2s-containers service
print("Forward ports of f2s-containers services")
services = subprocess.run(["kubectl", "get", "services", "-n", "f2s-containers", "-o=jsonpath=\"{range .items[*]}{.metadata.name}{\"\\n\"}{end}\""], capture_output=True, text=True).stdout.splitlines()
print(f"found {len(services)} f2s-containers services:")
for service_name in services:
    print(f"  Starting port forwarding for service: {service_name}")

    # Generate a random local port
    local_port = generate_random_port()

    # Forward the service port
    subprocess.Popen(["kubectl", "port-forward", "-n", "f2s-containers", f"service/{service_name}", f"{local_port}:9092"], stdout=subprocess.DEVNULL)

    # Store an environment variable
    service_name_with_underscores = service_name.replace("-", "_")
    uppercase_service_name = service_name_with_underscores.upper()
    print(f" => F2S_SERVICE_{uppercase_service_name}=127.0.0.1:{local_port}")
    os.environ[f"F2S_SERVICE_{uppercase_service_name}"] = f"127.0.0.1:{local_port}"

print("Port forwarding started for all services.")

# Start f2s application
print("Starting f2s...")
subprocess.run(["go", "run", "main.go"])
