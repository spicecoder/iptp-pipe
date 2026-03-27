function createAuthService() {
  return {
    id: "auth_service",
    runOnce: true,
    requires: [
      { name: "login_submitted", value: "Y" }
    ],
    async run() {
      await delay(300);
      return [
        { name: "user_authenticated", value: "Y" }
      ];
    }
  };
}

function delay(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

module.exports = { createAuthService };
