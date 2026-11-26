import socket, time

s = socket.create_connection(("localhost", 8080))
f = s.makefile("rw", newline="\n")
f.readline()
f.write("CONNECT probe%d\n" % int(time.time())); f.flush()
time.sleep(1)

times = []
for _ in range(50):
    t0 = time.perf_counter()
    f.write("who\n"); f.flush()
    while not f.readline().startswith("OK players="):
        pass
    times.append((time.perf_counter() - t0) * 1000)
    time.sleep(0.1)

print("moyenne : %.2f ms | max : %.2f ms" % (sum(times) / len(times), max(times)))
