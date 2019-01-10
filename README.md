# certcheck

Batch verify of TLS/SSL Certificates

Lookup is handled concurrently and each request times out after 3 seconds.

## Usage

Certcheck accepts input from stdin (piped) or from provided file (using the -filename flag).

## Examples

`cat ssdomains | ./certcheck`

`./certcheck -filename ssldomains`

## Options

 * -filename (string): The name of the file containing the target domain names, one per line.
 * -num-routines (int): The number of routines that will process this data concurrently. (default 4)

