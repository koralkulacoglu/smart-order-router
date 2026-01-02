import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
import json
import re

data = []
with open('benchmark_results.json', 'r') as f:
    for line in f:
        try:
            event = json.loads(line)
            if event.get('Action') == 'output' and 'Benchmark' in event.get('Output', ''):
                pass
            
            if event.get('Action') == 'run':
                continue
            
            if 'Benchmark' in event.get('Test', '') and event.get('Action') == 'pass':
                pass
                
        except ValueError:
            continue

parsed_records = []

with open('benchmark_results.json', 'r') as f:
    for line in f:
        entry = json.loads(line)
        text = entry.get('Output', '').strip()
        
        match = re.search(r'Benchmark(.+?)-(\d+)\s+(\d+)\s+(\d+(?:\.\d+)?)\s+ns/op\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op', text)
        
        if match:
            impl_name = match.group(1)
            cores = int(match.group(2))
            iterations = int(match.group(3))
            ns_op = float(match.group(4))
            bytes_op = int(match.group(5))
            allocs_op = int(match.group(6))
            
            parsed_records.append({
                "Implementation": impl_name,
                "Cores": cores,
                "Latency (ns)": ns_op,
                "Throughput (Ops/sec)": 1e9 / ns_op,
                "Memory (Bytes)": bytes_op,
                "Allocations": allocs_op
            })

df = pd.DataFrame(parsed_records)

sns.set_theme(style="whitegrid")
plt.figure(figsize=(16, 10))

# Latency vs Cores
plt.subplot(2, 2, 1)
sns.lineplot(data=df, x="Cores", y="Latency (ns)", hue="Implementation", style="Implementation", markers=True, dashes=False)
plt.title("Latency vs CPU Cores (Lower is Better)")
plt.ylabel("Nanoseconds per Operation")
plt.xticks([1, 2, 4, 8, 16, 24, 32])

# Throughput vs Cores
plt.subplot(2, 2, 2)
sns.lineplot(data=df, x="Cores", y="Throughput (Ops/sec)", hue="Implementation", style="Implementation", markers=True, dashes=False)
plt.title("Throughput vs CPU Cores (Higher is Better)")
plt.ylabel("Operations per Second")
plt.xticks([1, 2, 4, 8, 16, 24, 32])

# Memory Allocations
plt.subplot(2, 2, 3)
sns.barplot(data=df, x="Implementation", y="Allocations", hue="Implementation")
plt.title("Memory Allocations per Operation")
plt.ylabel("Allocs / Op")

# Bytes Allocated
plt.subplot(2, 2, 4)
sns.barplot(data=df, x="Implementation", y="Memory (Bytes)", hue="Implementation")
plt.title("Bytes Allocated per Operation")
plt.ylabel("Bytes / Op")

plt.tight_layout()
plt.savefig("benchmark_analysis.png")
print("Plots saved to benchmark_analysis.png")
