function createOrderService() {
  return {
    id: "order_service",
    runOnce: true,
    requires: [
      { name: "user_authenticated", value: "Y" },
      { name: "cart_valid", value: "Y" }
    ],
    async run() {
      await delay(200);
      return [
        { name: "order_created", value: "Y" }
      ];
    }
  };
}

function delay(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

module.exports = { createOrderService };
