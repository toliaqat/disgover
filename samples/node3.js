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
	id: "333333333333333333333333333=",
	transport: {
		host: 'localhost',
		port: 9003
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