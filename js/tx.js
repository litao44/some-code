/*
 * 测试TPS
 */

var Web3 = require('web3');
var web3 = new Web3();
web3.setProvider(new web3.providers.HttpProvider("http://localhost:8545"));

//web3.eth.personal.unlockAccount("0xe08f28a9e8dbd7b701542ed377911148fabe7669", "123456", 3600);

for (let index = 0; index < 100; index++) {
    web3.eth.sendTransaction({
            from: "0x521a8600cdda197be03db0c64874db17ee5e91b1",
            to: "0x1a69cc48ef98b6f79fce95efa69e876e52b959a1",
            value: web3.utils.toWei("1", "ether")
        })
        .once('transactionHash', console.log);
}
