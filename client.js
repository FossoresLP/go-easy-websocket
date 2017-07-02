const wsHandlers = {};

function registerHandler(cmd, fnc) {
	wsHandlers[cmd] = fnc;
}

function messageHandler(msg) {
	let cmdEnd = msg.indexOf(": ");
	if (cmdEnd > 0 && cmdEnd <= 255) {
		let cmd = msg.substring(0, cmdEnd);
		let data = msg.substring(cmdEnd + 2);
		wsHandlers[cmd](data);
	}
}
