const path = require("path");
const fs = require("fs");

// connection.json 객체화
const ccpPath = path.resolve(__dirname, "..", "ccp", "connection-org1.json");
const ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));
const FabricCAServices = require("fabric-ca-client");
const { Gateway, Wallets } = require("fabric-network");


async function makeAdminWallet() {
    const id = "admin";
    const pw = "adminpw";

    console.log(id, pw);

    try {
        // Create a new CA client for interacting with the CA.
        const caInfo = ccp.certificateAuthorities["ca.org1.example.com"];
        const caTLSCACerts = caInfo.tlsCACerts.pem;
        const ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), "wallet");
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the admin user.
        const identity = await wallet.get(id);
        if (identity) {
        console.log(`An identity for the admin user ${id} already exists in the wallet`);
        // const res_str = `{"result":"failed","msg":"An identity for the admin user ${id} already exists in the wallet"}`;
        // res.json(JSON.parse(res_str));
        return;
        }

        // Enroll the admin user, and import the new identity into the wallet.
        const enrollment = await ca.enroll({ enrollmentID: id, enrollmentSecret: pw });
        const x509Identity = {
        credentials: {
            certificate: enrollment.certificate,
            privateKey: enrollment.key.toBytes(),
        },
        mspId: "Org1MSP",
        type: "X.509",
        };
        await wallet.put(id, x509Identity);

        // response to client
        console.log('Successfully enrolled admin user "admin" and imported it into the wallet');
        // const res_str = `{"result":"success","msg":"Successfully enrolled admin user ${id} in the wallet"}`;
        // res.status(200).json(JSON.parse(res_str));
    } catch (error) {
        console.error(`Failed to enroll admin user ${id} : ${error}`);
        // const res_str = `{"result":"failed","msg":"failed to enroll admin user - ${id} : ${error}"}`;
        // res.json(JSON.parse(res_str));
    }
}

module.exports.makeAdminWallet = makeAdminWallet;