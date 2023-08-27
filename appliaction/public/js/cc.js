const path = require("path");
const fs = require("fs");

// connection.json 객체화
const ccpPath = path.resolve(__dirname, "..", "ccp", "connection-org1.json");
const ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));
const { Wallets, Gateway  } = require("fabric-network");


async function cc_call(id, fn_name, args) {
    try {
        const walletPath = path.join(process.cwd(), "wallet");
        console.log(`Wallet path in cc_call: ${walletPath}`);
        const wallet = await Wallets.newFileSystemWallet(walletPath);

        console.log(`id: ${id}`);

        const userExists = await wallet.get(id);
        if (!userExists) {
            console.log(
                `An identity for the user "${id}" does not exist in the wallet`
            );
            console.log("Run the registerUser.js application before retrying");
            return;
        }

        const gateway = new Gateway();
        await gateway.connect(ccp, {
            wallet,
            identity: id,
            discovery: { enabled: true, asLocalhost: true },
        });

        const network = await gateway.getNetwork("mychannel");
        const contract = network.getContract("fungi");

        var result;

        if (fn_name == "GetFungiByOwner") {
        
            result = await contract.evaluateTransaction("GetFungiByOwner");
        // } else if (fn_name == "addRating") {
        //     e = args[0];
        //     p = args[1];
        //     s = args[2];
        //     result = await contract.submitTransaction("addRating", e, p, s);
        } else if (fn_name == "CreateRandomFungus")
            result = await contract.submitTransaction("CreateRandomFungus", args);
        else result = "not supported function";

        console.log(`result in CC_call: ${result}`);

        return result;
        
    } catch (error) {
        console.error(error)
    }
    
}

module.exports.cc_call = cc_call;