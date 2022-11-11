package main

/*
	echo is a very simple implementation of an indexed storaged of block
	history.

	it connects to one or more validators and receive new blocks from them.
	blocks are stored on local files and indexes from involved tokens and their
	respective instructions positions on files are created.

	and instance of echo can be connected to another instance of echo in order
	to reeceive missing blocks. recent blocks can be downloaded directly from
	validators since the protocol requires them to hold information of most
	recent blocks.

	on the othe side, echo offers an API through which clients can manifest
	interest to follow instructions on one or more tokens. All instructions
	associated to them are sent to the clients together with instructions for
	new blocks. At connection client must offer the tokens and the most recent
	block to which information is available.

	clients can also request the hash history of blocks up to the genesis block.

*/
