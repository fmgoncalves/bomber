bomber
======

Simple program to send huge amounts of predefined email samples to an SMTP server.
Just an excuse to try out Go, really...

```
Usage of bomber:
  -c="all": Category of messages to send
  -l=false: Only list available categories.
  -n=100: Number of message to be sent
  -s="localhost": SMTP server to send the message to
  -samples="./samples/": Directory containing email message samples in JSON
  -throttle=0: Throttle the message flow (msg/second).
  -v=false: Prints details of execution.
```
