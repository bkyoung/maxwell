# Maxwell

Maxwell is a secrets agent.  It runs as a system daemon that can interact with many different secret storage backends.  

## Purpose

Maxwell is meant to be a specialized solution to the "secret zero" problem, in cases where node attestation may not be feasible or otherwise make
sense, but there still exists a need to securely connect a new, previously unknown system to another over an untrusted intermediate network 
unattended.
