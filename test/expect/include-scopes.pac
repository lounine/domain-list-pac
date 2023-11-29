function FindProxyForURL (url, host) {
	var h = host.toLowerCase();
	var p1 = 'SOCKS 192.168.41.1:1080';
	var p2 = 'PROXY 192.168.41.1:8080';

	if (dnsDomainIs(h, 'domain-socks-1.com')) return p1;
	if (dnsDomainIs(h, 'some-domain.com')) return p2;
	if (dnsDomainIs(h, 'other-domain.com')) return p2;
	if (dnsDomainIs(h, 'domain-socks-2.com')) return p1;

	return 'DIRECT';
}
