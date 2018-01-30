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
	id: "444444444444444444444444444=",
	transport: {
		host: 'localhost',
		port: 9004
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

	setTimeout(function() {
		discover.find("222222222222222222222222222=", function (error, contact) {
			if (error) {
				console.log(error)
			}
			else {
				console.log("\r\n\r\nTRACE: " + JSON.stringify(contact));
			}
		});
	}, 5000);
});