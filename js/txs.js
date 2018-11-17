var Web3 = require('web3');
var util = require('util');

var web3 = new Web3();
web3.setProvider(new web3.providers.HttpProvider("http://localhost:9045"));

var coinbase = "0xb72355fdbd78e5ce7f7a147178611d761805b5ee"; //马云爸爸,very very rich

var accounts = [ //需要解锁 personal.unlockAccount()
    "0xb72355fdbd78e5ce7f7a147178611d761805b5ee",
];

var receivers = [
    "0x94485e4ab3726b9fecc4fbcbd757deef6e28ac16",
    "0x3dbc4ba9084cc94b5e52f9a6a78ee2a324df6520",
    "0x5a01afd40880a3f1abb762cd2742139d67db8f33",
    "0xa3657125dd88c05d149fa4055cb088e909f75fb4",
    "0xda20d7b21c11fc10cd70f11da293e7e716601a27",
    "0x86f34627fdd1dbd437a8e1fae9e41f2369f23f45",
    "0x9f0aeba1bc72aa8b8c2e07b110944aec8a5d75ea",
];

function sendRandomTxs() { // TODO, 检测pendingTransactions，数量过多就不发了
    var num = getRandomIntInclusive(5, 10);
    console.log(util.format("本次发送 %d 个交易, hash:", num));
    for (var i = 0; i < num; i++) {
        sendRandomTx();
    }
}

function sendRandomTx() {
    var value = getRandomArbitrary(1, 20);
    sendTransaction(getArrayRandom(accounts), getArrayRandom(receivers), value.toString());
}

function getArrayRandom(arr) {
    return arr[getRandomIntInclusive(0, arr.length - 1)];
}

function sendTransaction(from, to, value) {
    web3.eth.sendTransaction({
            from: from,
            to: to,
            value: web3.utils.toWei(value, "ether")
        })
        .once('transactionHash', function (hash) {
            console.log(util.format("from: [%s],to: [%s], value: [%s], hash: [%s]", from, to, value, hash));
        })
        .on('error', function (error, receipt) {
            console.log(util.format("from: [%s],to: [%s], value: [%s], error: [%s]", from, to, value, error));
        });
}

function getRandomIntInclusive(min, max) {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min + 1)) + min; //The maximum is inclusive and the minimum is inclusive
}

function getRandomArbitrary(min, max) {
    return Math.random() * (max - min) + min;
}

function initAccount() {
    accounts.forEach(function (account) {
        sendTransaction(coinbase, account, "1000000000");
    });
}


initAccount();
setInterval(sendRandomTxs, 1000 * 2);
