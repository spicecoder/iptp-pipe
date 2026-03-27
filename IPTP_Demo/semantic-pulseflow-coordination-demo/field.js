function pulseKey(pulse) {
  return `${pulse.name}:${pulse.value}`;
}

class Field {
  constructor(initialPulses = []) {
    this.map = new Map();
    this.absorbMany(initialPulses);
  }

  has(pulse) {
    return this.map.has(pulseKey(pulse));
  }

  absorb(pulse) {
    const key = pulseKey(pulse);
    if (this.map.has(key)) {
      return false;
    }
    this.map.set(key, { ...pulse });
    return true;
  }

  absorbMany(pulses) {
    let changed = false;
    for (const pulse of pulses) {
      changed = this.absorb(pulse) || changed;
    }
    return changed;
  }

  all() {
    return [...this.map.values()].sort((a, b) =>
      pulseKey(a).localeCompare(pulseKey(b))
    );
  }

  print() {
    console.log("FIELD:");
    for (const pulse of this.all()) {
      console.log(`  ${pulse.name}:${pulse.value}`);
    }
    console.log("");
  }
}

module.exports = { Field };
