# Example Opentracing App

This is a trivial example app that demonstrates how OpenTracing can be
used with!

To run the program, run `./opentracing-example`. This will, by default,
start a new Appdash server and write all of your traces to it. However,
if you want to use a different tracer system, i.e. Lighstep, all you have
to do is pass the flag `--lighstep.token=ACCESS_TOKEN`.

## Building

Make sure that you have Go installed, then run `go build`.
