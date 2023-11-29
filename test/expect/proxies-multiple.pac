function FindProxyForURL (url, host) {
	var h = host.toLowerCase();
	var p1 = 'PROXY 192.168.41.1:8080';
	var p2 = 'SOCKS 192.168.41.1:1080';

	if (dnsDomainIs(h, 'domain-for-proxy-1.com')) return p1;
	if (dnsDomainIs(h, 'domain-for-proxy-2.com')) return p1;
	if (dnsDomainIs(h, 'domain-for-socks-1.com')) return p2;
	if (dnsDomainIs(h, 'domain-for-socks-2.com')) return p2;

	return 'DIRECT';
}
