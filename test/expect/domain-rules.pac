function FindProxyForURL (url, host) {
	var h = host.toLowerCase();
	var p = 'PROXY 192.168.41.1:8080';

	if (dnsDomainIs(h, 'sub-domains-1.com')) return p;
	if (dnsDomainIs(h, 'sub-domains-2.com')) return p;
	if (h == 'www.full-domain-1.com') return p;
	if (h == 'www.full-domain-2.com') return p;
	if (shExpMatch(h, '*some-keyword*')) return p;
	if (/^some-regexp(-[1-9]+)?\.com$/.test(h)) return p;

	return 'DIRECT';
}
