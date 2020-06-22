# acme-deploy
acme renew, deploy nginx/kong/istio

## Help
```bash
$ ./acme --help
NAME:
   acme - Let's Encrypt client written in Go

USAGE:
   acme [global options] command [command options] [arguments...]

VERSION:
   dev

COMMANDS:
   run      Register an account, then create and install a certificate
   dns      Shows additional help for the '--dns' global option
   deploy   Shows additional help for the '--deploy' global option
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --domains value, -d value   Add a domain to the process. Can be specified multiple times.
   --server value, -s value    CA hostname (and optionally :port). The server certificate must be trusted in order to avoid further modifications to the client. (default: "https://acme-v02.api.letsencrypt.org/directory")
   --email value, -m value     Email used for registration and recovery contact.
   --csr value, -c value       Certificate signing request filename, if an external CSR is to be used.
   --eab                       Use External Account Binding for account registration. Requires --kid and --hmac.
   --kid value                 Key identifier from External CA. Used for External Account Binding.
   --hmac value                MAC key from External CA. Should be in Base64 URL Encoding without padding format. Used for External Account Binding.
   --key-type value, -k value  Key type to use for private keys. Supported: rsa2048, rsa4096, rsa8192, ec256, ec384. (default: "rsa4096")
   --deploy value              Key type to use for deploy type. Supported: secret, kong, nginx. Run 'acme deploy' for help on usage.
   --dns value                 Solve a DNS challenge using the specified provider. Can be mixed with other types of challenges. Run 'acme dns' for help on usage.
   --dns.disable-cp            By setting this flag to true, disables the need to wait the propagation of the TXT record to all authoritative name servers.
   --dns.resolvers value       Set the resolvers to use for performing recursive DNS queries. Supported: host:port. The default is to use the system resolvers, or Google's DNS resolvers if the system's cannot be determined.
   --http-timeout value        Set the HTTP timeout value to a specific value in seconds. (default: 10)
   --dns-timeout value         Set the DNS timeout value to a specific value in seconds. Used only when performing authoritative name servers queries. (default: 10)
   --cert.timeout value        Set the certificate timeout value to a specific value in seconds. Only used when obtaining certificates. (default: 30)
   --days value                The number of days left on a certificate to renew it. (default: 20)
   --help, -h                  show help
   --version, -v               print the version
```
