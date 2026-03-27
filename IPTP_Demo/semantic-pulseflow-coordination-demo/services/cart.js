function createCartService() {
  return {
    id: "cart_service",
    runOnce: true,
    requires: [
      { name: "cart_submitted", value: "Y" }
    ],
    async run() {
      await delay(500);
      return [
        { name: "cart_valid", value: "Y" }
      ];
    }
  };
}

function delay(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

module.exports = { createCartService };
