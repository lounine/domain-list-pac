function FindProxyForURL (url, host) {
	var h = host.toLowerCase();
	var p = 'PROXY 192.168.41.1:8080';

	if (dnsDomainIs(h, 'domain.com')) return p;
	if (dnsDomainIs(h, 'some-domain.com')) return p;
	if (dnsDomainIs(h, 'other-domain.com')) return p;

	return 'DIRECT';
}
