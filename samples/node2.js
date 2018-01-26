var Discover = require('discover');
var TcpTransport = require('discover-tcp-transport');

var seedNode = {
	id: "111111111111111111111111111=",
	transport: {
		host: 'localhost',
		port: 9001
	}
};

var thisNode = {
	id: "222222222222222222222222222=",
	transport: {
		host: 'localhost',
		port: 9002
	}
};

var transport = TcpTransport.listen({port: thisNode.transport.port}, function () {
	var discover = new Discover({
		inlineTrace: true, 
		seeds: [
			seedNode
		],
		transport: transport
	});

	discover.register(thisNode);
});