const { Field } = require("./field");
const { Engine } = require("./engine");
const { createAuthService } = require("./services/auth");
const { createCartService } = require("./services/cart");
const { createOrderService } = require("./services/order");

async function main() {
  console.log("PulseFlow Coordination Demo");
  console.log("");
  console.log("Problem:");
  console.log("  In independently built systems, the logic that determines when something");
  console.log("  should happen often ends up spread across code, configs, or orchestration.");
  console.log("");
  console.log("This demo shows a different approach:");
  console.log("  - auth_service reacts to login_submitted:Y");
  console.log("  - cart_service reacts to cart_submitted:Y");
  console.log("  - order_service reacts only when both user_authenticated:Y and cart_valid:Y exist");
  console.log("");
  console.log("Important:");
  console.log("  - no service calls another service");
  console.log("  - no central workflow tells order_service when to run");
  console.log("  - the pass loop only gives each service a chance to check readiness");
  console.log("  - semantic dependency is based on field state, not service position");
  console.log("");
  console.log("Deliberately awkward service order:");
  console.log("  order_service -> cart_service -> auth_service");
  console.log("");

  const field = new Field([
    { name: "login_submitted", value: "Y" },
    { name: "cart_submitted", value: "Y" }
  ]);

  const services = [
    createOrderService(),
    createCartService(),
    createAuthService()
  ];

  const engine = new Engine(field, services);
  await engine.runUntilStable();

  console.log("FINAL FIELD:");
  field.print();
}

main().catch((error) => {
  console.error("Fatal error:", error);
  process.exit(1);
});