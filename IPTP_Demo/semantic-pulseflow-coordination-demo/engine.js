class Engine {
  constructor(field, services) {
    this.field = field;
    this.services = services;
    this.executed = new Set();
  }

  serviceReady(service) {
    if (service.runOnce && this.executed.has(service.id)) {
      return false;
    }
    return service.requires.every((pulse) => this.field.has(pulse));
  }

  async runService(service) {
    console.log(`[${service.id}] READY`);
    const emitted = await service.run();
    console.log(`[${service.id}] EXECUTED`);

    const changed = [];
    for (const pulse of emitted) {
      if (this.field.absorb(pulse)) {
        changed.push(pulse);
      }
    }

    this.executed.add(service.id);

    if (changed.length === 0) {
      console.log(`[${service.id}] emitted no new pulses\n`);
      return false;
    }

    console.log(`[${service.id}] EMITTED:`);
    for (const pulse of changed) {
      console.log(`  ${pulse.name}:${pulse.value}`);
    }
    console.log("");
    return true;
  }

  async runUntilStable() {
    let pass = 1;

    while (true) {
      console.log(`=== PASS ${pass} ===`);
      this.field.print();

      let progressed = false;
      let readyCount = 0;

      for (const service of this.services) {
        if (this.executed.has(service.id) && service.runOnce) {
          console.log(`[${service.id}] ALREADY EXECUTED`);
          continue;
        }

        if (this.serviceReady(service)) {
          readyCount += 1;
          const changed = await this.runService(service);
          progressed = changed || progressed;
        } else {
          console.log(`[${service.id}] NOT READY`);
        }
      }

      console.log("");

      if (readyCount === 0) {
        console.log("No more ready services. Execution complete.\n");
        break;
      }

      if (!progressed) {
        console.log("No field change in this pass. Stopping.\n");
        break;
      }

      pass += 1;
    }
  }
}

module.exports = { Engine };