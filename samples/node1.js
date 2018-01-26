var Discover = require('discover');
var TcpTransport = require('discover-tcp-transport');

var thisNode = {
	id: "111111111111111111111111111=",
	transport: {
		host: 'localhost',
		port: 9001
	}
};









var transport = TcpTransport.listen({port: thisNode.transport.port}, function () {
	var discover = new Discover({
		inlineTrace: true, 
		seeds: [

		],
		transport: transport
	});

	discover.register(thisNode);
});