'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class VehicleWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }

    // initializes the dataset for each test round
    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
    }

    // transaction subject to benchmarks
    async submitTransaction() {
        const randId = Math.floor(Math.random()*this.roundArguments.vehicles);
        const vehicleID = `${this.workerIndex}${randId}`
        const newArgs = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'CreateVehicle',
            invokerIdentity: 'peer1',
            contractArguments: ['blue', vehicleID, 'Alex', '766'],
            readOnly: false,
        };
        //console.log(`RandID ${randId}: Creating vehicle ${vehicleID}`);

        await this.sutAdapter.sendRequests(newArgs);
    }

    // cleaning up created vehicles
    async cleanupWorkloadModule() {
        for (let i=0; i<this.roundArguments.vehicles; i++) {
            const vehicleID = `${this.workerIndex}${i}`;
            console.log(`Worker ${this.workerIndex}: Deleting vehicle ${vehicleID}`);
            const req = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'DeleteVehicle',
                invokerIdentity: 'peer1',
                contractArguments: [vehicleID],
                readOnly: false
            };

            await this.sutAdapter.sendRequests(req);
        }
    }
}

function createWorkloadModule() {
    return new VehicleWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
