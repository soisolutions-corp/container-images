# Wait All

A custom init process written in Go that waits for all child processes to exit before exiting itself.

## Why?

Cribl Stream forks multiple Worker Processes on start up. This can cause issues with process management, especially when running in a container. When Cribl stops or restarts, the primary process will exit before all Worker Processes exit. This can result in data loss due to improper shutdown while waiting for data to flush to destinations.

By providing an init process that waits for all child processes to exit before exiting itself, we can ensure a clean shutdown of all Worker Processes.
