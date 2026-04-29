const loginOutput = document.getElementById("loginOutput");
const userOutput = document.getElementById("userOutput");
const productsList = document.getElementById("productsList");
const ordersOutput = document.getElementById("ordersOutput");
const paymentOutput = document.getElementById("paymentOutput");
const healthOutput = document.getElementById("healthOutput");

document.getElementById("loginBtn").addEventListener("click", login);
document.getElementById("loadUserBtn").addEventListener("click", loadUser);
document.getElementById("loadProductsBtn").addEventListener("click", loadProducts);
document.getElementById("loadOrdersBtn").addEventListener("click", loadOrders);
document.getElementById("payBtn").addEventListener("click", pay);
document.getElementById("checkHealthBtn").addEventListener("click", checkHealth);

loadProducts();
loadOrders();

async function login() {
  const username = document.getElementById("usernameInput").value;
  const password = document.getElementById("passwordInput").value;

  const data = await request("/api/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });
  show(loginOutput, data);
}

async function loadUser() {
  const data = await request("/api/users/me");
  show(userOutput, data);
}

async function loadProducts() {
  const products = await request("/api/products");
  productsList.innerHTML = "";

  if (!Array.isArray(products)) {
    productsList.textContent = "Could not load products.";
    return;
  }

  products.forEach((product) => {
    const card = document.createElement("div");
    card.className = "product";
    card.innerHTML = `
      <strong>${escapeHTML(product.name)}</strong>
      <div class="product-meta">ID: ${product.id} | Price: $${product.price}</div>
      <div class="order-row">
        <input type="number" min="1" value="1" aria-label="Quantity for ${escapeHTML(product.name)}">
        <button>Create Order</button>
      </div>
    `;

    const quantityInput = card.querySelector("input");
    const button = card.querySelector("button");
    button.addEventListener("click", () => createOrder(product.id, quantityInput.value));
    productsList.appendChild(card);
  });
}

async function createOrder(productID, quantityValue) {
  const quantity = Number(quantityValue);
  const data = await request("/api/orders", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      user_id: 1,
      product_id: productID,
      quantity,
    }),
  });
  show(ordersOutput, data);
  await loadOrders();
}

async function loadOrders() {
  const data = await request("/api/orders");
  show(ordersOutput, data);
}

async function pay() {
  const data = await request("/api/payments", { method: "POST" });
  show(paymentOutput, data);
}

async function checkHealth() {
  const frontend = await requestText("/health");
  const gateway = await requestText("/gateway-health");
  show(healthOutput, {
    frontend,
    gateway: parseMaybeJSON(gateway),
  });
}

async function request(url, options) {
  try {
    const response = await fetch(url, options);
    const text = await response.text();
    const body = parseMaybeJSON(text);

    if (!response.ok) {
      return {
        error: true,
        status: response.status,
        body,
      };
    }

    return body;
  } catch (error) {
    return {
      error: true,
      message: error.message,
    };
  }
}

async function requestText(url) {
  try {
    const response = await fetch(url);
    const text = await response.text();
    if (!response.ok) {
      return `HTTP ${response.status}: ${text}`;
    }
    return text.trim();
  } catch (error) {
    return error.message;
  }
}

function parseMaybeJSON(text) {
  try {
    return JSON.parse(text);
  } catch {
    return text;
  }
}

function show(element, data) {
  element.textContent = JSON.stringify(data, null, 2);
}

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}
