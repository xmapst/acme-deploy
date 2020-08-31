我是光年实验室高级招聘经理。
我在github上访问了你的开源项目，你的代码超赞。你最近有没有在看工作机会，我们在招软件开发工程师，拉钩和BOSS等招聘网站也发布了相关岗位，有公司和职位的详细信息。
我们公司在杭州，业务主要做流量增长，是很多大型互联网公司的流量顾问。公司弹性工作制，福利齐全，发展潜力大，良好的办公环境和学习氛围。
公司官网是http://www.gnlab.com,公司地址是杭州市西湖区古墩路紫金广场B座，若你感兴趣，欢迎与我联系，
电话是0571-88839161，手机号：18668131388，微信号：echo 'bGhsaGxoMTEyNAo='|base64 -D ,静待佳音。如有打扰，还请见谅，祝生活愉快工作顺利。

# acme-deploy
acme renew, deploy nginx/kong/istio, Automatically detect expiration time

## Install 
#### Binary installation
Download the binary file suitable for your platform from [release](https://github.com/xmapst/acme-deploy/releases)  
```bash
wget https://github.com/xmapst/acme-deploy/releases/acme_linux_amd64.tar.gz
tar zxr -C /usr/local/bin acme_linux_amd64.tar.gz
```
#### Source installation  
go version go1.14.4 linux/amd64
```bash
git clone https://github.com/xmapst/acme-deploy.git
cd acme-deploy
chmod +x build.sh
./build.sh
```

## deploy nginx(example dnspod)
```bash
export DNSPOD_API_KEY=xxxxxxxxxxxxx  
export DNSPOD_HTTP_TIMEOUT=60  
export CERT_PATH=/path/ssl.cert  
export KEY_PATH=/path/ssl.key  
acme --dns dnspod -d *.example.com -m username@example.com --deploy nginx run  
```

## deploy kong(example dnspod)
```bash
export DNSPOD_API_KEY=xxxxxxxxxxxxx  
export DNSPOD_HTTP_TIMEOUT=60  
export KONG_ADMIN_ADDR=http://kong.kong:8001  
acme --dns dnspod -d *.example.com -m username@example.com --deploy kong run  
```

## deploy istio(example dnspod)
```bash
export DNSPOD_API_KEY=xxxxxxxxxxxxx  
export DNSPOD_HTTP_TIMEOUT=60  
export KUBECONF=/path/kube.config  
export NAMESPACE=istio-system  
export SECRETNAME=istio-ssl  
export CERT=ssl-cert  
export KEY=ssl-key  
acme --dns dnspod -d *.example.com -m username@example.com --deploy secret run  
```

## Help
```bash
$ acme --help
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
