"use strict";

var assert = require('assert'),
	crypto = require('crypto'),
	Discover = require('discover'),
	TcpTransport = require('discover-tcp-transport'),
	util = require('util');

var transport1, transport2;
var discover1, discover2;
var id1, id2;

transport1 = TcpTransport.listen({port: 6741}, function () {
	transport2 = TcpTransport.listen({port: 6742}, function () {
		startLocalTest();
	});
});

function startLocalTest() {

	id1 = crypto.randomBytes(20).toString('base64');
	id2 = crypto.randomBytes(20).toString('base64');

	discover1 = new Discover({
		inlineTrace: false, 
		seeds: [
			{
				id: id1,
				data: 'discover1',
				transport: {
					host: 'localhost',
					port: 6741
				}
			},
			{
				id: id2,
				data: 'discover2',
				transport: {
					host: 'localhost',
					port: 6742
				}
			}
		],
		transport: transport1
	});
	discover2 = new Discover({
		inlineTrace: false, 
		seeds: [
			{
				id: id1,
				data: 'discover1',
				transport: {
					host: 'localhost',
					port: 6741
				}
			},
			{
				id: id2,
				data: 'discover2',
				transport: {
					host: 'localhost',
					port: 6742
				}
			}
		],
		transport: transport2
	});

	console.log('~script five discover instances running');
	console.log('~script starting self-registrations');

	var node1 = {id: id1, data: 'discover1', transport: {host: 'localhost', port: 6741}};
	var node2 = {id: id2, data: 'discover2', transport: {host: 'localhost', port: 6742}};

	discover1.register(node1);
	discover2.register(node2);

	console.log('~script self-registrations complete');
	console.log('~script allowing nodes to communicate and settle');

	setTimeout(continueLocalTest2, 1000);
};

var id6;

function continueLocalTest2() {
	id6 = crypto.randomBytes(20).toString('base64');
	var node6 = {id: id6, data: 'discover6', transport: {host: 'localhost', port: 6741}};

	console.log('~script multiple nodes per discover instance');

	discover1.register(node6);

	console.log('~script allowing nodes to communicate and settle');

	setTimeout(continueLocalTest3, 1000);
};

function continueLocalTest3() {
	console.log('~script retrieving node6 from discover2');
	discover2.find(id6, function (error, contact) {
		assert.ok(!error);

		console.log('~script recevied contact: ' + util.inspect(contact, false, null));
	});
};

function complete() {
	console.log('~script complete');
	setTimeout(function () { process.exit(); }, 1000);
};