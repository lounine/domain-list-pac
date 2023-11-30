function FindProxyForURL (url, host) {
	var h = host.toLowerCase();
	var p = 'PROXY 192.168.41.1:8080';

	if (dnsDomainIs(h, 'sub-domains-1.com')) return p;
	if (dnsDomainIs(h, 'sub-domains-2.com')) return p;
	if (h == 'www.full-domain.com') return p;
	if (shExpMatch(h, '*some-keyword*')) return p;

	return 'DIRECT';
}
