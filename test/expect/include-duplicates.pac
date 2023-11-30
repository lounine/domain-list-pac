function FindProxyForURL (url, host) {
	var h = host.toLowerCase();
	var p = 'PROXY 192.168.41.1:8080';

	if (dnsDomainIs(h, 'original-domain-1.com')) return p;
	if (dnsDomainIs(h, 'duplicate-domain-1.com')) return p;
	if (dnsDomainIs(h, 'duplicate-domain-2.com')) return p;
	if (dnsDomainIs(h, 'original-domain-2.com')) return p;

	return 'DIRECT';
}
