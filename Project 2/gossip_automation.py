import subprocess
import statistics
import re

def run_gossip(nodes, topology, algorithm, runs=10):
    times = []
    for _ in range(runs):
        result = subprocess.run(['powershell', '-Command', 
                                 f'Measure-Command {{./Gossip2 {nodes} {topology} {algorithm}}}'], 
                                capture_output=True, text=True)
        
        # Extract TotalMilliseconds from the output
        match = re.search(r'TotalMilliseconds\s*:\s*([\d.]+)', result.stdout)
        if match:
            time = float(match.group(1))
            times.append(time)
        else:
            print(f"Warning: Could not extract time for run with {nodes} nodes")
    
    if times:
        return statistics.mean(times), statistics.stdev(times)
    else:
        return None, None

node_counts = [100, 200, 400, 800, 1600, 3200, 6400]
topology = 'line'
algorithm = 'gossip'

for nodes in node_counts:
    mean, stdev = run_gossip(nodes, topology, algorithm)
    if mean is not None and stdev is not None:
        print(f"Nodes: {nodes}, Mean Time: {mean:.2f} ms, Std Dev: {stdev:.2f} ms")
    else:
        print(f"Nodes: {nodes}, Failed to get valid measurements")