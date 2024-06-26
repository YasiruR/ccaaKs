'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class HouseWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }

    // initializes the dataset for each test round
    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);

        for (let i=0; i<this.roundArguments.houses; i++) {
            const assetID = `${this.workerIndex}${i}`
            console.log(`Worker ${this.workerIndex}: Creating asset ${assetID}`);
            const req = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'CreateHouse',
                invokerIdentity: 'peer1',
                contractArguments: ['blue', assetID, 'Alex', '766'],
                readOnly: false
            };

            await this.sutAdapter.sendRequests(req);
        }
    }

    // transaction subject to benchmarks
    async submitTransaction() {
        const randId = Math.floor(Math.random()*this.roundArguments.houses);
        const newArgs = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'GetHouse',
            invokerIdentity: 'peer1',
            contractArguments: [`${this.workerIndex}${randId}`],
            readOnly: true,
        };

        await this.sutAdapter.sendRequests(newArgs);
    }

    // cleaning up created houses
    async cleanupWorkloadModule() {
        for (let i=0; i<this.roundArguments.houses; i++) {
            const assetID = `${this.workerIndex}${i}`;
            console.log(`Worker ${this.workerIndex}: Deleting asset ${assetID}`);
            const req = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'DeleteHouse',
                invokerIdentity: 'peer1',
                contractArguments: [assetID],
                readOnly: false
            };

            await this.sutAdapter.sendRequests(req);
        }
    }
}

function createWorkloadModule() {
    return new HouseWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
